package main

import (
	"context"
	"errors"
	"io"

	goGit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.f110.dev/xerrors"

	"go.f110.dev/mono/go/pkg/git"
	"go.f110.dev/mono/go/pkg/storage"
)

type gitDataService struct {
	repo map[string]*goGit.Repository
}

type repository struct {
	Name   string
	Prefix string
}

type ObjectStorageInterface interface {
	PutReader(ctx context.Context, name string, data io.Reader) error
	Delete(ctx context.Context, name string) error
	Get(ctx context.Context, name string) (io.ReadCloser, error)
	List(ctx context.Context, prefix string) ([]*storage.Object, error)
}

func newService(s ObjectStorageInterface, repositories []repository) (*gitDataService, error) {
	repo := make(map[string]*goGit.Repository)
	for _, r := range repositories {
		storer := git.NewObjectStorageStorer(s, r.Prefix)
		gitRepo, err := goGit.Open(storer, nil)
		if err != nil {
			return nil, xerrors.WithStack(err)
		}

		repo[r.Name] = gitRepo
	}

	return &gitDataService{repo: repo}, nil
}

func (g *gitDataService) ListRepositories(_ context.Context, _ *git.RequestListRepositories) (*git.ResponseListRepositories, error) {
	var list []string
	for k := range g.repo {
		list = append(list, k)
	}

	return &git.ResponseListRepositories{Repositories: list}, nil
}

func (g *gitDataService) ListReferences(_ context.Context, req *git.RequestListReferences) (*git.ResponseListReferences, error) {
	repo, ok := g.repo[req.Repo]
	if !ok {
		return nil, errors.New("repository not found")
	}

	refs, err := repo.References()
	if err != nil {
		return nil, err
	}

	res := &git.ResponseListReferences{}
	for {
		ref, err := refs.Next()
		if err == io.EOF {
			break
		}
		res.Refs = append(res.Refs, &git.Reference{
			Name:   ref.Name().String(),
			Hash:   ref.Hash().String(),
			Target: ref.Target().String(),
		})
	}
	return res, nil
}

func (g *gitDataService) GetCommit(_ context.Context, req *git.RequestGetCommit) (*git.ResponseGetCommit, error) {
	repo, ok := g.repo[req.Repo]
	if !ok {
		return nil, errors.New("repository not found")
	}
	if req.Sha == "" {
		return nil, errors.New("SHA field is required")
	}

	h := plumbing.NewHash(req.Sha)
	commit, err := repo.CommitObject(h)
	if err != nil {
		return nil, err
	}

	res := &git.ResponseGetCommit{
		Commit: &git.Commit{
			Sha:     commit.Hash.String(),
			Message: commit.Message,
			Committer: &git.Signature{
				Name:  commit.Committer.Name,
				Email: commit.Committer.Email,
			},
			Author: &git.Signature{
				Name:  commit.Author.Name,
				Email: commit.Author.Email,
			},
			Tree: commit.TreeHash.String(),
		},
	}
	if len(commit.ParentHashes) > 0 {
		parents := make([]string, len(commit.ParentHashes))
		for i := 0; i < len(commit.ParentHashes); i++ {
			parents[i] = commit.ParentHashes[i].String()
		}
		res.Commit.Parents = parents
	}

	return res, nil
}

func (g *gitDataService) GetTree(ctx context.Context, tree *git.RequestGetTree) (*git.ResponseGetTree, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gitDataService) GetBlob(ctx context.Context, blob *git.RequestGetBlob) (*git.ResponseGetBlob, error) {
	//TODO implement me
	panic("implement me")
}