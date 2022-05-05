package emaildelivery

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/edmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Ingredient struct {
	IngredientInterfaceImpl edmodels.EmailDeliveryInterface
}

func MakeIngredeint(config edmodels.TypeInputWithService) Ingredient {
	defaultSendEmail := func(input edmodels.EmailType, userContext supertokens.UserContext) error {
		return (*config.Service.SendEmail)(input, userContext)
	}

	result := Ingredient{
		IngredientInterfaceImpl: edmodels.EmailDeliveryInterface{
			SendEmail: &defaultSendEmail,
		},
	}

	if config.Override != nil {
		result.IngredientInterfaceImpl = (*config.Override)(result.IngredientInterfaceImpl)
	}

	return result
}
