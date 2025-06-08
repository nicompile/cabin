package docker

import (
	"errors"
	"os"
	"path"
	"text/template"

	"github.com/nicompile/cabin/internal/git"
	"github.com/nicompile/infra-library-go/pkg/config"
)

type DockerParams struct {
	GoVersion string
	Uid       string
}

func MakeDockerfiles(repo git.Repo, configs []config.Config) error {
	for _, conf := range configs {
		err := makeDockerfile(repo, conf)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path.Join(repo.Dir, "docker-compose.yml"))
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("templates/docker-compose.tmpl")
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, configs)
	if err != nil {
		return err
	}

	return nil
}

func makeDockerfile(repo git.Repo, conf config.Config) error {
	var templateFile string
	switch conf.Type {
	case config.SERVERLESS_FUNCTION:
		templateFile = "templates/lambda-dockerfile.tmpl"
	case config.ONEOFF_JOB:
		return errors.New("not implemented")
	default:
		return errors.New("not a valid deployable type")
	}

	f, err := os.Create(path.Join(repo.Dir, "target", conf.Uid, "Dockerfile"))
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, DockerParams{GoVersion: repo.GoVersion, Uid: conf.Uid})
	if err != nil {
		return err
	}

	return nil
}
