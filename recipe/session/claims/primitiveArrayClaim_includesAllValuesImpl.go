package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type includesAllValuesImpl struct {
	id              string
	claim           *PrimitiveArrayClaim
	maxAgeInSeconds *int64
	vals            []interface{}
}

func (impl *includesAllValuesImpl) GetID() string {
	return impl.id
}

func (impl *includesAllValuesImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *includesAllValuesImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]interface{})
	if !ok || val == nil {
		return true
	}
	if impl.maxAgeInSeconds != nil {
		return val["t"].(int64) < time.Now().Unix()-*impl.maxAgeInSeconds
	}
	return false
}

func (impl *includesAllValuesImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
	claimVal, claimValOk := impl.claim.GetValueFromPayload(payload, userContext).([]interface{})
	assertCondition(claimValOk, "claim value not an array")

	if claimVal == nil {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":           "value does not exist",
				"expectedToInclude": impl.vals,
				"actualValue":       claimVal,
			},
		}
	}
	ageInSeconds := time.Now().Unix() - impl.claim.GetLastRefetchTime(payload, userContext)
	if impl.maxAgeInSeconds != nil && ageInSeconds > *impl.maxAgeInSeconds {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":         "expired",
				"ageInSeconds":    ageInSeconds,
				"maxAgeInSeconds": *impl.maxAgeInSeconds,
			},
		}
	}

	isValid := true
	valsMap := map[interface{}]bool{}
	for _, v := range impl.vals {
		valsMap[v] = true
	}
	for _, v := range claimVal {
		if !valsMap[v] {
			isValid = false
			break
		}
	}

	if !isValid {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":           "wrong value",
				"expectedToInclude": impl.vals,
				"actualValue":       claimVal,
			},
		}
	}
	return ClaimValidationResult{
		IsValid: true,
	}
}
