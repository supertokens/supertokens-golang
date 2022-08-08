package claims

import (
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Any interface{}

type SessionClaim interface {
	FetchValue(userId string, userContext supertokens.UserContext) Any
	AddToPayload_internal(payload map[string]Any, value Any, userContext supertokens.UserContext) map[string]Any
	RemoveFromPayloadByMerge_internal(payload map[string]Any, userContext supertokens.UserContext) map[string]Any
	RemoveFromPayload(payload map[string]Any, userContext supertokens.UserContext) map[string]Any
	GetValueFromPayload(payload map[string]Any, userContext supertokens.UserContext) Any
}

func BuildSessionClaim(c SessionClaim, userId string, userContext supertokens.UserContext) map[string]Any {
	value := c.FetchValue(userId, userContext)
	if value == nil {
		return map[string]Any{}
	}

	return c.AddToPayload_internal(map[string]Any{}, value, userContext)
}

type SessionClaimValidator interface {
	GetID() string
	GetClaim() SessionClaim
	ShouldRefetch(payload map[string]Any, userContext supertokens.UserContext) bool
	Validate(payload map[string]Any, userContext supertokens.UserContext) ClaimValidationResult
}

type ClaimValidationResult struct {
	IsValid bool
	Reason  Any
}

type ClaimValidationError struct {
	ID     string
	Reason Any
}
