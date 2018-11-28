package demo1

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/pkg/errors"
)

// PluginLoader takes care of scanning a directory for plugins and loading them.
type PluginLoader struct {
	Path       string
	PlayerAddr string
	Plugins    []*Plugin
}

// Load all plugins found in Path. If one of the plugins fails to load, the loading stops and all the
// already loaded plugins are stopped.
func (p *PluginLoader) Load() error {
	list, err := listPlugins(p.Path)
	if err != nil {
		return err
	}

	for _, name := range list {
		plg := NewPlugin(name, p.PlayerAddr)
		err = plg.Start()
		if err != nil {
			break
		}

		p.Plugins = append(p.Plugins, plg)
	}

	if err != nil {
		for _, plg := range p.Plugins {
			err := plg.Stop()
			if err != nil {
				log.Printf("Plugin %s failed to stop correctly: '%s'\n", plg.Path, err)
			}
		}

		return err
	}

	return nil
}

func (p *PluginLoader) Stop() error {
	for _, plg := range p.Plugins {
		err := plg.Stop()
		if err != nil {
			log.Printf("Plugin %s failed to stop correctly: '%s'\n", plg.Path, err)
		}
	}

	return nil
}

func listPlugins(dir string) ([]string, error) {
	var list []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dir {
			return filepath.SkipDir
		}

		if strings.HasPrefix(info.Name(), "player-") {
			list = append(list, path)
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list plugins at path '%s'\n", dir)
	}

	return list, nil
}

var execCommand = exec.Command

// Plugin manages a command lifecycle.
type Plugin struct {
	Path       string
	PlayerAddr string
	cmd        *exec.Cmd
	mu         sync.Mutex
}

// NewPlugin creates a plugin.
func NewPlugin(path, playerAddr string) *Plugin {
	return &Plugin{Path: path, PlayerAddr: playerAddr}
}

// Start the command at Path. The command is started in its own group. The --player-addr flag is
// passed and must be parsed during startup. This method is non-blocking, to stop the plugin use the
// Stop method.
func (p *Plugin) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil {
		return errors.New("plugin already started")
	}

	// pass the player address to the plugin using a command line flag
	args := []string{
		"--player-addr", p.PlayerAddr,
	}

	p.cmd = execCommand(p.Path, args...)
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr
	// make sure the child process starts in a different process group.
	// this prevents the OS to send CTRL-C to the child process and allows
	// the parent process to control when the child processes must stop.
	p.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	return errors.Wrap(p.cmd.Start(), "failed to start plugin")
}

// Stop sends the Interrupt signal to the command and waits for completion.
func (p *Plugin) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return errors.Wrap(err, "failed to send interrupt signal to plugin")
	}

	_ = p.cmd.Wait()
	return nil
}
