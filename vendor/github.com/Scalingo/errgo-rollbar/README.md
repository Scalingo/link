## Rollbar - Errgo binding v0.2.0

Rollbar Go client: `https://github.com/stvp/rollbar`
Errgo: `https://github.com/go-errgo/errgo`

### Errgo masking system

Errgo is based on a error masking system which gather information
each time an error is masked and returned to the called.

From this stack of masked error, we can build the calling path of
the error and get the precise origin of a given error.

### Example

```go
func A() error {
  return errgo.Mask(err)
}

func B() error {
  return errgo.Mask(A())
}

func main() {
  err := B()
  rollbar.ErrorWithStack(rollbar.ERR, err, errgorollbar.BuildStack(err))
  rollbar.Wait()
}
```

Then your are able to see on your rollbal dashboard that the correct
stack trace has been sent.

## Release a New Version

Bump new version number in:

- `CHANGELOG.md`
- `README.md`

Commit, tag and create a new release:

```sh
git add CHANGELOG.md README.md
git commit -m "Bump v0.2.0"
git tag v0.2.0
git push origin master
git push --tags
```
