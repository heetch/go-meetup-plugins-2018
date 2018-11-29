package demo2

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setFakeExecCommand(t *testing.T) func() {
	execCommand = func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_RUN_HELPER_PROCESS=1"}
		return cmd
	}

	return func() {
		execCommand = exec.Command
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_RUN_HELPER_PROCESS") != "1" {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	require.Equal(t, "/fake/command", cmd)
	require.Len(t, args, 2)
	require.Equal(t, "--player-addr", args[0])

	<-c
	os.Exit(0)
}

func TestPlugin(t *testing.T) {
	cleanup := setFakeExecCommand(t)
	defer cleanup()

	p := NewPlugin("/fake/command", "127.0.0.1:5000")
	err := p.Start()
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	err = p.Stop()
	require.NoError(t, err)
}

func TestListPlugins(t *testing.T) {
	dir, err := ioutil.TempDir("", "player")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	err = ioutil.WriteFile(filepath.Join(dir, "player-plugin1"), nil, 0600)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir, "player-plugin2"), nil, 0600)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(dir, "plugin3"), nil, 0600)
	require.NoError(t, err)

	list, err := listPlugins(dir)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, filepath.Join(dir, "player-plugin1"), list[0])
	require.Equal(t, filepath.Join(dir, "player-plugin2"), list[1])
}
