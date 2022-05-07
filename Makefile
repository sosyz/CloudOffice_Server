all:
	go build -o main
put:
	export CGO_ENABLED=0
	go build -o build/CloudPrinter/main -a -ldflags '-extldflags "-static"' .
	serverless deploy --target build/CloudPrinter/ --force