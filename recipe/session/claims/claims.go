package claims

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SessionClaim(key string, fetchValue FetchValueFunc) *TypeSessionClaim {
	sessionClaim := &TypeSessionClaim{
		Key:        key,
		FetchValue: fetchValue,
	}

	sessionClaim.Build = func(userId string, userContext supertokens.UserContext) (map[string]interface{}, error) {
		value, err := sessionClaim.FetchValue(userId, userContext)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return map[string]interface{}{}, nil
		}
		return sessionClaim.AddToPayload_internal(map[string]interface{}{}, value, userContext), nil
	}

	return sessionClaim
}

type FetchValueFunc func(userId string, userContext supertokens.UserContext) (interface{}, error)

type TypeSessionClaim struct {
	Key                               string
	FetchValue                        FetchValueFunc
	AddToPayload_internal             func(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{}
	RemoveFromPayloadByMerge_internal func(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{}
	RemoveFromPayload                 func(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{}
	GetValueFromPayload               func(payload map[string]interface{}, userContext supertokens.UserContext) interface{}
	Build                             func(userId string, userContext supertokens.UserContext) (map[string]interface{}, error)
}

type SessionClaimValidator struct {
	ID            string
	Claim         *TypeSessionClaim
	ShouldRefetch func(payload map[string]interface{}, userContext supertokens.UserContext) bool
	Validate      func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult
}

type ClaimValidationResult struct {
	IsValid bool
	Reason  interface{} // This can be nil, add checks when used
}

type ClaimValidationError struct {
	ID     string      `json:"id"`
	Reason interface{} `json:"reason"` // This can be nil, add checks when used
}
