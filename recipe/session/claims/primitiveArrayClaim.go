package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type PrimitiveArrayClaim struct {
	Key        string
	fetchValue func(userId string, userContext supertokens.UserContext) interface{}
}

func (claim *PrimitiveArrayClaim) FetchValue(userId string, userContext supertokens.UserContext) interface{} {
	return claim.fetchValue(userId, userContext)
}

func (claim *PrimitiveArrayClaim) AddToPayload_internal(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{} {
	_, ok := value.([]interface{})
	assertCondition(ok, "value not an array")
	payload[claim.Key] = map[string]interface{}{
		"v": value,
		"t": time.Now().Unix(),
	}
	return payload
}

func (claim *PrimitiveArrayClaim) RemoveFromPayloadByMerge_internal(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
	payload[claim.Key] = nil
	return payload
}

func (claim *PrimitiveArrayClaim) RemoveFromPayload(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
	delete(payload, claim.Key)
	return payload
}

func (claim *PrimitiveArrayClaim) GetValueFromPayload(payload map[string]interface{}, userContext supertokens.UserContext) interface{} {
	if value, ok := payload[claim.Key].(map[string]interface{}); ok {
		return value["v"]
	}
	return nil
}

func (claim *PrimitiveArrayClaim) GetLastRefetchTime(payload map[string]interface{}, userContext supertokens.UserContext) int64 {
	if value, ok := payload[claim.Key].(map[string]interface{}); ok {
		return value["t"].(int64)
	}
	return 0
}

func (claim *PrimitiveArrayClaim) GetValidators() PrimitiveArrayClaimValidators {
	return PrimitiveArrayClaimValidators{
		Includes: func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-includes"
			if id != nil {
				claimId = *id
			}
			return &IncludesValueImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId, val: val}
		},
		Excludes: func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-excludes"
			if id != nil {
				claimId = *id
			}
			return &ExcludesValueImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId, val: val}
		},
	}
}

type IncludesValueImpl struct {
	id              string
	claim           *PrimitiveArrayClaim
	maxAgeInSeconds *int64
	val             interface{}
}

func (impl *IncludesValueImpl) GetID() string {
	return impl.id
}

func (impl *IncludesValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *IncludesValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]interface{})
	if !ok || val == nil {
		return true
	}
	if impl.maxAgeInSeconds != nil {
		return val["t"].(int64) < time.Now().Unix()-*impl.maxAgeInSeconds
	}
	return false
}

func (impl *IncludesValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
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
	if !includes(claimVal, impl.val) {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]interface{}{
				"message":           "wrong value",
				"expectedToInclude": impl.val,
				"actualValue":       claimVal,
			},
		}
	}
	return ClaimValidationResult{
		IsValid: true,
	}
}

type ExcludesValueImpl struct {
	id              string
	claim           *PrimitiveArrayClaim
	maxAgeInSeconds *int64
	val             interface{}
}

func (impl *ExcludesValueImpl) GetID() string {
	return impl.id
}

func (impl *ExcludesValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *ExcludesValueImpl) ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]interface{})
	if !ok || val == nil {
		return true
	}
	if impl.maxAgeInSeconds != nil {
		return val["t"].(int64) < time.Now().Unix()-*impl.maxAgeInSeconds
	}
	return false
}

func (impl *ExcludesValueImpl) Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
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

type PrimitiveArrayClaimValidators struct {
	Includes    func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	Excludes    func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	IncludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	ExcludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
}
