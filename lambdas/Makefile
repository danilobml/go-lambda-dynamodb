build:
	@GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -tags lambda.norpc -o bootstrap .

zip:
	@zip go_lambda.zip bootstrap

package: build zip

clean:
	@go clean
	@rm go_lambda.zip bootstrap
