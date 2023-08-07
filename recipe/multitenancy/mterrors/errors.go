package mterrors

type TenantDoesNotExistError struct {
	Msg string
}

func (err TenantDoesNotExistError) Error() string {
	return err.Msg
}

type RecipeDisabledForTenantError struct {
	Msg string
}

func (err RecipeDisabledForTenantError) Error() string {
	return err.Msg
}
