package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func PrimitiveClaim(key string, fetchValue FetchValueFunc) *TypePrimitiveClaim {
	sessionClaim := SessionClaim(key, fetchValue)

	sessionClaim.AddToPayload_internal = func(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{} {
		payload[sessionClaim.Key] = map[string]interface{}{
			"v": value,
			"t": time.Now().Unix(),
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

	primitiveClaim := &TypePrimitiveClaim{
		TypeSessionClaim: sessionClaim,
	}

	primitiveClaim.GetLastRefetchTime = func(payload map[string]interface{}, userContext supertokens.UserContext) int64 {
		if value, ok := payload[sessionClaim.Key].(map[string]interface{}); ok {
			return value["t"].(int64)
		}
		return 0
	}
	primitiveClaim.Validators = PrimitiveClaimValidators{
		HasValue: func(val interface{}, id *string) *SessionClaimValidator {
			validatorId := primitiveClaim.Key
			if id != nil {
				validatorId = *id
			}
			return &SessionClaimValidator{
				ID:    validatorId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val := primitiveClaim.GetValueFromPayload(payload, userContext)
					return val == nil
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal := primitiveClaim.GetValueFromPayload(payload, userContext)
					isValid := claimVal == val
					if isValid {
						return ClaimValidationResult{
							IsValid: true,
						}
					} else {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":       "wrong value",
								"expectedValue": val,
								"actualValue":   claimVal,
							},
						}
					}
				},
			}
		},
		HasFreshValue: func(val interface{}, maxAgeInSeconds int64, id *string) *SessionClaimValidator {
			validatorId := primitiveClaim.Key
			if id != nil {
				validatorId = *id
			}
			return &SessionClaimValidator{
				ID:    validatorId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).(map[string]interface{})
					if !ok || val == nil {
						return true
					}
					return val["t"].(int64) < time.Now().Unix()-maxAgeInSeconds
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
					ageInSeconds := time.Now().Unix() - primitiveClaim.GetLastRefetchTime(payload, userContext)
					if ageInSeconds > maxAgeInSeconds {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":         "expired",
								"ageInSeconds":    ageInSeconds,
								"maxAgeInSeconds": maxAgeInSeconds,
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

	return primitiveClaim
}

type TypePrimitiveClaim struct {
	*TypeSessionClaim
	GetLastRefetchTime func(payload map[string]interface{}, userContext supertokens.UserContext) int64
	Validators         PrimitiveClaimValidators
}

type PrimitiveClaimValidators struct {
	HasValue      func(val interface{}, id *string) *SessionClaimValidator
	HasFreshValue func(val interface{}, maxAgeInSeconds int64, id *string) *SessionClaimValidator
}
