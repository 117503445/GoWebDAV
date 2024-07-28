# Configuration

GoWebDAV can read configurations from multiple sources and is highly flexible. In this article, I will first provide some usage examples. If these do not meet your needs, you can check the *Description* section for more specific ways to configure.

## Examples

### Binary & Command Line Arguments

```sh
./gowebdav --address 0.0.0.0 --port 80 --dav "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"
```

### Binary & Environment Variables

```sh
ADDRESS=0.0.0.0 PORT=80 DAV="/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false" ./gowebdav
```

### Binary & Configuration File

Download [config.toml](https://github.com/117503445/GoWebDAV/releases/latest/download/config.toml) beside the binary file `gowebdav`, and then run:

```sh
./gowebdav
```

### Docker & `-e` Environment Variables

```sh
docker run -it -d -v ./data:/workspace/data -e dav="/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false" -p 80:80 --restart=unless-stopped 117503445/go_webdav
```

### Docker & CLI Arguments

```sh
docker run --rm 117503445/go_webdav --help # View help
docker run -it -d -v ./data:/workspace/data -p 80:80 --restart=unless-stopped 117503445/go_webdav --address 0.0.0.0 --port 80 --dav "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"
```

### Docker Compose & Configuration File

```yaml
# docker-compose.yml
services:
  go_webdav:
    image: 117503445/go_webdav
    restart: unless-stopped
    volumes:
      - ./config.toml:/workspace/config.toml
      - /data:/data
    ports:
      - "80:80"
```

Download [config.toml](https://github.com/117503445/GoWebDAV/releases/latest/download/config.toml) beside the `docker-compose.yml`, and then run:

```sh
docker compose up -d
```

## Description

GoWebDAV can read configurations from various sources, with a priority order from high to low as follows: CLI arguments, configuration files, environment variables, default values.

The default path for the configuration file is `config.toml`. Of course, this can be changed. You can specify the path to the configuration file through the CONFIG environment variable, such as `CONFIG=/path/to/config.toml ./gowebdav`. The CLI can also change the configuration file path, with a higher priority than environment variables, such as `./gowebdav --config /path/to/config.toml`. Moreover, you can pass multiple configuration file paths separated by a comma, such as `./gowebdav --config /path/to/config1.toml,/path/to/config2.toml`. If there are duplicate configuration items, the later configuration file will override the earlier ones. If the configuration file path is a relative path, such as `config.toml`, it will first attempt to resolve it as a file in the same directory as the executable; if it does not exist, it will then try to resolve it as a file in the current working directory.

The following configuration items support reading from CLI arguments, configuration files, environment variables, and default values:

| Configuration Item | Type | Default Value | Description | CLI | Configuration File | Environment Variable |
| --- | --- | --- | --- | --- | --- | --- |
| `address` | `string` | `0.0.0.0` | Listening address | `--address 0.0.0.0` | `address = "0.0.0.0"` | `ADDRESS=0.0.0.0` |
| `port` | `int` | `80` | Listening port | `--port 80` | `port = 80` | `PORT=80` |
| `dav` | `string` | `/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false` | WebDAV service configuration | `--dav "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"` | `dav = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"` | `DAV="/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"` |
| `secret_dav_list` | `bool` | `false` | Whether to hide the WebDAV service list | `--secret-dav-list` | `secret_dav_list = true` | `SECRET_DAV_LIST=true` |

`dav` is suitable for use in CLI and environment variables because it allows you to conveniently configure multiple WebDAV services in a single line. However, this format is less readable in configuration files. To enhance the usability of configuration files, you can use the `davs` field as follows:

```toml
# each [[davs]] block serves a folder to a sub WebDAV server
[[davs]]
prefix = "/public-writable" # URL prefix, visit http://localhost/public-writable to access in browser, or use WebDAV client
pathDir = "./data/public-writable" # local folder to serve
# username and password for basic auth, set both to "null" to disable auth
username = "null"
password = "null"
readOnly = false # whether to allow write operations(POST, PUT, DELETE)

[[davs]]
prefix = "/public-readonly"
pathDir = "./data/public-readonly"
username = "null"
password = "null"
readOnly = true

[[davs]]
prefix = "/private-writable"
pathDir = "./data/private-writable"
username = "user1"
password = "pass1"
readOnly = false
```

GoWebDAV will parse both the `davs` and `dav` fields and merge the results. If the `davs` field is defined and the `dav` field is the default value, then the `dav` field will be ignored; if the final dav configuration parsed is the same as the default, instance folders and files will be automatically created to facilitate quick start-up.

Admittedly, the configuration logic of GoWebDAV is somewhat complex, and the various configuration methods are not orthogonal, with many special cases. However, this design aims to provide a good user experience for all types of users. Enjoy using it!
