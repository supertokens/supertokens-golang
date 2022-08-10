package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type hasFreshValueImpl struct {
	id              string
	claim           SessionClaim
	val             interface{}
	maxAgeInSeconds int64
}

func (impl *hasFreshValueImpl) GetID() string {
	return impl.id
}

func (impl *hasFreshValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *hasFreshValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]interface{})
	if !ok || val == nil {
		return true
	}
	return val["t"].(int64) < time.Now().Unix()-impl.maxAgeInSeconds
}

func (impl *hasFreshValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
	claimVal := impl.claim.GetValueFromPayload(payload, userContext)
	primitiveClaim, isPrimitiveClaimOk := impl.claim.(isPrimitiveClaim)
	if !isPrimitiveClaimOk {
		panic("Claim is not a primitive claim")
	}
	if claimVal == nil {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":       "value does not exist",
				"expectedValue": impl.val,
				"actualValue":   claimVal,
			},
		}
	}
	ageInSeconds := time.Now().Unix() - primitiveClaim.getPrimitiveClaim().GetLastRefetchTime(payload, userContext)
	if ageInSeconds > impl.maxAgeInSeconds {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":         "expired",
				"ageInSeconds":    ageInSeconds,
				"maxAgeInSeconds": impl.maxAgeInSeconds,
			},
		}
	}
	if claimVal != impl.val {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":       "wrong value",
				"expectedValue": impl.val,
				"actualValue":   claimVal,
			},
		}
	}
	return ClaimValidationResult{
		IsValid: true,
	}
}
