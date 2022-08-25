package userroles

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	urclaims "github.com/supertokens/supertokens-golang/recipe/userroles/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func init() {
	urclaims.UserRoleClaim = NewUserRoleClaim()
	urclaims.PermissionClaim = NewPermissionClaim()
}

func NewUserRoleClaim() *urclaims.TypeUserRoleClaim {
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

	primitiveArrayClaim := claims.PrimitiveArrayClaim("st-role", fetchValue, nil)
	return &urclaims.TypeUserRoleClaim{
		TypePrimitiveArrayClaim: primitiveArrayClaim,
		Validators: &urclaims.TypeUserRoleClaimValidators{
			PrimitiveArrayClaimValidators: primitiveArrayClaim.Validators,
		},
	}
}

func NewPermissionClaim() *urclaims.TypePermissionClaim {
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

	primitiveArrayClaim := claims.PrimitiveArrayClaim("st-perm", fetchValue, nil)
	return &urclaims.TypePermissionClaim{
		TypePrimitiveArrayClaim: primitiveArrayClaim,
		Validators: &urclaims.TypePermissionClaimValidators{
			PrimitiveArrayClaimValidators: primitiveArrayClaim.Validators,
		},
	}
}
