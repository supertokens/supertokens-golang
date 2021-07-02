package constants

const (
	FormFieldPasswordID           = "password"
	FormFieldEmailID              = "email"
	SignUpAPI                     = "/signup"
	SignInAPI                     = "/signin"
	GeneratePasswordResetTokenAPI = "/user/password/reset/token"
	PasswordResetAPI              = "/user/password/reset"
	SignupEmailExistsAPI          = "/signup/email/exists"

	EmailAlreadyExistsError        = "EMAIL_ALREADY_EXISTS_ERROR"
	WrongCredentialsError          = "WRONG_CREDENTIALS_ERROR"
	UnknownUserID                  = "UNKNOWN_USER_ID"
	ResetPasswordInvalidTokenError = "RESET_PASSWORD_INVALID_TOKEN_ERROR"
)
