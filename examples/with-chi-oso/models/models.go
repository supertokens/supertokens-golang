package models

type Repository struct {
	Id       int
	Name     string
	IsPublic bool
}

type RepositoryRole struct {
	Role   string
	RepoId int
}

type User struct {
	Roles []RepositoryRole
	Email string
}
