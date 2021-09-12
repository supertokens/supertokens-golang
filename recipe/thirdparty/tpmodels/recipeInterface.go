package tpmodels

type RecipeInterface struct {
	GetUserByID             func(userID string) (*User, error)
	GetUsersByEmail         func(email string) ([]User, error)
	GetUserByThirdPartyInfo func(thirdPartyID string, thirdPartyUserID string) (*User, error)
	SignInUp                func(thirdPartyID string, thirdPartyUserID string, email EmailStruct) (SignInUpResponse, error)
}

type SignInUpResponse struct {
	OK *struct {
		CreatedNewUser bool
		User           User
	}
	FieldError *struct{ Error string }
}
