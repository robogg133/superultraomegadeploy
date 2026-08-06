package git_test

import (
	"os"
	"path/filepath"
	"testing"

	"git.servidordomal.lol/robogg133/superultraomegadeploy/internal/repo/git"
)

const DeployKey string = ``

func TestClonePublic(t *testing.T) {

	dir := filepath.Join(t.ArtifactDir(), "plain_clone_test")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	g := new(git.Git)
	if err := g.ClonePublic(t.Context(), "https://github.com/robogg133/gonion.git", dir); err != nil {
		t.Fatal(err)
	}
	s, err := g.HeadHashString()
	if err != nil {
		t.Error(err)
	}
	t.Log("HEAD hash id:", s)
}

func TestCloneDeployKey(t *testing.T) {

	dir := filepath.Join(t.ArtifactDir(), "deploy_key_plain_clone_test")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}

	g := new(git.Git)
	g.DeployKey = []byte(DeployKey)
	if err := g.Clone(t.Context(), "", dir); err != nil {
		t.Fatal(err)
	}
	s, err := g.HeadHashString()
	if err != nil {
		t.Error(err)
	}
	t.Log("HEAD hash id:", s)
}
