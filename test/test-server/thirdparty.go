package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func addThirdPartyRoutes(router *mux.Router) {
	router.HandleFunc("/test/thirdparty/manuallycreateorupdateuser", manuallyCreateOrUpdateUserHandler).Methods("POST")
	router.HandleFunc("/test/thirdparty/getuserbyid", getUserByIDHandler).Methods("POST")
	router.HandleFunc("/test/thirdparty/getusersbyemail", getUsersByEmailHandler).Methods("POST")
	router.HandleFunc("/test/thirdparty/getuserbythirdpartyinfo", getUserByThirdPartyInfoHandler).Methods("POST")
}

func manuallyCreateOrUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId         string                 `json:"tenantId"`
		ThirdPartyID     string                 `json:"thirdPartyId"`
		ThirdPartyUserID string                 `json:"thirdPartyUserId"`
		Email            string                 `json:"email"`
		UserContext      map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.TenantId == "" {
		body.TenantId = "public"
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	response, err := thirdparty.ManuallyCreateOrUpdateUser(body.TenantId, body.ThirdPartyID, body.ThirdPartyUserID, body.Email, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID      string                 `json:"userId"`
		UserContext map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	user, err := thirdparty.GetUserByID(body.UserID, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func getUsersByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId    string                 `json:"tenantId"`
		Email       string                 `json:"email"`
		UserContext map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.TenantId == "" {
		body.TenantId = "public"
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	users, err := thirdparty.GetUsersByEmail(body.TenantId, body.Email, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func getUserByThirdPartyInfoHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId         string                 `json:"tenantId"`
		ThirdPartyID     string                 `json:"thirdPartyId"`
		ThirdPartyUserID string                 `json:"thirdPartyUserId"`
		UserContext      map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.TenantId == "" {
		body.TenantId = "public"
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	user, err := thirdparty.GetUserByThirdPartyInfo(body.TenantId, body.ThirdPartyID, body.ThirdPartyUserID, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
