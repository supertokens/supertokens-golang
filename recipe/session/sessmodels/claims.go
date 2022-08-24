package sessmodels

import "github.com/supertokens/supertokens-golang/recipe/session/claims"

type ValidateClaimsResult struct {
	InvalidClaims            []claims.ClaimValidationError
	AccessTokenPayloadUpdate map[string]interface{}
}

type ValidateClaimsResponse struct {
	OK *struct {
		InvalidClaims []claims.ClaimValidationError
	}
	SessionDoesNotExistError *struct{}
}

type GetClaimValueResult struct {
	OK *struct {
		Value interface{}
	}
	SessionDoesNotExistError *struct{}
}
