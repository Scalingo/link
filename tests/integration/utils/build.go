package utils

import (
	"os"
	"os/exec"
	"testing"
)

func BuildLinKBinary(t *testing.T) string {
	t.Helper()

	destPath := t.TempDir()

	codePath := "/go/src/github.com/Scalingo/link"
	if os.Getenv("CODE_PATH") != "" {
		codePath = os.Getenv("CODE_PATH")
	}

	t.Logf("Building LinK binary from %s to %s", codePath, destPath)

	cmd := exec.Command("go", "build", "-buildvcs=false", "-o", destPath+"/link", codePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build LinK binary: %s\nOutput: %s", err.Error(), string(output))
	}

	return destPath + "/link"
}
