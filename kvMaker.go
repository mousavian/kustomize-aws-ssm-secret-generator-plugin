package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type plugin struct{}

var KVSource plugin

type options struct {
	AwsRegion             string `alias:"AWS_REGION"`
	AwsAccessKeyID        string `alias:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey    string `alias:"AWS_SECRET_ACCESS_KEY"`
	AwsSessionToken       string `alias:"AWS_SESSION_TOKEN"`
	AwsParameterStorePath string `alias:"AWS_SSM_PATH"`

	UppercaseKeys string `alias:"UPPERCASE_KEY"`
}

func (p plugin) Get(root string, args []string) (map[string]string, error) {
	r := make(map[string]string)
	opts := &options{}
	cfg := &aws.Config{}
	opts.parseArgs(&args)

	if opts.AwsParameterStorePath == "" {
		return r, fmt.Errorf("AWS_SSM_PATH is required")
	}

	if opts.AwsRegion != "" {
		cfg.Region = aws.String(opts.AwsRegion)
	}

	if opts.AwsAccessKeyID != "" && opts.AwsSecretAccessKey != "" {
		staticCreds := credentials.NewStaticCredentials(opts.AwsAccessKeyID, opts.AwsSecretAccessKey, opts.AwsSessionToken)
		cfg.WithCredentials(staticCreds)
	}

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	svc := ssm.New(sess)

	getParamsInput := &ssm.GetParametersByPathInput{
		Path:           aws.String(opts.AwsParameterStorePath),
		WithDecryption: aws.Bool(true),
	}

	for {
		resp, err := svc.GetParametersByPath(getParamsInput)

		if err != nil {
			return nil, err
		}

		params := resp.Parameters
		for _, p := range params {
			name := sanitizeKey(p.Name, opts.UppercaseKeys == "true")
			r[name] = sanitizeValue(p.Value)
		}

		nextToken := resp.NextToken
		if nextToken == nil {
			break
		}
	}

	return r, nil
}

func sanitizeKey(path *string, ensureUppercase bool) string {
	pathBrokenDown := strings.Split(aws.StringValue(path), "/")
	name := pathBrokenDown[len(pathBrokenDown)-1]

	if ensureUppercase {
		name = strings.ToUpper(name)
	}

	return name
}

func sanitizeValue(path *string) string {
	return aws.StringValue(path)
}

func (opts *options) parseArgs(args *[]string) error {
	ov := reflect.ValueOf(opts).Elem()
	typeOfOpts := ov.Type()

	for _, arg := range *args {
		argKeyValuePair := strings.Split(arg, "=")
		if len(argKeyValuePair) == 2 {
			for i := 0; i < ov.NumField(); i++ {
				field := ov.Field(i)
				fieldTag := typeOfOpts.Field(i).Tag

				if argKeyValuePair[0] == fieldTag.Get("alias") {
					field.SetString(argKeyValuePair[1])
				}
			}
		}
	}

	return nil
}
