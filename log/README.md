log
===

The `log` package is a collection of context utilities to make working with zap.Logger better.

[Docs](https://pkg.go.dev/github.com/drshriveer/gtools/log)

### Getting started

Install with:

```bash
go get -u github.com/drshriveer/gtools/log
```

### Features

-	Log field propagation (on primed contexts when fields are added added downstream before logging upstream)

### TODO

-	Deferred logging (so that fields can propagate AFTER a log.XXX() has been called)
-	Info/Debug capture and print IFF error encountered or flagged
