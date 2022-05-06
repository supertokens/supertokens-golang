package emaildelivery

import (
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type Ingredient struct {
	IngredientInterfaceImpl emaildeliverymodels.EmailDeliveryInterface
}

func MakeIngredient(config emaildeliverymodels.TypeInputWithService) Ingredient {
	defaultSendEmail := func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) error {
		return (*config.Service.SendEmail)(input, userContext)
	}

	result := Ingredient{
		IngredientInterfaceImpl: emaildeliverymodels.EmailDeliveryInterface{
			SendEmail: &defaultSendEmail,
		},
	}

	if config.Override != nil {
		result.IngredientInterfaceImpl = config.Override(result.IngredientInterfaceImpl)
	}

	return result
}
