package main

type Repository struct {
	Id       int
	Name     string
	IsPublic bool
}

var reposDb = map[string]Repository{
	"gmail": {Id: 0, Name: "gmail"},
	"react": {Id: 1, Name: "react", IsPublic: true},
	"oso":   {Id: 2, Name: "oso"},
}

func GetRepositoryByName(name string) Repository {
	return reposDb[name]
}

type RepositoryRole struct {
	Role   string
	RepoId int
}

type User struct {
	Roles []RepositoryRole
}

var usersDb = map[string]User{
	"larry":  {Roles: []RepositoryRole{{Role: "admin", RepoId: 0}}},
	"anne":   {Roles: []RepositoryRole{{Role: "maintainer", RepoId: 1}}},
	"graham": {Roles: []RepositoryRole{{Role: "contributor", RepoId: 2}}},
}

func GetCurrentUser(userID string) User {
	return usersDb[userID]
}
