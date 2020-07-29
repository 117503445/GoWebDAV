# GoWebdav

> A Simple, Powerful WebDAV Server By Golang.

```sh
docker rm go_webdav -f
docker rmi 117503445/go_webdav
docker run -it --name go_webdav -d -e dav="/dav1,./TestDir1,user1,pass1;/dav2,./TestDir2,user2,pass2" -p 80:80 --restart=always 117503445/go_webdav:latest
```

## 本地调试

把 config.yml.example 重命名为 config.yml， 在 config.yml 文件中配置，再按照常规操作运行

使用了分层构建，在 build 层 通过 go build 构筑了 可执行文件 app，再在 prod 层 进行运行。如果以后需要修改配置文件的结构，也需要修改 Dockerfile。

## 本地 Docker 运行

```sh
docker rm go_webdav -f
docker rmi 117503445/go_webdav

docker build -t 117503445/go_webdav . # 国外
docker build -t 117503445/go_webdav -f Dockerfile_cn . #国内,启用 go 镜像

docker run -it --name go_webdav -d -e dav="/dav1,./TestDir1,user1,pass1;/dav2,./TestDir2,user2,pass2" -p 80:80 --restart=always 117503445/go_webdav:latest
```
