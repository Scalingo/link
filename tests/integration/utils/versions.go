package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func DownloadLinKVersion(t *testing.T, version string) string {
	t.Helper()

	// https://github.com/Scalingo/link/releases/download/v3.0.2/link-v3.0.2-linux-amd64.tar.gz
	baseDir := t.TempDir()

	t.Logf("Downloading link version %s to %s", version, baseDir)

	resp, err := http.Get(fmt.Sprintf("https://github.com/Scalingo/link/releases/download/v%s/link-v%s-linux-amd64.tar.gz", version, version))
	require.NoError(t, err)
	defer resp.Body.Close()

	// decompress the tar.gz file
	gzipReader, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break // end of tar archive
			}
			require.NoError(t, err)
		}

		if header.Typeflag == tar.TypeReg { // regular file

			filename := filepath.Base(header.Name)
			if filename != "link" {
				continue // skip non-link files
			}

			file, err := os.Create(filepath.Join(baseDir, filename))
			require.NoError(t, err)
			defer file.Close()

			err = file.Chmod(0700) // make it executable
			if err != nil {
				require.NoError(t, err)
			}

			_, err = io.Copy(file, tarReader)
			require.NoError(t, err)

			t.Logf("LinK version %s downloaded and extracted to %s", version, filepath.Join(baseDir, filename))
			return filepath.Join(baseDir, "link")
		}
	}

	t.Fatal("Failed to find 'link' binary in the tar archive")
	return ""
}
