package multitenancyclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeMultitenancyClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
	CheckAccessToDomain func(allowedDomain string, maxAgeInSeconds *int64) claims.SessionClaimValidator
}

var AllowDomainsClaim *claims.TypeSessionClaim
var AllowDomainsClaimValidators TypeMultitenancyClaimValidators
