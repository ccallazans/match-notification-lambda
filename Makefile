build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/ cmd/main.go

zip:
	zip function.zip main