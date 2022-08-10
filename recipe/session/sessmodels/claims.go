package sessmodels

type ValidateClaimsResponse struct {
	InvalidClaims            []ClaimValidationError
	AccessTokenPayloadUpdate map[string]interface{}
}

type ClaimValidationError struct {
}
