package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type PrimitiveClaim struct {
	Key        string
	fetchValue func(userId string, userContext supertokens.UserContext) interface{}
}

func (claim *PrimitiveClaim) FetchValue(userId string, userContext supertokens.UserContext) interface{} {
	return claim.fetchValue(userId, userContext)
}

func (claim *PrimitiveClaim) AddToPayload_internal(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{} {
	payload[claim.Key] = map[string]interface{}{
		"v": value,
		"t": time.Now().Unix(),
	}

	return payload
}

func (claim *PrimitiveClaim) RemoveFromPayloadByMerge_internal(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
	payload[claim.Key] = nil
	return payload
}

func (claim *PrimitiveClaim) RemoveFromPayload(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
	delete(payload, claim.Key)
	return payload
}

func (claim *PrimitiveClaim) GetValueFromPayload(payload map[string]interface{}, userContext supertokens.UserContext) interface{} {
	if value, ok := payload[claim.Key].(map[string]interface{}); ok {
		return value["v"]
	}
	return nil
}

func (claim *PrimitiveClaim) GetLastRefetchTime(payload map[string]interface{}, userContext supertokens.UserContext) int64 {
	if value, ok := payload[claim.Key].(map[string]interface{}); ok {
		return value["t"].(int64)
	}
	return 0
}

func (claim *PrimitiveClaim) GetValidators() PrimitiveClaimValidators {
	return PrimitiveClaimValidators{
		HasValue: func(val interface{}, id *string) SessionClaimValidator {
			claimId := claim.Key
			if id != nil {
				claimId = *id
			}
			return &HasValueImpl{
				id:    claimId,
				claim: claim,
				val:   val,
			}
		},
		HasFreshValue: func(val interface{}, maxAgeInSeconds int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-freshVal"
			if id != nil {
				claimId = *id
			}
			return &HasFreshValueImpl{
				id:              claimId,
				claim:           claim,
				maxAgeInSeconds: maxAgeInSeconds,
				val:             val,
			}
		},
	}
}

func (claim *PrimitiveClaim) getPrimitiveClaim() *PrimitiveClaim {
	return claim
}

type PrimitiveClaimValidators struct {
	HasValue      func(val interface{}, id *string) SessionClaimValidator
	HasFreshValue func(val interface{}, maxAgeInSeconds int64, id *string) SessionClaimValidator
}

type HasValueImpl struct {
	id    string
	claim SessionClaim
	val   interface{}
}

func (impl *HasValueImpl) GetID() string {
	return impl.id
}

func (impl *HasValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *HasValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val := impl.claim.GetValueFromPayload(payload, userContext)
	return val == nil
}

func (impl *HasValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
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

type HasFreshValueImpl struct {
	id              string
	claim           SessionClaim
	val             interface{}
	maxAgeInSeconds int64
}

func (impl *HasFreshValueImpl) GetID() string {
	return impl.id
}

func (impl *HasFreshValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *HasFreshValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]interface{})
	if !ok || val == nil {
		return true
	}
	return val["t"].(int64) < time.Now().Unix()-impl.maxAgeInSeconds
}

func (impl *HasFreshValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
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

type isPrimitiveClaim interface {
	getPrimitiveClaim() *PrimitiveClaim
}
