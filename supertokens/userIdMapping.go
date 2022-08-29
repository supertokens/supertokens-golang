package supertokens

import (
	"errors"
)

var userIdTypeDefault = "ANY"

type CreateUserIdMappingResult struct {
	OK                              *struct{}
	UnknownSupertokensUserIdError   *struct{}
	UserIdMappingAlreadyExistsError *struct {
		DoesSuperTokensUserIdExist bool
		DoesExternalUserIdExist    bool
	}
}

func CreateUserIdMapping(supertokensUserId string, externalUserId string, externalUserIdInfo *string, force bool) (CreateUserIdMappingResult, error) {
	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return CreateUserIdMappingResult{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return CreateUserIdMappingResult{}, err
	}
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return CreateUserIdMappingResult{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	data := map[string]interface{}{
		"superTokensUserId": supertokensUserId,
		"externalUserId":    externalUserId,
		"force":             force,
	}
	if externalUserIdInfo != nil {
		data["externalUserIdInfo"] = *externalUserIdInfo
	}
	resp, err := querier.SendPostRequest("/recipe/userid/map", data)
	if err != nil {
		return CreateUserIdMappingResult{}, err
	}
	if resp["status"] == "OK" {
		return CreateUserIdMappingResult{
			OK: &struct{}{},
		}, nil
	} else if resp["status"] == "UNKNOWN_SUPERTOKENS_USER_ID_ERROR" {
		return CreateUserIdMappingResult{
			UnknownSupertokensUserIdError: &struct{}{},
		}, nil
	} else {
		return CreateUserIdMappingResult{
			UserIdMappingAlreadyExistsError: &struct {
				DoesSuperTokensUserIdExist bool
				DoesExternalUserIdExist    bool
			}{
				DoesSuperTokensUserIdExist: resp["doesSuperTokensUserIdExist"].(bool),
				DoesExternalUserIdExist:    resp["doesExternalUserIdExist"].(bool),
			},
		}, nil
	}
}

type GetUserIdMappingResult struct {
	OK *struct {
		SupertokensUserId  string
		ExternalUserId     string
		ExternalUserIdInfo *string
	}
	UnknownMappingError *struct{}
}

func GetUserIdMapping(userId string, userIdType *string) (GetUserIdMappingResult, error) {
	if userIdType == nil {
		userIdType = &userIdTypeDefault
	}
	err := validateUserIdType(*userIdType)
	if err != nil {
		return GetUserIdMappingResult{}, err
	}

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return GetUserIdMappingResult{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return GetUserIdMappingResult{}, err
	}
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return GetUserIdMappingResult{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	resp, err := querier.SendGetRequest("/recipe/userid/map", map[string]string{
		"userId":     userId,
		"userIdType": *userIdType,
	})
	if err != nil {
		return GetUserIdMappingResult{}, err
	}
	if resp["status"] == "OK" {
		var externalUserIdInfo *string = nil
		if v, ok := resp["externalUserIdInfo"].(string); ok {
			externalUserIdInfo = &v
		}
		return GetUserIdMappingResult{
			OK: &struct {
				SupertokensUserId  string
				ExternalUserId     string
				ExternalUserIdInfo *string
			}{
				SupertokensUserId:  resp["superTokensUserId"].(string),
				ExternalUserId:     resp["externalUserId"].(string),
				ExternalUserIdInfo: externalUserIdInfo,
			},
		}, nil
	} else {
		return GetUserIdMappingResult{
			UnknownMappingError: &struct{}{},
		}, nil
	}
}

type DeleteUserIdMappingResult struct {
	OK *struct {
		DidMappingExist bool
	}
}

func DeleteUserIdMapping(userId string, userIdType *string, force bool) (DeleteUserIdMappingResult, error) {
	if userIdType == nil {
		userIdType = &userIdTypeDefault
	}
	err := validateUserIdType(*userIdType)
	if err != nil {
		return DeleteUserIdMappingResult{}, err
	}

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return DeleteUserIdMappingResult{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return DeleteUserIdMappingResult{}, err
	}
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return DeleteUserIdMappingResult{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	resp, err := querier.SendPostRequest("/recipe/userid/map/remove", map[string]interface{}{
		"userId":     userId,
		"userIdType": userIdType,
		"force":      force,
	})
	if err != nil {
		return DeleteUserIdMappingResult{}, err
	}
	return DeleteUserIdMappingResult{
		OK: &struct{ DidMappingExist bool }{
			DidMappingExist: resp["didMappingExist"].(bool),
		},
	}, nil
}

type UpdateOrDeleteUserIdMappingInfoResult struct {
	OK                  *struct{}
	UnknownMappingError *struct{}
}

func UpdateOrDeleteUserIdMappingInfo(userId string, userIdType *string, externalUserIdInfo *string) (UpdateOrDeleteUserIdMappingInfoResult, error) {
	if userIdType == nil {
		userIdType = &userIdTypeDefault
	}
	err := validateUserIdType(*userIdType)
	if err != nil {
		return UpdateOrDeleteUserIdMappingInfoResult{}, err
	}

	querier, err := GetNewQuerierInstanceOrThrowError("")
	if err != nil {
		return UpdateOrDeleteUserIdMappingInfoResult{}, err
	}
	cdiVersion, err := querier.GetQuerierAPIVersion()
	if err != nil {
		return UpdateOrDeleteUserIdMappingInfoResult{}, err
	}
	if maxVersion(cdiVersion, "2.15") != cdiVersion {
		return UpdateOrDeleteUserIdMappingInfoResult{}, errors.New("Please upgrade the SuperTokens core to >= 3.15.0")
	}

	resp, err := querier.SendPutRequest("/recipe/userid/external-user-id-info", map[string]interface{}{
		"userId":             userId,
		"userIdType":         *userIdType,
		"externalUserIdInfo": externalUserIdInfo,
	})
	if err != nil {
		return UpdateOrDeleteUserIdMappingInfoResult{}, err
	}

	if resp["status"] == "OK" {
		return UpdateOrDeleteUserIdMappingInfoResult{
			OK: &struct{}{},
		}, nil
	} else {
		return UpdateOrDeleteUserIdMappingInfoResult{
			UnknownMappingError: &struct{}{},
		}, nil
	}
}

func validateUserIdType(userIdType string) error {
	if userIdType != "SUPERTOKENS" && userIdType != "EXTERNAL" && userIdType != "ANY" {
		return errors.New("userIdType must be one of SUPERTOKENS, EXTERNAL or ANY")
	}
	return nil
}
