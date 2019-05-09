# Kustomize Secret Generator Plugin for AWS SSM Parameter Store

This plugin can be attached to [Kustomize](https://kustomize.io/) to generate Kubernetes secrets automatically from parameters in [Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html) of AWS Simple System Manager Service (SSM).

This can be useful for CI/CD.

## Usage
Download latest `aws-ssm.so` file from releases and put it in your Kustomize plugin directory (default to `~/.config/kustomize/plugin/kvSources/`)

and use it in your kustomization file:
```yaml
secretGenerator:
  - name: my-secret-name
    kvSources:
      - name: aws-ssm
        pluginType: go
        args:
          - AWS_SSM_PATH=/path/to/my/secrets/ # Required
          - AWS_REGION=ap-southeast-1         # Optional
          - AWS_ACCESS_KEY_ID=                # Optional
          - AWS_SECRET_ACCESS_KEY=            # Optional
          - AWS_SESSION_TOKEN=                # Optional
          - UPPERCASE_KEY=true                # Optional
```

Assuming you have two parameter under `/path/to/my/secrets/` such as:

`/path/to/my/secrets/key1` with value of `value1` and

`/path/to/my/secrets/key2` with value of `value2`

the output will be

```yaml
apiVersion: v1
data:
  KEY1: dmFsdWUx
  KEY2: dmFsdWUy
kind: Secret
metadata:
  name: my-secret-name-someRandomHash
type: Opaque
```


### Note
Note that this feature of Kustomize is alpha and is not released yet.
So to test you have to build it from master branch and run it with `enable_alpha_goplugins_accept_panic_risk` parameter like:

```
kustomize --enable_alpha_goplugins_accept_panic_risk build ./kustomization.yaml | kubectl apply -f -
```
