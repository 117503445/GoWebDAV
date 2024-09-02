# GoWebDAV

> Share local files using WebDAV, lightweight and easy to use

English | [简体中文](./README_zh_CN.md)

## Features

- Implemented in Golang for high performance.

- Finally compiled into a single binary file, no need for Apache or similar environments, with few dependencies.

- Supports browser access.

- Multiple WebDAV services can be enabled on the same port, each with different mount directories, usernames, and passwords.

- Good Docker support.

## Quick Start

Download the binary file from <https://github.com/117503445/GoWebDAV/releases>

Run

```sh
./gowebdav
```

GoWebDAV will automatically create sample files under the `./data` path, with the following file structure

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

Access <http://localhost:80> in your browser to see the 3 different GoWebDAV services.

![index](./doc/index.png)

Among them, <http://localhost:80/public-writable> is the `public-writable` service, mapping to the local `./data/public-writable` folder. It is unauthenticated and writable. You can view file contents in the browser and perform operations like upload and delete.

![public-writable](./doc/public-writable.png)

<http://localhost:80/public-readonly> is the `public-readonly` service, mapping to the local `./data/public-readonly` folder. It is unauthenticated and read-only. You can view file contents in the browser but cannot upload, delete, etc.

![public-readonly](./doc/public-readonly.png)

<http://localhost:80/private-writable> is the `private-writable` service, mapping to the local `./data/private-writable` folder. It requires user authentication and is writable. After logging in with `user1` and `pass1`, you can view file contents in the browser and perform operations like upload and delete.

![private-writable](./doc/private-writable.png)

Besides using a browser, you can also access it using other WebDAV client tools.

You can configure the local path, user authentication, read-only status, and other properties of the WebDAV service by specifying the `dav` parameter. For details, see the *Usage* section.

## Usage

```sh
./gowebdav --help # View help

./gowebdav --addr 127.0.0.1 # Listen on 127.0.0.1, default is 0.0.0.0
./gowebdav --port 8080 # Listen on port 8080, default is port 80

./gowebdav --dav "/dir1,/data/dir1,user1,pass1,true" # Configure folder path and properties
```

The `dav` parameter can specify the local path, user authentication, read-only status, and other properties of the WebDAV service.

Each local path can be configured for a WebDAV service, separated by semicolons. For example:

- `"/dir1,/data/dir1,user1,pass1,true;/dir2,/data/dir2,null,null,false"` describes 2 services, mapping the folder `/data/dir1` to the WebDAV service `/dir1` and the folder `/data/dir2` to the WebDAV service `/dir2`.

For each service, you need to separate 5 parameters with commas: `service path, local path, username, password, read-only status`. When both the username and password are `null`, no authentication is required. For example:

- `"/dir1,/data/dir1,user1,pass1,true"` describes mapping `/data/dir1` to the `/dir1` service, where access requires the username and password `user1` and `pass1`, respectively, and is read-only (prohibits upload, update, delete).
- `"/dir2,/data/dir2,null,null,false"` describes mapping `/data/dir2` to the `/dir2` service, where no authentication is required and it is read-write.
- `"/dir3,/data/dir3,null,null,true"` describes mapping `/data/dir3` to the `/dir3` service, where no authentication is required and it is read-only.

In particular, if there is only one service named `/`, you can access <http://localhost:80> directly without specifying a service name. For example:

- `"/,/data/dir1,user1,pass1,true"` describes mapping `/data/dir1` to the `/` service, where access requires the username and password `user1` and `pass1`, respectively, and is read-only.

When `dav` is not specified, the default `dav` parameter used by GoWebDAV is `/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false`.

## Docker

Prepare the local folder paths to be shared as `/data/dir1` and `/data/dir2`.

```sh
docker run -it -d -v /data:/data -e dav="/dir1,/data/dir1,user1,pass1,true;/dir2,/data/dir2,null,null,false" -p 80:80 --restart=unless-stopped 117503445/go_webdav
```

Open <http://localhost/dir1> and <http://localhost/dir2> in your browser to access disk files in WebDAV format.

Pass the `data` parameter through the environment variable `dav` and specify the mapped port with `-p 80:80`.

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

If you need to run Gowebdav using Docker in a non-root user environment, you can refer to [nonroot](./doc/nonroot.md).

## Configuration

GoWebDAV supports configuring the WebDAV service through environment variables, command-line arguments, configuration files, and other methods. The examples mentioned above are typical uses of GoWebDAV. If these examples do not meet your needs, you can refer to the [Configuration](./doc/config.md) documentation.

## Security

GoWebDAV uses HTTP Basic Auth for authentication, with account passwords transmitted in plaintext, lacking security. If dealing with important files or passwords, be sure to use a layer of HTTPS with Nginx or Traefik proxy servers.

GoWebDAV currently does not have plans to directly support HTTPS, as I believe that tasks like domain names and certificate renewal should be handled at the higher-level proxy server.

## Development

Refer to [dev.md](./doc/dev.md)

## Acknowledgements

<https://github.com/dom111/webdav-js> provides frontend support