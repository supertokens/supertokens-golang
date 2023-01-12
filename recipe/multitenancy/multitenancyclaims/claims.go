package multitenancyclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeMultitenancyClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
	CheckAccessToDomain func(allowedDomain string, maxAgeInSeconds *int64) claims.SessionClaimValidator
}

var MultitenancyTenantIdClaim *claims.TypeSessionClaim
var MultitenancyDomainsClaim *claims.TypeSessionClaim
var MultitenancyValidators TypeMultitenancyClaimValidators
