package git

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/mod/modfile"
)

type Repo struct {
	Dir       string
	Modname   string
	GoVersion string
}

func Clone(url string, branch string) (Repo, error) {
	dir := path.Join("/tmp", "cabin-"+uuid.NewString())

	err := exec.Command("mkdir", "-p", dir).Run()
	if err != nil {
		return Repo{}, err
	}

	err = exec.Command("git", "clone", url, dir).Run()
	if err != nil {
		return Repo{}, err
	}

	err = exec.Command("git", "-C", dir, "checkout", branch).Run()
	if err != nil {
		return Repo{}, err
	}

	modBytes, err := os.ReadFile(path.Join(dir, "go.mod"))
	if err != nil {
		return Repo{}, err
	}

	mod, err := modfile.Parse("go.mod", modBytes, nil)
	if err != nil {
		return Repo{}, err
	}

	return Repo{Dir: dir, Modname: mod.Module.Mod.Path, GoVersion: mod.Go.Version}, nil
}

func Cleanup(repo Repo) error {
	if !strings.HasPrefix(repo.Dir, "/tmp/") {
		panic("Tried to delete " + repo.Dir)
	}

	err := exec.Command("rm", "-rf", repo.Dir).Run()
	if err != nil {
		return err
	}
	return nil
}
