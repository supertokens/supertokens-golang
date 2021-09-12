package sessmodels

type APIInterface struct {
	RefreshPOST   func(options APIOptions) error
	SignOutPOST   func(options APIOptions) (SignOutPOSTResponse, error)
	VerifySession func(verifySessionOptions *VerifySessionOptions, options APIOptions) (*SessionContainer, error)
}

type SignOutPOSTResponse struct {
	OK *struct{}
}
