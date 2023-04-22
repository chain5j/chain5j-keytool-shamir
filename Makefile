linux:
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build -o keytools genetool/main.go
	@echo "Done building."
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o keytools main.go
	@echo "Done building."
mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o keytools main.go
	@echo "Done building."

