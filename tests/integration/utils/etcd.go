package utils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/link/v3/models"
)

func CleanupEtcdData(t *testing.T) {
	t.Helper()

	client, err := etcd.ClientFromEnv()
	require.NoError(t, err)

	dir := models.EtcdLinkDirectory

	// t.Context() has already been canceled since we're in a cleanup function.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Logf("Removing directory %s from etcd", dir)
	_, err = client.KV.Delete(ctx, dir, etcdv3.WithPrefix())
	require.NoError(t, err)
}

func StartEtcd(version string) (func(), error) {
	fmt.Printf("Downloading etcd version %s\n", version)
	etcdBinary, err := downloadEtcd(version)
	if err != nil {
		return nil, fmt.Errorf("download etcd: %w", err)
	}

	dataDir := filepath.Join(os.TempDir(), "etcd-data"+strconv.Itoa(int(time.Now().UnixMilli())))
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("create etcd data directory: %w", err)
	}

	cmd := exec.Command(etcdBinary,
		"--name", "etcd-test",
		"--data-dir", dataDir,
		"--listen-client-urls", "http://0.0.0.0:2379",
		"--listen-peer-urls", "http://0.0.0.0:2380",
		"--advertise-client-urls", "http://localhost:2379",
	)

	fmt.Printf("Starting etcd with command: %s\n", cmd.String())
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("start etcd: %w", err)
	}

	waitTimeout := 10 * time.Second
	startTime := time.Now()
	fmt.Println("Waiting for etcd to be ready...")
	for {
		if time.Since(startTime) > waitTimeout {
			return nil, fmt.Errorf("etcd did not start within %s", waitTimeout)
		}
		conn, err := net.DialTimeout("tcp", "localhost:2379", 5*time.Second)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("etcd is ready")

	return func() {
		fmt.Printf("Stopping etcd process with PID %d\n", cmd.Process.Pid)
		cmd.Process.Kill()
		fmt.Printf("Removing etcd data directory: %s\n", dataDir)
		os.RemoveAll(dataDir)
		fmt.Printf("Removing etcd binary: %s\n", etcdBinary)
		os.Remove(etcdBinary)
	}, nil
}

func downloadEtcd(version string) (string, error) {
	_, err := os.Open("/tmp/etcd-" + version)
	if err == nil {
		return "/tmp/etcd-" + version, nil
	}

	url := "https://github.com/etcd-io/etcd/releases/download/v" + version + "/etcd-v" + version + "-linux-amd64.tar.gz"
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	gReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	defer gReader.Close()

	etcdBinary := os.TempDir()
	err = os.MkdirAll(etcdBinary, 0755)
	if err != nil {
		return "", err
	}

	tarReader := tar.NewReader(gReader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break // end of tar archive
			}
			return "", err
		}

		if header.Typeflag == tar.TypeReg {

			path := filepath.Base(header.Name)
			if path != "etcd" {
				continue // skip non-etcd files
			}

			file, err := os.Create(filepath.Join(etcdBinary, "etcd-"+version))
			if err != nil {
				return "", err
			}
			defer file.Close()

			_, err = io.Copy(file, tarReader)
			if err != nil {
				return "", err
			}

			err = file.Chmod(0700) // make it executable
			if err != nil {
				return "", err
			}

			return "/tmp/etcd-" + version, nil
		}
	}
	return "", fmt.Errorf("etcd binary not found in the tar archive")
}
