package evclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeEmailVerificationClaimValidators struct {
	claims.BooleanClaimValidators
	IsVerified func(refetchTimeOnFalseInSeconds *int64, maxAgeInSeconds *int64) claims.SessionClaimValidator
}

var EmailVerificationClaim *claims.TypeSessionClaim

var EmailVerificationClaimValidators TypeEmailVerificationClaimValidators
