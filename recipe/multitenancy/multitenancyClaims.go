package multitenancy

import (
	"errors"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancyclaims"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func NewMultitenancyClaims() (*claims.TypeSessionClaim, *claims.TypeSessionClaim, multitenancyclaims.TypeMultitenancyClaimValidators) {
	fetchTenantId := func(userId string, userContext supertokens.UserContext) (interface{}, error) {
		instance, err := GetRecipeInstanceOrThrowError()
		if err != nil {
			return nil, err
		}
		tenantIdRes, err := instance.GetTenantIdForUserID(userId, userContext)
		if err != nil {
			return false, err
		}
		if tenantIdRes.OK != nil {
			return tenantIdRes.OK.TenantId, nil
		} else {
			return "", errors.New("UNKNOWN_USER_ID")
		}
	}

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
			tenantId := tenantIdRes.OK.TenantId
			domains, err := instance.GetDomainsForTenantId(tenantId, userContext)
			if err != nil {
				return nil, err
			}
			return domains, nil
		} else {
			// This may be Unknown user id error or the user may belong to a non thirdparty recipe, in which case, the domains can be assumed to be an empty array
			return []interface{}{}, nil
		}
	}

	var defaultMaxAge int64 = 300
	mtDomainClaim, arrayClaimValidators := claims.PrimitiveArrayClaim("st-tenant-domains", fetchDomains, &defaultMaxAge)
	mtTenantIdClaim, _ := claims.PrimitiveClaim("st-tenantid", fetchTenantId, &defaultMaxAge)

	getLastRefetchTime := func(payload map[string]interface{}, userContext supertokens.UserContext) *int64 {
		if value, ok := payload[mtDomainClaim.Key].(map[string]interface{}); ok {
			switch t := value["t"].(type) {
			case int64:
				return &t
			case float64:
				it := int64(t)
				return &it
			}
		}
		return nil
	}

	validators := multitenancyclaims.TypeMultitenancyClaimValidators{
		PrimitiveArrayClaimValidators: arrayClaimValidators,
		CheckAccessToDomain: func(allowedDomain string, refetchTimeOnFalseInSeconds *int64, maxAgeInSeconds *int64) claims.SessionClaimValidator {
			if refetchTimeOnFalseInSeconds == nil {
				var defaultTimeout int64 = 10
				refetchTimeOnFalseInSeconds = &defaultTimeout
			}
			if maxAgeInSeconds == nil {
				var defaultTimeout int64 = 300
				maxAgeInSeconds = &defaultTimeout
			}

			claimValidator := arrayClaimValidators.Includes(allowedDomain, maxAgeInSeconds, nil)
			claimValidator.ShouldRefetch = func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
				value := mtDomainClaim.GetValueFromPayload(payload, userContext)
				return value == nil || (*getLastRefetchTime(payload, userContext) < time.Now().UnixNano()/1000000-*maxAgeInSeconds*1000) || (*getLastRefetchTime(payload, userContext) < time.Now().UnixNano()/1000000-*refetchTimeOnFalseInSeconds*1000)
			}
			return claimValidator
		},
	}
	return mtTenantIdClaim, mtDomainClaim, validators
}
