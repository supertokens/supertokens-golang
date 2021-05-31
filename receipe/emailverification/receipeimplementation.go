package emailverification

import "github.com/supertokens/supertokens-golang/supertokens"

type RecipeImplementation struct {
	Querier supertokens.Querier
}

func NewRecipeImplementation(querier supertokens.Querier) *RecipeImplementation {
	return &RecipeImplementation{
		Querier: querier,
	}
}

func (r *RecipeImplementation) createEmailVerificationToken(userId, email string) {}

func (r *RecipeImplementation) verifyEmailUsingToken(token string) {}

func (r *RecipeImplementation) isEmailVerified(userId, email string) bool {
	return false
}
