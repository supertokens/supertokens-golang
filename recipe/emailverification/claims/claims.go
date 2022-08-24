package claims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeEmailVerificationClaim struct {
	*claims.TypeBooleanClaim

	Validators *EmailVerificationClaimValidators
}

type EmailVerificationClaimValidators struct {
	*claims.BooleanClaimValidators
	IsVerified func(refetchTimeOnFalseInSeconds *int64) *claims.SessionClaimValidator
}

var EmailVerificationClaim *TypeEmailVerificationClaim
