.PHONY: clean build test

clean:
	@rm -rf ./build
	@mkdir -p ./build

build:
	@go build -buildmode plugin -o ./build/aws-ssm.so ./kvMaker.go
	@cp ./build/aws-ssm.so ~/.config/kustomize/plugin/kvSources/

test:
	@./bin/kustomize --enable_alpha_goplugins_accept_panic_risk build .

all: clean build test
