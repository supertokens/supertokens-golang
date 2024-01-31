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
			},
		},
	}

	jsonified, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "{\"id\":\"123\",\"timeJoined\":123,\"isPrimaryUser\":true,\"tenantIds\":[\"t1\",\"t2\"],\"emails\":[\"e1\",\"e2\"],\"phoneNumbers\":[\"p1\",\"p2\"],\"thirdParty\":[{\"id\":\"tp1\",\"userId\":\"tpu1\"},{\"id\":\"tp2\",\"userId\":\"tpu2\"}],\"loginMethods\":[{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"123\",\"recipeId\":\"emailpassword\",\"email\":\"e1\",\"phoneNumber\":\"p1\",\"thirdParty\":{\"id\":\"tp1\",\"userId\":\"tpu1\"},\"verified\":true},{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"haha\",\"recipeId\":\"emailpassword\",\"email\":\"e2\",\"phoneNumber\":\"p2\",\"thirdParty\":{\"id\":\"tp2\",\"userId\":\"tpu2\"},\"verified\":true}]}", string(jsonified))

	// now we test unmarshaling
	var user2 User
	err = json.Unmarshal(jsonified, &user2)
	if err != nil {
		t.Error(err)
		return
	}

	// check that user2 and user are the same
	assert.Equal(t, user, user2)

	// we marshall it once again
	jsonified, err = json.Marshal(user2)
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
			},
		},
	}

	jsonified, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "{\"id\":\"123\",\"timeJoined\":123,\"isPrimaryUser\":true,\"tenantIds\":[],\"emails\":[],\"phoneNumbers\":[],\"thirdParty\":[],\"loginMethods\":[{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"123\",\"recipeId\":\"emailpassword\",\"verified\":true},{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"456\",\"recipeId\":\"emailpassword\",\"verified\":true}]}", string(jsonified))

	// now we test unmarshaling
	var user2 User
	err = json.Unmarshal(jsonified, &user2)
	if err != nil {
		t.Error(err)
		return
	}

	// check that user2 and user are the same
	assert.Equal(t, user, user2)

	// we marshall it once again
	jsonified, err = json.Marshal(user2)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "{\"id\":\"123\",\"timeJoined\":123,\"isPrimaryUser\":true,\"tenantIds\":[],\"emails\":[],\"phoneNumbers\":[],\"thirdParty\":[],\"loginMethods\":[{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"123\",\"recipeId\":\"emailpassword\",\"verified\":true},{\"tenantIds\":[\"t1\",\"t2\"],\"timeJoined\":123,\"recipeUserId\":\"456\",\"recipeId\":\"emailpassword\",\"verified\":true}]}", string(jsonified))
}

func TestHelperFunctions(t *testing.T) {
	recipeUserId1, _ := NewRecipeUserID("123")
	recipeUserId2, _ := NewRecipeUserID("haha")
	email1 := "e1"
	phoneNumber1 := "+36701234123"
	thirdParty1 := ThirdParty{
		ID:     "tp1",
		UserID: "tpu1",
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
							Email:       nil,
							PhoneNumber: nil,
							ThirdParty:  nil,
						},
					},
				},
			},
		},
	}

	{
		email := "e1"
		assert.True(t, user.LoginMethods[0].HasSameEmailAs(&email))
	}
	{
		email := "E1"
		assert.True(t, user.LoginMethods[0].HasSameEmailAs(&email))
	}
	{
		email := "E1 "
		assert.True(t, user.LoginMethods[0].HasSameEmailAs(&email))
	}
	{
		email := "e"
		assert.False(t, user.LoginMethods[0].HasSameEmailAs(&email))
		assert.False(t, user.LoginMethods[0].HasSameEmailAs(nil))
	}
	{
		email := "e1"
		assert.False(t, user.LoginMethods[1].HasSameEmailAs(&email))
		assert.False(t, user.LoginMethods[1].HasSameEmailAs(nil))
	}

	{
		phoneNumber := "+36701234123"
		assert.True(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
	}
	{
		phoneNumber := " \t+36-70/1234 123  "
		assert.True(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
	}
	{
		phoneNumber := " \t+36701234123  "
		assert.True(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
	}
	{
		phoneNumber := "36701234123"
		assert.False(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
	}
	{
		phoneNumber := "0036701234123"
		assert.False(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
	}
	{
		phoneNumber := "06701234123"
		assert.False(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
	}
	{
		phoneNumber := "p36701234123"
		assert.False(t, user.LoginMethods[0].HasSamePhoneNumberAs(&phoneNumber))
		assert.False(t, user.LoginMethods[0].HasSamePhoneNumberAs(nil))
	}
	{
		phoneNumber := "+36701234123"
		assert.False(t, user.LoginMethods[1].HasSamePhoneNumberAs(&phoneNumber))
		assert.False(t, user.LoginMethods[1].HasSamePhoneNumberAs(nil))
	}
	{
		thirdParty := ThirdParty{
			ID:     "tp1",
			UserID: "tpu1",
		}
		assert.True(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(&thirdParty))
	}
	{
		thirdParty := ThirdParty{
			ID:     " tp1  ",
			UserID: "  tpu1\t",
		}
		assert.True(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(&thirdParty))
	}
	{
		thirdParty := ThirdParty{
			ID:     " tp1  ",
			UserID: "  Tpu1\t",
		}
		assert.False(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(&thirdParty))
	}

	{
		thirdParty := ThirdParty{
			ID:     " Tp1  ",
			UserID: "  tpu1\t",
		}
		assert.False(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(&thirdParty))
	}
	{
		thirdParty := ThirdParty{
			ID:     " tp12  ",
			UserID: "  tpU1\t",
		}
		assert.False(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(&thirdParty))
	}
	{
		thirdParty := ThirdParty{
			ID:     " tp1  ",
			UserID: "  tpU2\t",
		}
		assert.False(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(&thirdParty))
		assert.False(t, user.LoginMethods[0].HasSameThirdPartyInfoAs(nil))
	}
	{
		thirdParty := ThirdParty{
			ID:     "tp1",
			UserID: "tpu1",
		}
		assert.False(t, user.LoginMethods[1].HasSameThirdPartyInfoAs(&thirdParty))
		assert.False(t, user.LoginMethods[1].HasSameThirdPartyInfoAs(nil))
	}
}
