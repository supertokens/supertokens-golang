package session

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func TrueClaim() (*claims.TypeSessionClaim, claims.BooleanClaimValidators) {
	claim, validators := claims.BooleanClaim(
		"st-true",
		func(userId string, tenantId *string, userContext supertokens.UserContext) (interface{}, error) {
			return true, nil
		},
		nil,
	)
	return claim, validators
}

func NilClaim() (*claims.TypeSessionClaim, claims.PrimitiveClaimValidators) {
	claim, validators := claims.PrimitiveClaim(
		"st-nil",
		func(userId string, tenantId *string, userContext supertokens.UserContext) (interface{}, error) {
			return nil, nil
		},
		nil,
	)
	return claim, validators
}

func StubClaim(validate func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult) (*claims.TypeSessionClaim, StubValidator) {
	claim, validators := claims.PrimitiveClaim(
		"st-stub",
		func(userId string, tenantId *string, userContext supertokens.UserContext) (interface{}, error) {
			return "stub", nil
		},
		nil,
	)

	return claim, StubValidator{
		PrimitiveClaimValidators: validators,
		Stub: func() claims.SessionClaimValidator {
			return claims.SessionClaimValidator{
				ID:    claim.Key,
				Claim: claim,
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
					return validate(payload, userContext)
				},
			}
		},
	}
}

func StubClaimWithRefetch(validate func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult) (*claims.TypeSessionClaim, StubValidator) {
	claim, validators := claims.PrimitiveClaim(
		"st-stub",
		func(userId string, tenantId *string, userContext supertokens.UserContext) (interface{}, error) {
			return "stub", nil
		},
		nil,
	)

	return claim, StubValidator{
		PrimitiveClaimValidators: validators,
		Stub: func() claims.SessionClaimValidator {
			return claims.SessionClaimValidator{
				ID:    claim.Key,
				Claim: claim,
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) claims.ClaimValidationResult {
					return validate(payload, userContext)
				},
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					return true
				},
			}
		},
	}
}

type StubValidator struct {
	claims.PrimitiveClaimValidators
	Stub func() claims.SessionClaimValidator
}
