package multitenancy

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func NewAllowedDomainsClaim() (*claims.TypeSessionClaim, claims.PrimitiveArrayClaimValidators) {
	fetchDomains := func(userId string, tenantId string, userContext supertokens.UserContext) (interface{}, error) {
		instance, err := GetRecipeInstanceOrThrowError()
		if err != nil {
			return nil, err
		}

		if instance.GetAllowedDomainsForTenantId == nil {
			return nil, nil // User did not provide a function to get allowed domains, but is using a validator. So we don't allow any domains by default
		}

		domains, err := instance.GetAllowedDomainsForTenantId(tenantId, userContext)
		if err != nil {
			return nil, err
		}

		// Need this conversion because the array claim expects an array of type []interface{}
		domainsArray := make([]interface{}, len(domains))
		for i, domain := range domains {
			domainsArray[i] = domain
		}
		return domainsArray, nil
	}

	var defaultMaxAge int64 = 3600
	allowedDomainsClaim, allowedDomainsClaimValidators := claims.PrimitiveArrayClaim("st-t-dmns", fetchDomains, &defaultMaxAge)

	return allowedDomainsClaim, allowedDomainsClaimValidators
}
