package githuborg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"log"
)

type JsonRepo struct {
	Name string
	Repo
}

func loadRepo(filename string) (Repo, error) {
	content, err := os.ReadFile(filename)
	repo := Repo{}
	JsonRepo := map[string]Repo{}

	if err != nil {
		fmt.Println(err)
		return repo, err
	}

	err = json.Unmarshal(content, &JsonRepo)

	if err != nil {
		fmt.Println(err)
		return repo, err
	}

	for _, v := range JsonRepo {
		repo = v
	}

	if len(repo.Topics) == 0 {
		repo.Topics = []string{}
	}

	return repo, nil
}

func LoadConfig(path string) (map[string]Repo, error) {
	repos := map[string]Repo{}
	f, err := os.Open(path)

	if err != nil {
		log.Println(err)
		return repos, err
	}

	files, err := f.ReadDir(0)

	if err != nil {
		log.Println(err)
		return repos, err
	}

	for _, v := range files {
		if !v.IsDir() {
			repo, err := loadRepo(filepath.Join(path, v.Name()))
			if err != nil {
				log.Println(err)
				continue
			}
			repos[*repo.Name] = repo
		}
	}

	return repos, nil
}
