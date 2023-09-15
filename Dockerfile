FROM golang:1.21 as build
WORKDIR /go/release
COPY go.mod .
COPY go.sum .
COPY mod.sh .
RUN ./mod.sh
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o app
FROM alpine as prod
EXPOSE 80
WORKDIR /root
COPY --from=build /go/release/app app
ENTRYPOINT ./app --dav=$dav