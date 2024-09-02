# GoWebDAV

> 使用 WebDAV 分享本地文件，轻量、易于使用

[English](./README.md) | 简体中文

## 特性

- 基于 Golang 实现，性能高。

- 最终编译为单二进制文件，不需要 Apache 等环境，依赖少。

- 支持浏览器访问。

- 可以在同个端口下启用多个 WebDAV 服务，各自有不同的挂载目录、用户名密码。

- 良好的 Docker 支持。

## 快速开始

从 <https://github.com/117503445/GoWebDAV/releases> 下载二进制文件

运行

```sh
./gowebdav
```

GoWebDAV 会自动在 `./data` 路径下创建示例文件，文件结构如下

```sh
> tree ./data
./data
├── public-writable
│   └── 1.txt
├── public-readonly
│   └── 1.txt
└── private-writable
    └── 1.txt
```

使用浏览器访问 <http://localhost:80>，就可以看到 3 个不同的 GoWebDAV 服务了。

![index](./doc/index.png)

其中 <http://localhost:80/public-writable> 是 `public-writable` 服务，映射了本地的 `./data/public-writable` 文件夹。它是无用户验证的、可写的。可以在浏览器中查看文件内容，也可以进行上传、删除等操作。

![public-writable](./doc/public-writable.png)

<http://localhost:80/public-readonly> 是 `public-readonly` 服务，映射了本地的 `./data/public-readonly` 文件夹。它是无用户验证的、只读的。可以在浏览器中查看文件内容，但不可以进行上传、删除等操作。

![public-readonly](./doc/public-readonly.png)

<http://localhost:80/private-writable> 是 `private-writable` 服务，映射了本地的 `./data/private-writable` 文件夹。它是有用户验证的、可写的。在使用 `user1` 和 `pass1` 登录以后，可以在浏览器中查看文件内容，也可以进行上传、删除等操作。

![private-writable](./doc/private-writable.png)

当然，除了浏览器，也可以使用其他 WebDAV 客户端工具进行访问。

可以通过指定 `dav` 参数来配置 WebDAV 服务的本地路径、用户验证、是否只读等属性，详情见 *使用* 章节。

## 使用

```sh
./gowebdav --help # 查看帮助

./gowebdav --addr 127.0.0.1 # 在 127.0.0.1 监听，默认监听 0.0.0.0
./gowebdav --port 8080 # 在 8080 端口监听，默认监听 80 端口

./gowebdav --dav "/dir1,/data/dir1,user1,pass1,true" # 配置文件夹路径及属性
```

`dav` 参数可以指定 WebDAV 服务的本地路径、用户验证、是否只读等属性。

每个本地路径都可以配置一个 WebDAV 服务，使用分号分隔。例子：

- `"/dir1,/data/dir1,user1,pass1,true;/dir2,/data/dir2,null,null,false"` 描述了 2 个服务，分别是将文件夹 `/data/dir1` 映射至 WebDAV 服务 `/dir1`，将文件夹 `/data/dir2` 映射至 WebDAV 服务 `/dir2`。

对于每个服务，需要使用逗号分隔 5 个参数，分别是 `服务路径,本地路径,用户名,密码,是否只读`。其中用户名和密码都为 `null` 时表示不需要验证。例子：

- `"/dir1,/data/dir1,user1,pass1,true"` 描述了将 `/data/dir1` 映射至 `/dir1` 服务，访问需要的用户名和密码分别为 `user1` 和 `pass1`，只读(禁止上传、更新、删除)。
- `"/dir2,/data/dir2,null,null,false"` 描述了将 `/data/dir2` 映射至 `/dir2` 服务，访问不需要验证，可读写。
- `"/dir3,/data/dir3,null,null,true"` 描述了将 `/data/dir3` 映射至 `/dir3` 服务，访问不需要验证，只读。

特别的，如果只有 1 个服务且名为 `/`，则可以直接访问 <http://localhost:80> 而不需要指定服务名。例子：

- `"/,/data/dir1,user1,pass1,true"` 描述了将 `/data/dir1` 映射至 `/` 服务，访问需要的用户名和密码分别为 `user1` 和 `pass1`，只读。

当 `dav` 未指定时，GoWebDAV 默认使用的 `dav` 参数为 `/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false`。

## Docker

准备分享的本地文件夹路径为 `/data/dir1` 和 `/data/dir2`。

```sh
docker run -it -d -v /data:/data -e dav="/dir1,/data/dir1,user1,pass1,true;/dir2,/data/dir2,null,null,false" -p 80:80 --restart=unless-stopped 117503445/go_webdav
```

在浏览器中打开 <http://localhost/dir1> 和 <http://localhost/dir2>，就能以 WebDAV 的形式访问磁盘文件了。

通过环境变量 `dav` 传递 `data` 参数，通过 `-p 80:80` 指定映射的端口。

## Docker Compose

```yaml
services:
  go_webdav:
    image: 117503445/go_webdav
    restart: unless-stopped
    volumes:
      - /data:/data
    environment:
      - "dav=/dir1,/data/dir1,user1,pass1,true;/dir2,/data/dir2,null,null,false"
    ports:
      - "80:80"
```

如果需要在非 root 用户的环境下使用 Docker 运行 Gowebdav，可以参考 [nonroot](./doc/nonroot_zh_CN.md)。

## 配置

GoWebDAV 支持通过环境变量、命令行参数、配置文件等方式配置 WebDAV 服务，本文的上述例子是 GowebDAV 的典型使用方式。如果上述例子无法满足你的需求，可以参考 [配置](./doc/config_zh_CN.md)。

## 安全

GoWebDAV 使用 HTTP Basic Auth 进行验证，账号密码未经加密，毫无安全性可言。如果涉及重要文件、重要密码，请务必用 Nginx 或 Traefik 等代理服务器套一层 HTTPS。

GoWebDAV 目前没有直接支持 HTTPS 的计划，因为我认为 HTTPS 涉及域名、证书定期申请，这些工作都应当在上层代理服务器中完成。

## 开发

见 [dev.md](./doc/dev_zh_CN.md)

## 致谢

<https://github.com/dom111/webdav-js> 提供了前端支持
