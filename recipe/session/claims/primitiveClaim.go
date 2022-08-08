package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type PrimitiveClaim struct {
	Key        string
	fetchValue func(userId string, userContext supertokens.UserContext) Any
}

func (claim *PrimitiveClaim) FetchValue(userId string, userContext supertokens.UserContext) Any {
	return claim.fetchValue(userId, userContext)
}

func (claim *PrimitiveClaim) AddToPayload_internal(payload map[string]Any, value Any, userContext supertokens.UserContext) map[string]Any {
	payload[claim.Key] = map[string]Any{
		"v": value,
		"t": time.Now().Unix(),
	}
	return payload
}

func (claim *PrimitiveClaim) RemoveFromPayloadByMerge_internal(payload map[string]Any, userContext supertokens.UserContext) map[string]Any {
	payload[claim.Key] = nil
	return payload
}

func (claim *PrimitiveClaim) RemoveFromPayload(payload map[string]Any, userContext supertokens.UserContext) map[string]Any {
	delete(payload, claim.Key)
	return payload
}

func (claim *PrimitiveClaim) GetValueFromPayload(payload map[string]Any, userContext supertokens.UserContext) Any {
	if value, ok := payload[claim.Key].(map[string]Any); ok {
		return value["v"]
	}
	return nil
}

func (claim *PrimitiveClaim) GetLastRefetchTime(payload map[string]Any, userContext supertokens.UserContext) int64 {
	if value, ok := payload[claim.Key].(map[string]Any); ok {
		return value["t"].(int64)
	}
	return 0
}

func (claim *PrimitiveClaim) GetValidators() PrimitiveClaimValidators {
	return PrimitiveClaimValidators{
		HasValue: func(val Any, id *string) SessionClaimValidator {
			claimId := claim.Key
			if id != nil {
				claimId = *id
			}
			return &HasValueImpl{claim: claim, val: val, id: claimId}
		},
		HasFreshValue: func(val Any, maxAgeInSeconds int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-freshVal"
			if id != nil {
				claimId = *id
			}
			return &HasFreshValueImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId}
		},
	}
}

func (claim *PrimitiveClaim) getPrimitiveClaim() *PrimitiveClaim {
	return claim
}

type PrimitiveClaimValidators struct {
	HasValue      func(val Any, id *string) SessionClaimValidator
	HasFreshValue func(val Any, maxAgeInSeconds int64, id *string) SessionClaimValidator
}

type HasValueImpl struct {
	id    string
	claim SessionClaim
	val   Any
}

func (impl *HasValueImpl) GetID() string {
	return impl.id
}

func (impl *HasValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *HasValueImpl) ShouldRefetch(payload map[string]Any, userContext supertokens.UserContext) bool {
	val := impl.claim.GetValueFromPayload(payload, userContext)
	return val == nil
}

func (impl *HasValueImpl) Validate(payload map[string]Any, userContext supertokens.UserContext) ClaimValidationResult {
	claimVal := impl.claim.GetValueFromPayload(payload, userContext)
	isValid := claimVal == impl.val
	if isValid {
		return ClaimValidationResult{
			IsValid: true,
		}
	} else {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]Any{
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
	val             Any
	maxAgeInSeconds int64
}

func (impl *HasFreshValueImpl) GetID() string {
	return impl.id
}

func (impl *HasFreshValueImpl) GetClaim() SessionClaim {
	return impl.claim
}

func (impl *HasFreshValueImpl) ShouldRefetch(payload map[string]Any, userContext supertokens.UserContext) bool {
	val, ok := impl.claim.GetValueFromPayload(payload, userContext).(map[string]Any)
	if !ok || val == nil {
		return true
	}
	return val["t"].(int64) < time.Now().Unix()-impl.maxAgeInSeconds
}

func (impl *HasFreshValueImpl) Validate(payload map[string]Any, userContext supertokens.UserContext) ClaimValidationResult {
	claimVal := impl.claim.GetValueFromPayload(payload, userContext)
	primitiveClaim, isPrimitiveClaimOk := impl.claim.(isPrimitiveClaim)
	if !isPrimitiveClaimOk {
		panic("Claim is not a primitive claim")
	}
	if claimVal == nil {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]Any{
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
			Reason: map[string]Any{
				"message":         "expired",
				"ageInSeconds":    ageInSeconds,
				"maxAgeInSeconds": impl.maxAgeInSeconds,
			},
		}
	}
	if claimVal != impl.val {
		return ClaimValidationResult{
			IsValid: false,
			Reason: map[string]Any{
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
