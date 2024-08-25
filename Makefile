BINARY_NAME=ginServer
.DEFAULT_GOAL := dev
dev: 
	go run ./cmd/main.go
