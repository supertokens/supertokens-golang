package multitenancyclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeMultitenancyClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
	CheckAccessToDomain func(allowedDomain string, maxAgeInSeconds *int64) claims.SessionClaimValidator
}

var MultitenancyDomainsClaim *claims.TypeSessionClaim
var MultitenancyValidators TypeMultitenancyClaimValidators
