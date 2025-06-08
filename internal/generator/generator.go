package generator

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/dave/jennifer/jen"
	"github.com/nicompile/cabin/internal/git"
	"github.com/nicompile/infra-library-go/pkg/config"
)

func GenerateAll(repo git.Repo, configs []config.Config) error {
	for _, config := range configs {
		err := generate(repo, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func generate(repo git.Repo, conf config.Config) error {
	dir := path.Join(repo.Dir, "build", conf.Uid)
	err := exec.Command("mkdir", "-p", dir).Run()
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(dir, conf.Uid+".go"))
	if err != nil {
		return err
	}
	defer f.Close()

	j := jen.NewFile("main")
	switch conf.Type {
	case config.SERVERLESS_FUNCTION:
		targetPackage := path.Join(repo.Modname, conf.Target)
		j.ImportAlias(targetPackage, "target")
		j.Func().Id("main").Params().Block(
			jen.Qual("github.com/nicompile/infra-library-go/pkg/serverlessfunction", "StartLambda").Params(
				jen.Qual(targetPackage, "Function"),
			),
		)
	case config.ONEOFF_JOB:
		return errors.New("not implemented")
	}

	err = j.Render(f)
	if err != nil {
		panic(err)
	}

	return nil
}
