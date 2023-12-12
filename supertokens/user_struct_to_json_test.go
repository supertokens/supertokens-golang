package supertokens

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserJSONStruct(t *testing.T) {
	recipeUserId1, _ := NewRecipeUserID("123")
	recipeUserId2, _ := NewRecipeUserID("haha")
	email1 := "e1"
	phoneNumber1 := "p1"
	thirdParty1 := ThirdParty{
		ID:     "tp1",
		UserID: "tpu1",
	}
	email2 := "e2"
	phoneNumber2 := "p2"
	thirdParty2 := ThirdParty{
		ID:     "tp2",
		UserID: "tpu2",
	}

	user := User{
		ID:            "123",
		TimeJoined:    123,
		IsPrimaryUser: true,
		TenantIDs:     []string{"t1", "t2"},
		Emails:        []string{"e1", "e2"},
		PhoneNumbers:  []string{"p1", "p2"},
		ThirdParty: []ThirdParty{
			{
				ID:     "tp1",
				UserID: "tpu1",
			},
			{
				ID:     "tp2",
				UserID: "tpu2",
			},
		},
		LoginMethods: []LoginMethods{
			{
				Verified: true,
				RecipeLevelUser: RecipeLevelUser{
					TenantIDs:    []string{"t1", "t2"},
					TimeJoined:   123,
					RecipeUserID: recipeUserId1,
					AccountInfoWithRecipeID: AccountInfoWithRecipeID{
						RecipeID: EmailPasswordRID,
						AccountInfo: AccountInfo{
							Email:       &email1,
							PhoneNumber: &phoneNumber1,
							ThirdParty:  &thirdParty1,
						},
					},
				},
				HasSameEmailAs: func(email string) bool {
					return true
				},
				HasSamePhoneNumberAs: func(phoneNumber string) bool {
					return true
				},
				HasSameThirdPartyInfoAs: func(thirdParty *ThirdParty) bool {
					return true
				},
			},
			{
				Verified: true,
				RecipeLevelUser: RecipeLevelUser{
					TenantIDs:    []string{"t1", "t2"},
					TimeJoined:   123,
					RecipeUserID: recipeUserId2,
					AccountInfoWithRecipeID: AccountInfoWithRecipeID{
						RecipeID: EmailPasswordRID,
						AccountInfo: AccountInfo{
							Email:       &email2,
							PhoneNumber: &phoneNumber2,
							ThirdParty:  &thirdParty2,
						},
					},
				},
				HasSameEmailAs: func(email string) bool {
					return true
				},
				HasSamePhoneNumberAs: func(phoneNumber string) bool {
					return true
				},
				HasSameThirdPartyInfoAs: func(thirdParty *ThirdParty) bool {
					return true
				},
			},
		},
	}

	jsonified, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "{\"id\":\"123\",\"timeJoined\":123,\"isPrimaryUser\":true,\"tenantIds\":[\"t1\",\"t2\"],\"emails\":[\"e1\",\"e2\"],\"phoneNumbers\":[\"p1\",\"p2\"],\"thirdParty\":[{\"id\":\"tp1\",\"userId\":\"tpu1\"},{\"id\":\"tp2\",\"userId\":\"tpu2\"}],\"loginMethods\":[{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"123\",\"recipeId\":\"emailpassword\",\"email\":\"e1\",\"phoneNumber\":\"p1\",\"thirdParty\":{\"id\":\"tp1\",\"userId\":\"tpu1\"},\"verified\":true},{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"haha\",\"recipeId\":\"emailpassword\",\"email\":\"e2\",\"phoneNumber\":\"p2\",\"thirdParty\":{\"id\":\"tp2\",\"userId\":\"tpu2\"},\"verified\":true}]}", string(jsonified))
}

func TestUserJSONStructWithMissingAccountInfo(t *testing.T) {
	recipeUserId1, _ := NewRecipeUserID("123")
	recipeUserId2, _ := NewRecipeUserID("456")
	user := User{
		ID:            "123",
		TimeJoined:    123,
		IsPrimaryUser: true,
		TenantIDs:     []string{},
		Emails:        []string{},
		PhoneNumbers:  []string{},
		ThirdParty:    []ThirdParty{},
		LoginMethods: []LoginMethods{
			{
				Verified: true,
				RecipeLevelUser: RecipeLevelUser{
					TenantIDs:    []string{"t1", "t2"},
					TimeJoined:   123,
					RecipeUserID: recipeUserId1,
					AccountInfoWithRecipeID: AccountInfoWithRecipeID{
						RecipeID: EmailPasswordRID,
						AccountInfo: AccountInfo{
							Email: nil,
						},
					},
				},
			},
			{
				Verified: true,
				RecipeLevelUser: RecipeLevelUser{
					TenantIDs:    []string{"t1", "t2"},
					TimeJoined:   123,
					RecipeUserID: recipeUserId2,
					AccountInfoWithRecipeID: AccountInfoWithRecipeID{
						RecipeID: EmailPasswordRID,
						AccountInfo: AccountInfo{
							PhoneNumber: nil,
						},
					},
				},
				HasSameEmailAs: func(email string) bool {
					return true
				},
				HasSamePhoneNumberAs: func(phoneNumber string) bool {
					return true
				},
				HasSameThirdPartyInfoAs: func(thirdParty *ThirdParty) bool {
					return true
				},
			},
		},
	}

	jsonified, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "{\"id\":\"123\",\"timeJoined\":123,\"isPrimaryUser\":true,\"tenantIds\":[],\"emails\":[],\"phoneNumbers\":[],\"thirdParty\":[],\"loginMethods\":[{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"123\",\"recipeId\":\"emailpassword\",\"verified\":true},{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"456\",\"recipeId\":\"emailpassword\",\"verified\":true}]}", string(jsonified))
}
