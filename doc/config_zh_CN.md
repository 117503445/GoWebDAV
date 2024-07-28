# 配置

GoWebDAV 能从多种来源读取配置，且具有高度的灵活性。本文中，我会先给出一些使用的例子。如果无法满足你的需求，可以查看 *说明* 章节，了解更具体的配置方式。

## 例子

### 二进制 & 命令行参数

```sh
./gowebdav --address 0.0.0.0 --port 80 --dav "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"
```

### 二进制 & 环境变量

```sh
ADDRESS=0.0.0.0 PORT=80 DAV="/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false" ./gowebdav
```

### 二进制 & 配置文件

下载 [config.toml](https://github.com/117503445/GoWebDAV/releases/latest/download/config.toml)，并将其放置在二进制文件 `gowebdav` 旁边，然后运行：

```sh
./gowebdav
```

### Docker & `-e` 环境变量

```sh
docker run -it -d -v ./data:/workspace/data -e dav="/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false" -p 80:80 --restart=unless-stopped 117503445/go_webdav
```

### Docker & CLI 参数

```sh
docker run --rm 117503445/go_webdav --help # 查看帮助
docker run -it -d -v ./data:/workspace/data -p 80:80 --restart=unless-stopped 117503445/go_webdav --address 0.0.0.0 --port 80 --dav "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"
```

### Docker Compose & 配置文件

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

下载 [config.toml](https://github.com/117503445/GoWebDAV/releases/latest/download/config.toml)，并将其放置在二进制文件 `gowebdav` 旁边，然后运行：

```sh
docker compose up -d
```

## 说明

GoWebDAV 能从多种来源读取配置，优先级从高到低依次是：CLI 参数、配置文件、环境变量、默认值。

其中，配置文件路径的默认值是 `config.toml`。当然，配置文件路径也是可以更改的。通过 CONFIG 环境变量可以指定配置文件路径，如 `CONFIG=/path/to/config.toml ./gowebdav`。CLI 也可以改变配置文件路径，且具有比环境变量更高的优先级，如 `./gowebdav --config /path/to/config.toml`。甚至，你还可以传入多个配置文件路径，用 `,` 分隔，如 `./gowebdav --config /path/to/config1.toml,/path/to/config2.toml`。如果出现了重复的配置项，后面的配置文件会覆盖前面的。如果配置文件路径是相对路径的形式，如 `config.toml`，那么会优先尝试解析为可执行文件同目录下的文件，如果不存在，再尝试解析为当前工作目录下的文件。

对于下列配置项，同时支持从 CLI 参数、配置文件、环境变量、默认值中读取：

| 配置项 | 类型 | 默认值 | 说明 | CLI | 配置文件 | 环境变量 |
| --- | --- | --- | --- | --- | --- | --- |
| `address` | `string` | `0.0.0.0` | 监听地址 | `--address 0.0.0.0` | `address = "0.0.0.0"` | `ADDRESS=0.0.0.0` |
| `port` | `int` | `80` | 监听端口 | `--port 80` | `port = 80` | `PORT=80` |
| `dav` | `string` | `/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false` | WebDAV 服务配置 | `--dav "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"` | `dav = "/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"` | `DAV="/public-writable,./data/public-writable,null,null,false;/public-readonly,./data/public-readonly,null,null,true;/private-writable,./data/private-writable,user1,pass1,false"` |
| `secret_dav_list` | `bool` | `false` | 是否隐藏 WebDAV 服务列表 | `--secret-dav-list` | `secret_dav_list = true` | `SECRET_DAV_LIST=true` |

其中，`dav` 适合在 CLI 和 环境变量中使用，可以用单行的形式简便地配置多个 WebDAV 服务。但是在配置文件中，这种写法的可读性较差。为了提升配置文件的使用体验，你可以在配置文件中使用 `davs` 字段，如下：

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

GoWebDAV 会同时解析 `davs` 和 `dav` 字段，并将解析结果合并。如果定义了 `davs` 字段，且 `dav` 字段是默认值，那么会忽略 `dav` 字段；如果解析出的最终 dav 配置和默认值相同，就会自动创建实例文件夹及文件，方便用户快速上手。

好吧，我承认，GoWebDAV 的配置逻辑比较复杂，且各个配置方式不正交，存在很多特例。但是，这样的设计是为了让各种用户都有良好的使用体验。祝使用愉快！
