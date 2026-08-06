package git

import (
	"fmt"
	"io"
	"os/exec"
)

func (g *Git) GenPatch(path string) (string, error) {
	cmd := exec.Command("git", "diff")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		return "", fmt.Errorf("generate patch error: run cmd error: %v, cmd output: %v", err, out)
	}

	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	patch, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(patch), nil
}

func (g *Git) ApplyPatch(path, file string) (string, error) {
	cmd := exec.Command("git", "apply", file)
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		out, _ := cmd.CombinedOutput()
		return "", fmt.Errorf("apply patch error: run cmd error: %v, cmd output: %v", err, out)
	}

	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	patch, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(patch), nil
}
