

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o build/Dude_darwin_amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o build/Dude_darwin_arm64 main.go

	GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o build/Dude_windows_386 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/Dude_windows_amd64 main.go
	
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/Dude_linux_amd64 main.go
	