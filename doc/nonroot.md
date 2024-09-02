# nonroot

If you want to use Gowebdav with Docker under a non-root user, you can follow these steps.

## Steps

Using Docker Compose as an example, let's say you want to share `./data/dir1` and `./data/dir2` directories. First, prepare a `docker-compose.yml` file:

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
    user: "nonroot" # Specify the user inside the container as nonroot
```

Next, create the directories and set the permissions to 777:

```bash
mkdir -p ./data/dir1 ./data/dir2
chmod 777 ./data/dir1 ./data/dir2
```

Finally, start the container:

```bash
docker compose up -d
```

## Notes

`117503445/go_webdav` is based on [gcr.io/distroless/static-debian12](https://github.com/GoogleContainerTools/distroless). The `nonroot` user inside the image is a non-root user with a UID of 65532 and has write permissions to the `/home/nonroot` directory.

- If you do not create the `data` directory in advance, it will be automatically created when the container starts. However, this will be done by the Docker Daemon with root privileges, which may lead to permission issues.
- If you do not create the `dir1` and `dir2` directories in advance, they will be created by `GoWebdav` when the container starts. Since these directories will belong to the `nonroot` user, external regular users will not be able to write to them.
- If you do not set the permissions to 777 in advance, the `nonroot` user in `GoWebdav` will not be able to write to these directories.

In the scenarios described above, both inside and outside the container are regular users. If you only require the container to run as a regular user and the host to run as root, or vice versa, the setup might be simpler.
