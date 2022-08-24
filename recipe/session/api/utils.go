package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func getRequiredClaimValidators(
	claimValidatorsAddedByOtherRecipes []*claims.SessionClaimValidator,
	sessionRecipeImpl sessmodels.RecipeInterface,
	sessionContainer *sessmodels.SessionContainer,
	overrideGlobalClaimValidators func(globalClaimValidators []*claims.SessionClaimValidator, sessionContainer *sessmodels.SessionContainer, userContext supertokens.UserContext) []*claims.SessionClaimValidator,
	userContext supertokens.UserContext,
) ([]*claims.SessionClaimValidator, error) {
	globalClaimValidators, err := (*sessionRecipeImpl.GetGlobalClaimValidators)(sessionContainer.GetUserID(), claimValidatorsAddedByOtherRecipes, userContext)
	if err != nil {
		return nil, err
	}
	if overrideGlobalClaimValidators != nil {
		globalClaimValidators = overrideGlobalClaimValidators(globalClaimValidators, sessionContainer, userContext)
	}
	return globalClaimValidators, nil
}
