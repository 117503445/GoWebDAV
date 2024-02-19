FROM golang:1.22 as build
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app

FROM gcr.io/distroless/static-debian12 as prod
WORKDIR /workspace
COPY --from=build /workspace/app app
ENTRYPOINT [ "./app"]