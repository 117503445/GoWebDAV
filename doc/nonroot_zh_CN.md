# nonroot

如果你希望在非 root 用户的环境下使用 Docker 运行 Gowebdav，可以参考以下步骤。

## 步骤

以 Docker Compose 为例，假设你希望共享 `./data/dir1` 和 `./data/dir2` 两个目录。首先，准备好您的目录

```bash
mkdir -p ./data # 这里只是一个示例，您可以用任何方式创建目录
```

接下来，获取目录的 UID 和 GID

```bash
ls -nd ./data | awk '{ print $3":"$4 }'
```

然后，创建一个 `docker-compose.yml` 文件：

```yaml
services:
  go_webdav:
    image: 117503445/go_webdav
    restart: unless-stopped
    volumes:
      - ./data:/data
    environment:
      - "dav=/dir1,/data/dir1,null,null,false;/dir2,/data/dir2,null,null,false"
    ports:
      - "80:80"
    user: "1000:1000" # 填写正确的UID和GID以确保以正确的用户执行
```

最后，启动容器：

```bash
docker compose up -d
```

## 注意事项

docker 支持通过`--user "UID:GID"`的方式指定用户，所以我们可以用这个功能让容器运行在非 root 用户下  
不过您仍需提前创建 `data` 目录，以防止 Docker Daemon 使用 root 权限创建，导致权限问题。  
在上述情景中，容器内外都是普通用户。如果你只要求容器内是普通用户、外部是 root 用户，或者容器内是 root 用户、外部是普通用户，可能会更简单一些。
