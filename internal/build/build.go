package build

import (
	"path"

	"github.com/nicompile/cabin/internal/docker"
	"github.com/nicompile/cabin/internal/generator"
	"github.com/nicompile/cabin/internal/git"
	"github.com/nicompile/infra-library-go/pkg/config"
)

func Build(url string, branch string) error {
	repo, err := git.Clone(url, branch)
	if err != nil {
		return err
	}

	configs, err := config.GetAll(path.Join(repo.Dir, ".nicompile"))
	if err != nil {
		return err
	}

	err = generator.GenerateAll(repo, configs)
	if err != nil {
		return err
	}

	err = docker.MakeDockerfiles(repo, configs)
	if err != nil {
		return err
	}

	//err = git.Cleanup(repo)
	//if err != nil {
	//return err
	//}
	return nil
}
