package upgrades

import (
	"testing"

	"github.com/Scalingo/link/v3/tests/integration/utils"
)

func TestMain(m *testing.M) {
	stopper, err := utils.StartEtcd("3.5.17")
	if err != nil {
		panic(err)
	}
	defer stopper()

	m.Run()
}
