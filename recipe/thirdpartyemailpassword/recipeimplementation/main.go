package recipeimplementation

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	epm "github.com/supertokens/supertokens-golang/recipe/emailpassword/models"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(emailPasswordQuerier supertokens.Querier, thirdPartyQuerier *supertokens.Querier) models.RecipeInterface {
	emailPasswordImplementation := emailpassword.MakeRecipeImplementation(emailPasswordQuerier)
	var thirdPartyImplementation *tpmodels.RecipeInterface
	if thirdPartyQuerier != nil {
		thirdPartyImplementationTemp := thirdparty.MakeRecipeImplementation(*thirdPartyQuerier)
		thirdPartyImplementation = &thirdPartyImplementationTemp
	}
	return models.RecipeInterface{
		SignUp: func(email, password string) (models.SignUpResponse, error) {
			response, err := emailPasswordImplementation.SignUp(email, password)
			if err != nil {
				return models.SignUpResponse{}, err
			}
			if response.EmailAlreadyExistsError != nil {
				return models.SignUpResponse{
					EmailAlreadyExistsError: &struct{}{},
				}, nil
			}
			return models.SignUpResponse{
				OK: &struct{ User models.User }{
					User: models.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
						ThirdParty: nil,
					},
				},
			}, nil
		},

		SignIn: func(email, password string) (models.SignInResponse, error) {
			response, err := emailPasswordImplementation.SignIn(email, password)
			if err != nil {
				return models.SignInResponse{}, err
			}
			if response.WrongCredentialsError != nil {
				return models.SignInResponse{
					WrongCredentialsError: &struct{}{},
				}, nil
			}
			return models.SignInResponse{
				OK: &struct{ User models.User }{
					User: models.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
						ThirdParty: nil,
					},
				},
			}, nil
		},

		SignInUp: func(thirdPartyID, thirdPartyUserID string, email models.EmailStruct) (models.SignInUpResponse, error) {
			if thirdPartyImplementation == nil {
				return models.SignInUpResponse{}, errors.New("no thirdparty provider configured")
			}
			result, err := (*thirdPartyImplementation).SignInUp(thirdPartyID, thirdPartyUserID, tpmodels.EmailStruct{
				ID:         email.ID,
				IsVerified: email.IsVerified,
			})
			if err != nil {
				return models.SignInUpResponse{}, err
			}
			if result.FieldError != nil {
				return models.SignInUpResponse{
					FieldError: &struct{ Error string }{
						Error: result.FieldError.Error,
					},
				}, nil
			}
			return models.SignInUpResponse{
				OK: &struct {
					CreatedNewUser bool
					User           models.User
				}{
					CreatedNewUser: result.OK.CreatedNewUser,
					User: models.User{
						ID:         result.OK.User.ID,
						Email:      result.OK.User.Email,
						TimeJoined: result.OK.User.TimeJoined,
						ThirdParty: &result.OK.User.ThirdParty,
					},
				},
			}, nil
		},

		GetUserByID: func(userID string) (*models.User, error) {
			user, err := emailPasswordImplementation.GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			if user != nil {
				return &models.User{
					ID:         user.ID,
					Email:      user.Email,
					TimeJoined: user.TimeJoined,
					ThirdParty: nil,
				}, nil
			}
			if thirdPartyImplementation == nil {
				return nil, nil
			}

			userinfo, err := thirdPartyImplementation.GetUserByID(userID)
			if err != nil {
				return nil, err
			}

			if userinfo != nil {
				return &models.User{
					ID:         userinfo.ID,
					Email:      userinfo.Email,
					TimeJoined: userinfo.TimeJoined,
					ThirdParty: &userinfo.ThirdParty,
				}, nil
			}
			return nil, nil
		},

		GetUsersByEmail: func(email string) ([]models.User, error) {
			fromEP, err := emailPasswordImplementation.GetUserByEmail(email)
			if err != nil {
				return []models.User{}, err
			}

			fromTP := []tpmodels.User{}
			if thirdPartyImplementation != nil {
				fromTP, err = (*thirdPartyImplementation).GetUsersByEmail(email)
				if err != nil {
					return []models.User{}, err
				}
			}
			finalResult := []models.User{}

			if fromEP != nil {
				finalResult = append(finalResult, models.User{
					ID:         fromEP.ID,
					TimeJoined: fromEP.TimeJoined,
					Email:      fromEP.Email,
					ThirdParty: nil,
				})
			}

			for _, tpUser := range fromTP {
				finalResult = append(finalResult, models.User{
					ID:         tpUser.ID,
					TimeJoined: tpUser.TimeJoined,
					Email:      tpUser.Email,
					ThirdParty: &tpUser.ThirdParty,
				})
			}

			return finalResult, nil
		},

		GetUserByThirdPartyInfo: func(thirdPartyID string, thirdPartyUserID string) (*models.User, error) {
			if thirdPartyImplementation == nil {
				return nil, nil
			}

			userinfo, err := thirdPartyImplementation.GetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID)
			if err != nil {
				return nil, err
			}

			if userinfo != nil {
				return &models.User{
					ID:         userinfo.ID,
					Email:      userinfo.Email,
					TimeJoined: userinfo.TimeJoined,
					ThirdParty: &userinfo.ThirdParty,
				}, nil
			}
			return nil, nil
		},

		CreateResetPasswordToken: func(userID string) (epm.CreateResetPasswordTokenResponse, error) {
			return emailPasswordImplementation.CreateResetPasswordToken(userID)
		},
		ResetPasswordUsingToken: func(token, newPassword string) (epm.ResetPasswordUsingTokenResponse, error) {
			return emailPasswordImplementation.ResetPasswordUsingToken(token, newPassword)
		},
		UpdateEmailOrPassword: func(userId string, email, password *string) (epm.UpdateEmailOrPasswordResponse, error) {
			return emailPasswordImplementation.UpdateEmailOrPassword(userId, email, password)
		},
	}
}
