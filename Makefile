chat-server:
	CHAT_PORT=3000 go run ./cmd/chat/* server

chat-client:
	CHAT_PORT=3000 go run ./cmd/chat/* client

sneaky-server:
	SNEAKY_PORT=3000 go run ./cmd/sneaky/*.go

test:
	go test ./...
