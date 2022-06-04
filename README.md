# GoWebdav

> A Simple, Powerful WebDAV Server By Golang.

[中文](./README_CN.md)

## feature

- Based on Golang implementation, high performance

- Finally compiled into a single binary file, no need for Apache and other environments, more stable

- Support browser access

- Multiple WebDAV services can be enabled under the same port, each with a different mount directory, user name, and password

- Docker is well supported

## quickstart

### bin

Go to <https://github.com/117503445/GoWebDAV/releases> to download the latest binaries.

Then run `./gowebdav "/dav1,/root/dir1,user1,pass1,true;/dav2,/root/dir2,null,null,false"`

### Docker

```sh
docker run -it --name go_webdav -d -v /root/dir1:/root/dir1 -v /root/dir2:/root/dir2 -e dav="/dav1,/root/dir1,user1,pass1,true;/dav2,/root/dir2,null,null,false" -p 80:80 --restart=unless-stopped 117503445/go_webdav
```

```sh
-e dav="/dav1,/root/dir1,user1,pass1,true;/dav2,/root/dir2,null,null,false"
```

Indicates passing a configuration string into the Docker image.

Then open <http://localhost/dav1> and <http://localhost/dav2> in the browser or webdav client like [raidrive](https://www.raidrive.com/).

## Configuration String

Use a semicolon to separate each WebDAV service configuration, which means that the above string describes 2 services, which are

> /dav1,/root/dir1,user1,pass1,false

and

> /dav2,/root/dir2,null,null,true

Use a semicolon to separate each WebDAV service configuration, which means that the above string describes 2 services, which are

> /dav1,/root/dir1,user1,pass1,false

and

> /dav2,/root/dir2,null,null,true

The first service will mount the `/root/dir1` directory of the Docker image under `/dav1`. The required username and password for access are `user1` and `pass1` respectively.

Then, according to the previous `-v /root/dir1:/root/dir1`, the mapping relationship with `/root/dir1` of the physical machine can be completed and accessed.

The fifth parameter `false` indicates that this is a non-read-only service that supports addition, deletion, modification and query.

The second service will mount the `/root/dir2` directory of the Docker image under `/dav2`. The user name and password required for access are `null` and `null` respectively. At this time, it means that the service can be accessed without a password. .

Then according to the previous `-v /root/dir2:/root/dir2`, you can complete the mapping relationship with `/root/dir2` of the physical machine and access it.

The fifth parameter `true` indicates that this is a read-only service, only supports GET, does not support additions, deletions and modifications.

This method is recommended for file sharing without confidentiality requirements.

Note that the first argument cannot be `/static`.

## Background introduction

`GoWebdav` is used to build a WebDAV-based file sharing server.

### Reasons to use WebDAV

1. Samba is inconvenient to use on Windows clients, and it is difficult to use non-default ports.

2. FTP mount trouble.

3. NextCloud is too heavy and difficult to share files on the server.

### Reasons to reinvent the wheel of a WebDAV Server

I haven't seen a server implementation that can meet the above characteristics.

## local debugging

Rename `config.yml.example` to `config.yml`, configure in `config.yml` file

`go run .`

## Local Docker run

Using a layered build, the executable app is built through `go build` in the build layer, and then run in the prod layer. If you need to modify the structure of the configuration file later, you will also need to modify the Dockerfile.

````sh
docker build -t 117503445/go_webdav .
docker run --name go_webdav -d -v ${PWD}/TestDir1:/root/TestDir1 -v ${PWD}/TestDir2:/root/TestDir2 -e dav="/dav1,/root/TestDir1,user1,pass1 ,false;/dav2,/root/TestDir2,user2,pass2,true" -p 80:80 --restart=unless-stopped 117503445/go_webdav
````

## safety

HTTP Basic Auth is used for authentication, and the account password is sent in clear text, which has no security at all. If important files or passwords are involved, be sure to use a gateway such as Nginx or Traefik to provide HTTPS.

## THANKS

<https://github.com/dom111/webdav-js> provides front-end support