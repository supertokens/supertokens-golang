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
			return &hasValueImpl{
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
			return &hasFreshValueImpl{
				id:              claimId,
				claim:           claim,
				maxAgeInSeconds: maxAgeInSeconds,
				val:             val,
			}
		},
	}
}

func (claim *PrimitiveClaim) Build(userId string, userContext supertokens.UserContext) map[string]interface{} {
	value := claim.fetchValue(userId, userContext)
	if value == nil {
		return map[string]interface{}{}
	}
	return claim.AddToPayload_internal(map[string]interface{}{}, value, userContext)
}

func (claim *PrimitiveClaim) getPrimitiveClaim() *PrimitiveClaim {
	return claim
}

type PrimitiveClaimValidators struct {
	HasValue      func(val interface{}, id *string) SessionClaimValidator
	HasFreshValue func(val interface{}, maxAgeInSeconds int64, id *string) SessionClaimValidator
}

type isPrimitiveClaim interface {
	getPrimitiveClaim() *PrimitiveClaim
}
