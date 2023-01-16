package multitenancyclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

type TypeAllowedDomainsClaimValidators struct {
	claims.PrimitiveArrayClaimValidators
	CheckAccessToDomain func(allowedDomain string, maxAgeInSeconds *int64) claims.SessionClaimValidator
}

var AllowedDomainsClaim *claims.TypeSessionClaim
var AllowedDomainsClaimValidators TypeAllowedDomainsClaimValidators
