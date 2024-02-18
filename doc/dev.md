# Development

## Setting Up Development Environment

```sh
docker compose up -d
```

Then attach to the `gowebdav-dev` container using VSCode and start development in the `/workspace` directory.

## Common Development Operations

```sh
go run . # Run
go run . --port 8080 # Run and specify port

go build . # Build binary
go test ./... # Test

docker build -t 117503445/go_webdav . # Build Docker image
```