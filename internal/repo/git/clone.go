package git

import (
	"context"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

type Git struct {
	DeployKey []byte
	repo      *git.Repository
}

func (g *Git) ClonePublic(ctx context.Context, url, dst string) error {

	var err error
	g.repo, err = git.PlainCloneContext(ctx, dst, &git.CloneOptions{
		URL:   url,
		Depth: 1,
	})
	return err
}
func (g *Git) Clone(ctx context.Context, url, dst string) error {
	pks, err := ssh.NewPublicKeys("git", g.DeployKey, "")
	if err != nil {
		return err
	}

	repo, err := git.PlainCloneContext(ctx, dst, &git.CloneOptions{
		URL:   url,
		Depth: 1,
		ClientOptions: []client.Option{
			client.WithSSHAuth(pks),
		},
	})
	g.repo = repo

	return err
}

func (g *Git) HeadHashString() (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", err
	}
	return head.Hash().String(), nil
}

func (g *Git) Checkout(ctx context.Context, branch string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(plumbing.NewBranchReferenceName(branch)),
		Force:  true,
	})
}
