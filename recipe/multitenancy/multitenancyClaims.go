package multitenancy

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancyclaims"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func NewAllowedDomainsClaim() (*claims.TypeSessionClaim, multitenancyclaims.TypeAllowedDomainsClaimValidators) {
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
	allowedDomainsClaim, arrayClaimValidators := claims.PrimitiveArrayClaim("st-tenant-domains", fetchDomains, &defaultMaxAge)

	allowedDomainsClaimValidators := multitenancyclaims.TypeAllowedDomainsClaimValidators{
		PrimitiveArrayClaimValidators: arrayClaimValidators,
		CheckAccessToDomain: func(allowedDomain string, maxAgeInSeconds *int64) claims.SessionClaimValidator {
			if maxAgeInSeconds == nil {
				var defaultTimeout int64 = 3600
				maxAgeInSeconds = &defaultTimeout
			}

			claimValidator := arrayClaimValidators.Includes(allowedDomain, maxAgeInSeconds, nil)
			oShouldRefetch := claimValidator.ShouldRefetch
			claimValidator.ShouldRefetch = func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
				claimVal := allowedDomainsClaim.GetValueFromPayload(payload, userContext)
				if claimVal == nil {
					return false
				}
				return oShouldRefetch(payload, userContext)
			}
			return claimValidator
		},
	}
	return allowedDomainsClaim, allowedDomainsClaimValidators
}
