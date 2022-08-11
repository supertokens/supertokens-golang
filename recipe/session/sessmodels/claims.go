package sessmodels

type ValidateClaimsResult struct {
	InvalidClaims            []ClaimValidationError
	AccessTokenPayloadUpdate map[string]interface{}
}

type ValidateClaimsResponse struct {
	OK *struct {
		InvalidClaims []ClaimValidationError
	}
	SessionDoesNotExistError *struct{}
}

type ClaimValidationError struct {
	ID     string
	Reason interface{}
}
