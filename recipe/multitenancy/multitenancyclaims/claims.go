package multitenancyclaims

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
)

var AllowedDomainsClaim *claims.TypeSessionClaim
var AllowedDomainsClaimValidators claims.PrimitiveArrayClaimValidators
