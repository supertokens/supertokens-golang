package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type PrimitiveArrayClaim struct {
	Key        string
	fetchValue func(userId string, userContext supertokens.UserContext) interface{}
}

func (claim *PrimitiveArrayClaim) GetKey() string {
	return claim.Key
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
			return &includesValueImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId, val: val}
		},
		Excludes: func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-excludes"
			if id != nil {
				claimId = *id
			}
			return &excludesValueImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId, val: val}
		},
		IncludesAll: func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-includesAll"
			if id != nil {
				claimId = *id
			}
			return &includesAllValuesImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId, vals: vals}
		},
		ExcludesAll: func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			claimId := claim.Key + "-excludesAll"
			if id != nil {
				claimId = *id
			}
			return &excludesAllValuesImpl{claim: claim, maxAgeInSeconds: maxAgeInSeconds, id: claimId, vals: vals}
		},
	}
}

func (claim *PrimitiveArrayClaim) Build(userId string, userContext supertokens.UserContext) map[string]interface{} {
	value := claim.fetchValue(userId, userContext)
	if value == nil {
		return map[string]interface{}{}
	}
	return claim.AddToPayload_internal(map[string]interface{}{}, value, userContext)
}

type PrimitiveArrayClaimValidators struct {
	Includes    func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	Excludes    func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	IncludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	ExcludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
}
