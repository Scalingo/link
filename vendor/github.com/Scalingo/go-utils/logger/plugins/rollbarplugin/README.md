# Rollbar Plugin

This plugin will send every log with a level >= Error to rollbar.

## Configuration

This plugin need two different environment variables:

* `ROLLBAR_TOKEN`: The rollbar token used to authenticate against the Rollbar API
* `GO_ENV`: The go environment

# Usage

To use it add this to your main:

```
rollbarplugin.Register()
```
