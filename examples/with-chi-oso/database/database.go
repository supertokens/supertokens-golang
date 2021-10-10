package database

import (
	"fmt"

	"github.com/supertokens/supertokens-golang/examples/with-chi-oso/models"
)

type Db struct {
	users map[string]models.User
	repos map[string]models.Repository
}

func NewDb() *Db {
	return &Db{
		users: map[string]models.User{
			"larry@supertokens.io":  {Roles: []models.RepositoryRole{{Role: "admin", RepoId: 0}}},
			"anne@supertokens.io":   {Roles: []models.RepositoryRole{{Role: "maintainer", RepoId: 1}}},
			"graham@supertokens.io": {Roles: []models.RepositoryRole{{Role: "contributor", RepoId: 2}}},
		},
		repos: map[string]models.Repository{
			"gmail": {Id: 0, Name: "gmail"},
			"react": {Id: 0, Name: "react", IsPublic: true},
			"oso":   {Id: 0, Name: "oso"},
		},
	}
}

func (d *Db) GetRepositoryByName(name string) (models.Repository, error) {
	if r, ok := d.repos[name]; ok {
		return r, nil
	}
	return models.Repository{}, fmt.Errorf("repository not found: %s", name)
}

func (d *Db) GetCurrentUser(email string) (models.User, error) {
	if u, ok := d.users[email]; ok {
		return u, nil
	}
	return models.User{}, fmt.Errorf("user not found: %s", email)
}
