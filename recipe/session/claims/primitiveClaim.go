package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func PrimitiveClaim(key string, fetchValue FetchValueFunc, defaultMaxAgeInSeconds *int64) (*TypeSessionClaim, PrimitiveClaimValidators) {
	sessionClaim := SessionClaim(key, fetchValue)

	sessionClaim.AddToPayload_internal = func(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{} {
		payload[sessionClaim.Key] = map[string]interface{}{
			"v": value,
			"t": time.Now().UnixNano() / 1000000,
		}

		return payload
	}

	sessionClaim.RemoveFromPayloadByMerge_internal = func(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
		payload[sessionClaim.Key] = nil
		return payload
	}

	sessionClaim.RemoveFromPayload = func(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
		delete(payload, sessionClaim.Key)
		return payload
	}

	sessionClaim.GetValueFromPayload = func(payload map[string]interface{}, userContext supertokens.UserContext) interface{} {
		if value, ok := payload[sessionClaim.Key].(map[string]interface{}); ok {
			return value["v"]
		}
		return nil
	}

	getLastRefetchTime := func(payload map[string]interface{}, userContext supertokens.UserContext) *int64 {
		if value, ok := payload[sessionClaim.Key].(map[string]interface{}); ok {
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

	validators := PrimitiveClaimValidators{
		HasValue: func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			if maxAgeInSeconds == nil {
				maxAgeInSeconds = defaultMaxAgeInSeconds
			}
			validatorId := sessionClaim.Key
			if id != nil {
				validatorId = *id
			}
			return SessionClaimValidator{
				ID:    validatorId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val := sessionClaim.GetValueFromPayload(payload, userContext)
					if val == nil {
						return true
					}
					return maxAgeInSeconds != nil && *getLastRefetchTime(payload, userContext) < time.Now().UnixNano()/1000000-*maxAgeInSeconds*1000
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal := sessionClaim.GetValueFromPayload(payload, userContext)

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":       "value does not exist",
								"expectedValue": val,
								"actualValue":   claimVal,
							},
						}
					}
					ageInSeconds := (time.Now().UnixNano()/1000000 - *getLastRefetchTime(payload, userContext)) / 1000
					if maxAgeInSeconds != nil && ageInSeconds > *maxAgeInSeconds {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":         "expired",
								"ageInSeconds":    ageInSeconds,
								"maxAgeInSeconds": *maxAgeInSeconds,
							},
						}
					}
					if claimVal != val {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":       "wrong value",
								"expectedValue": val,
								"actualValue":   claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
	}

	return sessionClaim, validators
}

type PrimitiveClaimValidators struct {
	HasValue func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
}
