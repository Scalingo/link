# Rollbar Hook for [Logrus](https://github.com/sirupsen/logrus) v1.4.2

## Setup

```sh
go get github.com/Scalingo/logrus-rollbar
```

## Example

```go
package main

import (
	"fmt"
	"net/http"
	"math/rand"

	"github.com/sirupsen/logrus"
	"github.com/stvp/rollbar"

	logrusrollbar "github.com/Scalingo/logrus-rollbar"
)


func main() {
	rollbar.ApiKey = "123456ABCD"
	rollbar.Environment = "testing"

	logger := logrus.New()
	logger.Hooks.Add(logrusrollbar.Hook{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("something wrong happened in the database")

		logger.WithFields(
			logrus.Fields{"req": r, "error": err, "extra-data", rand.Int()},
		).Error("Something is really wrong")
	})

	http.ListenAndServe(":31313", nil)
}
```


## Release a New Version

Bump new version number in:

- `CHANGELOG.md`
- `README.md`

Commit, tag and create a new release:

```sh
version="1.4.2"

git switch --create release/${version}
git add CHANGELOG.md README.md
git commit --message="Bump v${version}"
git push --set-upstream origin release/${version}
gh pr create --reviewer=EtienneM --title "$(git log -1 --pretty=%B)"
```

Once the pull request merged, you can tag the new release.

```sh
git tag v${version}
git push origin master v${version}
gh release create v${version}
```

The title of the release should be the version number and the text of the
release is the same as the changelog.
