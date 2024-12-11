# Plugins

Extend the functionality of GoWebDAV through plugins.

## PreRequestHook

The PreRequestHook plugin is a hook plugin for GoWebDAV that allows you to perform certain operations before processing WebDAV requests. When the PreRequestHook plugin is enabled, GoWebDAV will call the plugin's `PreRequest` method before handling WebDAV requests, bypassing the built-in authentication logic.

The plugin is written in Go, and you can refer to the [example](../assets/Plugins/PreRequestExample).

Generally, you need to implement a `PreRequestHook` method:

```go
func PreRequest(cfg *gowebdav.HandlerConfig, r *http.Request) *gowebdav.PreRequestResult {
}
```

Here, `cfg` is the configuration for the matched WebDAV service, and `r` is the HTTP request object. The `PreRequest` method returns a `PreRequestResult` struct, which indicates how GoWebDAV should handle the request.

Then, at startup, pass the plugin path using the `--pre_request_hook` parameter.

Example:

```sh
./GoWebDAV --pre_request_hook assets/Hooks/PreRequestExample --dav /dav1,./dir1,null,null,true;/dav2,./dir2,null,null,true
```
