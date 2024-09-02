# nonroot

如果你希望在非 root 用户的环境下使用 Docker 运行 Gowebdav，可以参考以下步骤。

## 步骤

以 Docker Compose 为例，假设你希望共享 `./data/dir1` 和 `./data/dir2` 两个目录。首先，准备一个 `docker-compose.yml` 文件：

```yaml
services:
  go_webdav:
    image: 117503445/go_webdav
    restart: unless-stopped
    volumes:
      - ./data:/home/nonroot
    environment:
      - "dav=/dir1,/home/nonroot/dir1,null,null,false;/dir2,/home/nonroot/dir2,null,null,false"
    ports:
      - "80:80"
    user: "nonroot" # 指定容器内的用户为 nonroot
```

接下来，创建目录并赋予 777 权限：

```bash
mkdir -p ./data/dir1 ./data/dir2
chmod 777 ./data/dir1 ./data/dir2
```

最后，启动容器：

```bash
docker compose up -d
```

## 注意事项

`117503445/go_webdav` 是基于 [gcr.io/distroless/static-debian12](https://github.com/GoogleContainerTools/distroless) 制作的。镜像中的 `nonroot` 用户是非 root 用户，UID 为 65532，并对 `/home/nonroot` 目录有写入权限。

- 如果你没有提前创建 `data` 目录，容器启动后会自动创建该目录。但由于这是由具有 root 权限的 Docker Daemon 创建的，可能会导致权限问题。
- 如果你没有提前创建 `dir1` 和 `dir2` 目录，容器启动后会由 `GoWebdav` 创建这些目录。但由于它们属于 `nonroot` 用户，所以外部的普通用户无法写入。
- 如果你没有提前赋予 `777` 权限，`GoWebdav` 的 `nonroot` 用户将无法写入这些目录。

在上述情景中，容器内外都是普通用户。如果你只要求容器内是普通用户、外部是 root 用户，或者容器内是 root 用户、外部是普通用户，可能会更简单一些。
