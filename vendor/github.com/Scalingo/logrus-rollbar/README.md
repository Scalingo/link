# rollbar hook for [Logrus](https://github.com/sirupsen/logrus)

## Setup

```sh
go get gopkg.in/Scalingo/logrus-rollbar.v1
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
	logrusrollbar "gopkg.in/Scalingo/logrus-rollbar.v1"
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
