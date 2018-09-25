package githubprobe

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGithubProbe(t *testing.T) {
	Convey("When github respond healthy", t, func() {
		response := GithubStatusResponse{
			Status:        "nothing",
			LastUpdatedAt: time.Now(),
		}

		checker := GithubChecker{}

		buffer := new(bytes.Buffer)

		err := json.NewEncoder(buffer).Encode(&response)
		So(err, ShouldBeNil)

		err = checker.Check(buffer)
		So(err, ShouldBeNil)
	})

	Convey("When github respond not healthy", t, func() {
		response := GithubStatusResponse{
			Status:        "major",
			LastUpdatedAt: time.Now(),
		}

		checker := GithubChecker{}
		buffer := new(bytes.Buffer)
		err := json.NewEncoder(buffer).Encode(&response)
		So(err, ShouldBeNil)

		err = checker.Check(buffer)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Github is probably down")
	})
}
