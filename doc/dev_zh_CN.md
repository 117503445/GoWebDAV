# 开发

## 开发环境搭建

```sh
docker compose up -d
```

然后使用 VSCode 附加至 `gowebdav-dev` 容器，进入 `/workspace` 开发。

## 常用开发操作

```sh
go run . # 运行
go run . --port 8080 # 运行并指定端口

go build . # 构建二进制
go test ./... # 测试

docker build -t 117503445/go_webdav . # 构建 Docker 镜像
```
