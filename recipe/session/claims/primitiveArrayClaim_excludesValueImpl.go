package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type excludesValueImpl struct {
	id              string
	claim           *PrimitiveArrayClaim
	maxAgeInSeconds *int64
	val             interface{}
}

func (impl *excludesValueImpl) GetID() string {
	return impl.id
}

func (impl *excludesValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *excludesValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]interface{})
	if !ok || val == nil {
		return true
	}
	if impl.maxAgeInSeconds != nil {
		return val["t"].(int64) < time.Now().Unix()-*impl.maxAgeInSeconds
	}
	return false
}

func (impl *excludesValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
	claimVal, claimValOk := impl.claim.GetValueFromPayload(payload, userContext).([]interface{})
	assertCondition(claimValOk, "claim value not an array")

	if claimVal == nil {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":           "value does not exist",
				"expectedToInclude": impl.val,
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
	if includes(claimVal, impl.val) {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":           "wrong value",
				"expectedToExclude": impl.val,
				"actualValue":       claimVal,
			},
		}
	}
	return ClaimValidationResult{
		IsValid: true,
	}
}
