package claims

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SessionClaim interface {
	FetchValue(userId string, userContext supertokens.UserContext) interface{}
	AddToPayload_internal(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{}
	RemoveFromPayloadByMerge_internal(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{}
	RemoveFromPayload(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{}
	GetValueFromPayload(payload map[string]interface{}, userContext supertokens.UserContext) interface{}
	Build(userId string, userContext supertokens.UserContext) map[string]interface{}
}

func BuildSessionClaim(c SessionClaim, userId string, userContext supertokens.UserContext) map[string]interface{} {
	value := c.FetchValue(userId, userContext)
	if value == nil {
		return map[string]interface{}{}
	}

	return c.AddToPayload_internal(map[string]interface{}{}, value, userContext)
}

type SessionClaimValidator interface {
	GetID() string
	GetClaim() SessionClaim
	ShouldRefetch(payload map[string]interface{}, userContext supertokens.UserContext) bool
	Validate(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult
}

type ClaimValidationResult struct {
	IsValid bool
	Reason  interface{}
}

type ClaimValidationError struct {
	ID     string
	Reason interface{}
}
