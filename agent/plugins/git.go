package plugins

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/lodastack/agent/agent/common"
)

// Update updates plugin
func Update(namespace, gitPath string, pull bool) error {
	pluginDir := path.Join(common.Conf.PluginsDir, namespace)
	if !common.Exists(pluginDir) {
		if err := os.MkdirAll(pluginDir, 0755); err != nil {
			return err
		}
	}
	s := strings.Split(gitPath, "/")
	last := s[len(s)-1]
	pluginName := last[:len(last)-4]
	dir := path.Join(pluginDir, pluginName)
	if !common.Exists(dir) {
		cmd := exec.Command("git", "clone", gitPath)
		cmd.Dir = pluginDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("clone failed: %s", err)
		}
		if err := checkoutRelease(dir); err != nil {
			return fmt.Errorf("can not checkout to release: %s", err)
		}
		return nil
	} else {
		if err := checkBranch(dir); err != nil {
			if err = updateBranches(dir); err != nil {
				return fmt.Errorf("can not update remote branches: %s", err)
			}
			if err = checkoutRelease(dir); err != nil {
				return fmt.Errorf("can not checkout to release: %s", err)
			}
		}
		if !pull {
			return nil
		}
		cmd := exec.Command("git", "pull", "origin", "release")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to pull from release: %s", err)
		}
		return nil
	}
}

func updateBranches(dir string) error {
	cmd := exec.Command("git", "remote", "update", "origin", "--prune")
	cmd.Dir = dir
	return cmd.Run()
}

func checkoutRelease(dir string) error {
	cmd := exec.Command("git", "checkout", "release")
	cmd.Dir = dir
	return cmd.Run()
}

func checkBranch(dir string) error {
	cmd := exec.Command("git", "branch")
	cmd.Dir = dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return err
	} else {
		s := string(stdout.Bytes())
		if !strings.Contains(s, "* release") {
			return fmt.Errorf("%s", "not on branch release")
		} else {
			return nil
		}
	}
}
