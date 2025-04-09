# Rollbar Plugin

This plugin will send every log with a level >= Error to rollbar.

## Configuration

This plugin needs two different environment variables:

* `ROLLBAR_TOKEN`: The rollbar token used to authenticate against the Rollbar API
* `GO_ENV`: The go environment

# Usage

To use this plugin, register it during the initialization of your program.

/!\ This line MUST be written BEFORE the first initialization of the logger.

```go
rollbarplugin.Register()
defer rollbarplugin.Close()

log := logger.Default()
ctx := logger.ToCtx(context.Background(), log)
```
