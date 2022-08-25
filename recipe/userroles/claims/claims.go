package claims

import "github.com/supertokens/supertokens-golang/recipe/session/claims"

type TypeUserRoleClaim struct {
	*claims.TypePrimitiveArrayClaim
	Validators *TypeUserRoleClaimValidators
}

type TypeUserRoleClaimValidators struct {
	*claims.PrimitiveArrayClaimValidators
}

var UserRoleClaim *TypeUserRoleClaim

type TypePermissionClaim struct {
	*claims.TypePrimitiveArrayClaim
	Validators *TypePermissionClaimValidators
}

type TypePermissionClaimValidators struct {
	*claims.PrimitiveArrayClaimValidators
}

var PermissionClaim *TypePermissionClaim
