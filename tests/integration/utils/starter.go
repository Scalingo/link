package utils

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var verboseLinkLogging *bool = flag.Bool("link.verbose", false, "Enable verbose logging for LinK processes")

type LinKProcess struct {
	binaryPath string
	port       int
	cmd        *exec.Cmd
}

type StartLinkOpt func(cmd *exec.Cmd)

func WithEnv(key, value string) StartLinkOpt {
	return func(cmd *exec.Cmd) {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
}

func StartLinK(t *testing.T, binaryPath string, opts ...StartLinkOpt) LinKProcess {
	t.Helper()

	port := rand.Intn(10000) + 10000 // Random port between 10000 and 20000

	t.Logf("Starting LinK process with binary %s on port %d", binaryPath, port)

	// Start the LinK process here in another goroutine
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(binaryPath)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", port))
	cmd.Env = append(cmd.Env, "ETCD_HOSTS=http://localhost:2379")
	cmd.Env = append(cmd.Env, "INTERFACE=eth0")
	cmd.Env = append(cmd.Env, "HOSTNAME=test-host")
	cmd.Env = append(cmd.Env, "KEEPALIVE_INTERVAL=100ms")
	cmd.Env = append(cmd.Env, "GO_ENV=test")

	if verboseLinkLogging == nil || !*verboseLinkLogging {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	} else {
		cmd.Env = append(cmd.Env, "LOGGER_LEVEL=debug")
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	}

	for _, opt := range opts {
		opt(cmd)
	}

	err := cmd.Start()
	require.NoError(t, err)

	processExited := make(chan any, 1)

	go func() {
		cmd.Wait()
		processExited <- struct{}{}
		close(processExited)
	}()

	res := LinKProcess{
		binaryPath: binaryPath,
		port:       port,
		cmd:        cmd,
	}

	t.Logf("LinK process started with PID %d, waiting for it to bind to port %d", cmd.Process.Pid, port)
	t.Cleanup(func() {
		t.Log("Cleaning up LinK process")
		res.Stop(t)
	})

	// wait for the process to bind the port

	bindTimeout := 30 * time.Second
	startTime := time.Now()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-processExited:
			stdout := stdout.String()
			stderr := stderr.String()
			t.Fatalf("LinK process with PID %d exited unexpectedly. Stdout:\n%s\nStderr: \n%s", cmd.Process.Pid, stdout, stderr)
			return LinKProcess{}
		}

		if time.Since(startTime) > bindTimeout {
			t.Fatalf("LinK process with PID %d did not bind to port %d within %s timeout, stdout:\n%s\nStderr:\n%s", cmd.Process.Pid, port, bindTimeout, stdout.String(), stderr.String())
		}

		time.Sleep(100 * time.Millisecond)

		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			conn.Close()
			break // port is bound
		}
	}

	t.Logf("LinK process with PID %d is bound to port %d", cmd.Process.Pid, port)

	// wait for LinK to initialize
	time.Sleep(300 * time.Millisecond)

	return res
}

func (p LinKProcess) URL() string {
	return fmt.Sprintf("localhost:%d", p.port)
}

func (p LinKProcess) Stop(t *testing.T) {
	t.Helper()

	t.Logf("Stopping LinK process with PID %d", p.cmd.Process.Pid)

	err := p.cmd.Process.Kill()
	if err != os.ErrProcessDone {
		require.NoError(t, err)
	}
	t.Logf("LinK process with PID %d stopped", p.cmd.Process.Pid)
}
