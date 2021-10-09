package database

import (
	"fmt"
	"os"
	"path"

	"github.com/supertokens/supertokens-golang/examples/with-chi-oso/models"
	"github.com/viant/dsc"
)

type Db struct {
	manager dsc.Manager
}

func NewDb(name string) *Db {
	f := dsc.NewManagerFactory()
	cwd, _ := os.Getwd()
	cfg := dsc.NewConfig(
		"ndjson",
		"",
		fmt.Sprintf("url:%s,ext:json,dateFormat:yyyy-MM-dd hh:mm:ss", path.Join(cwd, name)),
	)
	manager, err := f.Create(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize db: %s", err)
		os.Exit(1)
	}
	return &Db{manager: manager}
}

func (d *Db) GetRepositoryByName(name string) (models.Repository, error) {
	var repo models.Repository
	success, err := d.manager.ReadSingle(
		&repo,
		"SELECT id, name, ispublic FROM repos WHERE name = ?",
		[]interface{}{name},
		nil,
	)
	if err != nil {
		return repo, err
	}
	if !success {
		return models.Repository{}, fmt.Errorf("repository not found: %s", name)
	}
	return repo, nil
}

func (d *Db) GetCurrentUser(email string) (models.User, error) {
	var u models.User
	success, err := d.manager.ReadSingle(
		&u,
		"SELECT email, roles FROM users WHERE email = ?",
		[]interface{}{email},
		nil,
	)
	if err != nil {
		return u, err
	}
	if !success {
		return models.User{}, fmt.Errorf("user not found: %s", email)
	}
	return u, nil
}
