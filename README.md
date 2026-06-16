# hotconf

Watch a config file and reload it without restarting your process.

```go
watcher, _ := hotconf.NewWatcher(200 * time.Millisecond)

watcher.Watch("/etc/app/config.json", func(path string) {
    cfg, _ := hotconf.Load[Config](path, json.Unmarshal)
    appConfig.Set(cfg)
})

watcher.Start(ctx)
defer watcher.Stop()
```

## Install

```sh
go get github.com/ntentasd/hotconf
```

## How it works

`hotconf` wraps [fsnotify](https://github.com/fsnotify/fsnotify) and adds debouncing so rapid-fire write events (common with editors and the kubelet) don't trigger redundant reloads. `Config[T]` is an `RWMutex`-protected holder. It is safe to be read concurrently while the watcher writes to it.

Works with any format; pass your own unmarshal function.

## Kubernetes

Mount your ConfigMap as a volume. The kubelet syncs changes to the file on disk; `hotconf` picks them up without a pod restart.

```yaml
volumeMounts:
  - name: config
    mountPath: /etc/app
---
volumes:
  - name: config
    configMap:
      name: app-config
```
