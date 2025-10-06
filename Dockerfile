FROM golang:1.25.1 AS build
ARG GOPROXY=https://proxy.golang.org,direct
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go test ./...
RUN CGO_ENABLED=0 go build -o app

FROM scratch
COPY --from=build /workspace/app /app
ENV dav="/,/data,null,null,false"
EXPOSE 80
VOLUME [ "/data" ]
ENTRYPOINT [ "/app"]