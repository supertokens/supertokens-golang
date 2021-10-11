allow(actor, action, resource) if
  has_permission(actor, action, resource);

actor User {}

resource Repository {
	permissions = ["read", "push", "delete"];
	roles = ["contributor", "maintainer", "admin"];

	"read" if "contributor";
	"push" if "maintainer";
	"delete" if "admin";

	"maintainer" if "admin";
	"contributor" if "maintainer";
}

has_role(user: User, roleName: String, repository: Repository) if
  role in user.Roles and
  role.Role = roleName and
  role.RepoId = repository.Id;