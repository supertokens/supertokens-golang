package multitenancyclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeMultitenancyClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
	HasAccessToCurrentDomain func(allowedDomain string, refetchTimeOnFalseInSeconds *int64, maxAgeInSeconds *int64) claims.SessionClaimValidator
}

var MultitenancyTenantIdClaim *claims.TypeSessionClaim
var MultitenancyDomainsClaim *claims.TypeSessionClaim
var MultitenancyValidators TypeMultitenancyClaimValidators
