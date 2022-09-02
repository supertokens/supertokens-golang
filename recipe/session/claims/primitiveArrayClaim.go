package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func PrimitiveArrayClaim(key string, fetchValue FetchValueFunc, defaultMaxAgeInSeconds *int64) (*TypeSessionClaim, PrimitiveArrayClaimValidators) {
	if defaultMaxAgeInSeconds == nil {
		val := int64(300)
		defaultMaxAgeInSeconds = &val
	}

	sessionClaim := SessionClaim(key, fetchValue)

	sessionClaim.AddToPayload_internal = func(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{} {
		payload[sessionClaim.Key] = map[string]interface{}{
			"v": value.([]interface{}),
			"t": time.Now().UnixMilli(),
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

	validators := PrimitiveArrayClaimValidators{
		Includes: func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			if maxAgeInSeconds == nil {
				maxAgeInSeconds = defaultMaxAgeInSeconds
			}
			claimId := sessionClaim.Key
			if id != nil {
				claimId = *id
			}
			return SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					claimVal, ok := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					if !ok || claimVal == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return *getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*maxAgeInSeconds*1000
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "value does not exist",
								"expectedToInclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					ageInSeconds := (time.Now().UnixMilli() - *getLastRefetchTime(payload, userContext)) / 1000
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
					if !includes(claimVal, val) {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "wrong value",
								"expectedToInclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
		Excludes: func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			if maxAgeInSeconds == nil {
				maxAgeInSeconds = defaultMaxAgeInSeconds
			}
			claimId := sessionClaim.Key
			if id != nil {
				claimId = *id
			}
			return SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					if !ok || val == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return *getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*maxAgeInSeconds*1000
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "value does not exist",
								"expectedToInclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					ageInSeconds := (time.Now().UnixMilli() - *getLastRefetchTime(payload, userContext)) / 1000
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
					if includes(claimVal, val) {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "wrong value",
								"expectedToExclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
		IncludesAll: func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			if maxAgeInSeconds == nil {
				maxAgeInSeconds = defaultMaxAgeInSeconds
			}
			claimId := sessionClaim.Key
			if id != nil {
				claimId = *id
			}
			return SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					if !ok || val == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return *getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*maxAgeInSeconds*1000
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "value does not exist",
								"expectedToInclude": vals,
								"actualValue":       claimVal,
							},
						}
					}
					ageInSeconds := (time.Now().UnixMilli() - *getLastRefetchTime(payload, userContext)) / 1000
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

					if !includesAll(claimVal, vals) {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "wrong value",
								"expectedToInclude": vals,
								"actualValue":       claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
		ExcludesAll: func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator {
			if maxAgeInSeconds == nil {
				maxAgeInSeconds = defaultMaxAgeInSeconds
			}
			claimId := sessionClaim.Key
			if id != nil {
				claimId = *id
			}
			return SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					if !ok || val == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return *getLastRefetchTime(payload, userContext) < time.Now().UnixMilli()-*maxAgeInSeconds*1000
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":              "value does not exist",
								"expectedToNotInclude": vals,
								"actualValue":          claimVal,
							},
						}
					}
					ageInSeconds := (time.Now().UnixMilli() - *getLastRefetchTime(payload, userContext)) / 1000
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

					if !excludesAll(claimVal, vals) {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":              "wrong value",
								"expectedToNotInclude": vals,
								"actualValue":          claimVal,
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

type PrimitiveArrayClaimValidators struct {
	Includes    func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	Excludes    func(val interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	IncludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
	ExcludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) SessionClaimValidator
}
