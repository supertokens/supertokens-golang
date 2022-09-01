package userroles

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/userroles/userrolesclaims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func init() {
	// automatically called when this package is imported
	userrolesclaims.UserRoleClaim, userrolesclaims.UserRoleClaimValidators = NewUserRoleClaim()
	userrolesclaims.PermissionClaim, userrolesclaims.PermissionClaimValidators = NewPermissionClaim()
}

func NewUserRoleClaim() (*claims.TypeSessionClaim, userrolesclaims.TypeUserRoleClaimValidators) {
	fetchValue := func(userId string, userContext supertokens.UserContext) (interface{}, error) {
		recipe, err := getRecipeInstanceOrThrowError()
		if err != nil {
			return nil, err
		}
		roles, err := (*recipe.RecipeImpl.GetRolesForUser)(userId, userContext)
		if err != nil {
			return nil, err
		}

		rolesArray := make([]interface{}, len(roles.OK.Roles))
		for i, role := range roles.OK.Roles {
			rolesArray[i] = role
		}
		return rolesArray, nil
	}

	userRoleClaim, primitiveArrayClaimValidators := claims.PrimitiveArrayClaim("st-role", fetchValue, nil)
	return userRoleClaim, userrolesclaims.TypeUserRoleClaimValidators{
		PrimitiveArrayClaimValidators: primitiveArrayClaimValidators,
	}
}

func NewPermissionClaim() (*claims.TypeSessionClaim, userrolesclaims.TypePermissionClaimValidators) {
	fetchValue := func(userId string, userContext supertokens.UserContext) (interface{}, error) {
		recipe, err := getRecipeInstanceOrThrowError()
		if err != nil {
			return nil, err
		}
		roles, err := (*recipe.RecipeImpl.GetRolesForUser)(userId, userContext)
		if err != nil {
			return nil, err
		}

		permissionSet := map[string]bool{}
		for _, role := range roles.OK.Roles {
			permissions, err := (*recipe.RecipeImpl.GetPermissionsForRole)(role, userContext)
			if err != nil {
				return nil, err
			}
			for _, permission := range permissions.OK.Permissions {
				permissionSet[permission] = true
			}
		}

		result := []interface{}{}

		for perm := range permissionSet {
			result = append(result, perm)
		}

		return result, nil
	}

	permissionClaim, primitiveArrayClaimValidators := claims.PrimitiveArrayClaim("st-perm", fetchValue, nil)
	return permissionClaim, userrolesclaims.TypePermissionClaimValidators{
		PrimitiveArrayClaimValidators: primitiveArrayClaimValidators,
	}
}
