# Really Simpler Logger
------------

### Package
```
go get github.com/saarwasserman/really-simple-logger@1.0.0
```

### Example
```
logger := reallysimplelogger.New(os.Stdout, reallysimplelogger.LevelInfo)

logger.Info("This is a call to space", map[string]string{"method": "post"})
logger.Error(errors.New("got an error"), nil)
```
