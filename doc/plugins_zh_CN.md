# plugins

通过插件扩展 GoWebDAV 的功能。

## PreRequestHook

PreRequestHook 插件是一个 GoWebDAV 的钩子插件，用于在处理 WebDAV 请求之前执行一些操作。在启用 PreRequestHook 插件后，GoWebDAV 会在处理 WebDAV 请求之前调用插件的 `PreRequest` 方法，并忽略内置的身份验证逻辑。

插件使用 Go 语言编写，可以参考 [example](../assets/Plugins/PreRequestExample.go)。

大致上，你需要实现一个 `PreRequestHook` 方法：

```go
func PreRequest(cfg *gowebdav.HandlerConfig, r *http.Request) *gowebdav.PreRequestResult {
}
```

其中 `cfg` 是匹配到的 WebDAV 服务的配置，`r` 是 HTTP 请求对象。`PreRequest` 方法返回一个 `PreRequestResult` 结构体，用于指示 GoWebDAV 如何处理这个请求。

然后在启动时，将插件路径传入 `--pre_request_hook` 参数即可。

例子:

```sh
./GoWebDAV --pre_request_hook assets/Hooks/PreRequestExample.go --dav /dav1,./dir1,null,null,true;/dav2,./dir2,null,null,true
```
