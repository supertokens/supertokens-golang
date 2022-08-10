package claims

import "github.com/supertokens/supertokens-golang/supertokens"

type hasValueImpl struct {
	id    string
	claim SessionClaim
	val   interface{}
}

func (impl *hasValueImpl) GetID() string {
	return impl.id
}

func (impl *hasValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *hasValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val := impl.claim.GetValueFromPayload(payload, userContext)
	return val == nil
}

func (impl *hasValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
	claimVal := impl.claim.GetValueFromPayload(payload, userContext)
	isValid := claimVal == impl.val
	if isValid {
		return ClaimValidationResult{
			IsValid: true,
		}
	} else {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":       "wrong value",
				"expectedValue": impl.val,
				"actualValue":   claimVal,
			},
		}
	}
}
