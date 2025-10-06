# nonroot

If you wish to run Gowebdav in Docker under a non-root user environment, please follow the steps below.

## Steps

Taking Docker Compose as an example, suppose you want to share two directories: `./data/dir1` and `./data/dir2`. First, prepare your directories:

```bash
mkdir -p ./data  # This is just an example; you can create directories in any way you prefer
```

Next, obtain the UID and GID of the directory:

```bash
ls -nd ./data | awk '{ print $3":"$4 }'
```

Then, create a `docker-compose.yml` file:

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
    user: "1000:1000"  # Replace with the correct UID and GID to ensure execution under the proper user
```

Finally, start the container:

```bash
docker compose up -d
```

## Notes

Docker supports specifying a user via the `--user "UID:GID"` option, allowing you to run containers as a non-root user.  
However, you must create the `data` directory in advance to prevent the Docker daemon from creating it with root permissions, which could lead to permission issues.  
In the scenario above, both inside and outside the container use a regular (non-root) user. If you only require the container to run as a regular user while the host uses root (or vice versa), the setup might be simpler.
