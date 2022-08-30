package userrolesclaims

import "github.com/supertokens/supertokens-golang/recipe/session/claims"

type TypeUserRoleClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
}

var UserRoleClaim claims.TypeSessionClaim
var UserRoleClaimValidators TypeUserRoleClaimValidators

type TypePermissionClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
}

var PermissionClaim claims.TypeSessionClaim
var PermissionClaimValidators TypePermissionClaimValidators
