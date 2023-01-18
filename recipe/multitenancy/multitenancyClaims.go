package multitenancy

import (
	"errors"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func NewAllowedDomainsClaim() (*claims.TypeSessionClaim, claims.PrimitiveArrayClaimValidators) {
	fetchDomains := func(userId string, userContext supertokens.UserContext) (interface{}, error) {
		instance, err := GetRecipeInstanceOrThrowError()
		if err != nil {
			return nil, err
		}

		tenantIdRes, err := instance.GetTenantIdForUserID(userId, userContext)
		if err != nil {
			return false, err
		}

		if tenantIdRes.OK != nil {
			if instance.GetAllowedDomainsForTenantId == nil {
				return []interface{}{}, nil // User did not provide a function to get allowed domains, but is using a validator. So we don't allow any domains by default
			}

			domains, err := instance.GetAllowedDomainsForTenantId(tenantIdRes.OK.TenantId, userContext)
			if err != nil {
				return nil, err
			}

			// Need this conversion because the array claim expects an array of type []interface{}
			domainsArray := make([]interface{}, len(domains))
			for i, domain := range domains {
				domainsArray[i] = domain
			}
			return domainsArray, nil
		} else {
			return nil, errors.New("UNKNOWN_USER_ID")
		}
	}

	var defaultMaxAge int64 = 3600
	allowedDomainsClaim, allowedDomainsClaimValidators := claims.PrimitiveArrayClaim("st-tenant-domains", fetchDomains, &defaultMaxAge)

	oGetValueFromPayload := allowedDomainsClaim.GetValueFromPayload
	allowedDomainsClaim.GetValueFromPayload = func(payload map[string]interface{}, userContext supertokens.UserContext) interface{} {
		value := oGetValueFromPayload(payload, userContext)

		if value == nil {
			return []interface{}{}
		}

		return value
	}

	oGetLastRefetchTime := allowedDomainsClaim.GetLastRefetchTime
	allowedDomainsClaim.GetLastRefetchTime = func(payload map[string]interface{}, userContext supertokens.UserContext) *int64 {
		value := oGetLastRefetchTime(payload, userContext)
		if value == nil {
			val := int64(time.Now().UnixNano() / 1000000)
			return &val
		}

		return value
	}

	return allowedDomainsClaim, allowedDomainsClaimValidators
}
