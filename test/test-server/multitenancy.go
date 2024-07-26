package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy"
	"github.com/supertokens/supertokens-golang/recipe/multitenancy/multitenancymodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func addMultitenancyRoutes(router *mux.Router) {
	router.HandleFunc("/test/multitenancy/createorupdatetenant", createOrUpdateTenantHandler).Methods("POST")
	router.HandleFunc("/test/multitenancy/deletetenant", deleteTenantHandler).Methods("POST")
	router.HandleFunc("/test/multitenancy/gettenant", getTenantHandler).Methods("POST")
	router.HandleFunc("/test/multitenancy/listalltenants", listAllTenantsHandler).Methods("GET")
	router.HandleFunc("/test/multitenancy/createorupdatethirdpartyconfig", createOrUpdateThirdPartyConfigHandler).Methods("POST")
	router.HandleFunc("/test/multitenancy/deletethirdpartyconfig", deleteThirdPartyConfigHandler).Methods("POST")
	router.HandleFunc("/test/multitenancy/associateusertotenant", associateUserToTenantHandler).Methods("POST")
	router.HandleFunc("/test/multitenancy/disassociateuserfromtenant", disassociateUserFromTenantHandler).Methods("POST")
}

func createOrUpdateTenantHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId    string                          `json:"tenantId"`
		Config      multitenancymodels.TenantConfig `json:"config"`
		UserContext map[string]interface{}          `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	response, err := multitenancy.CreateOrUpdateTenant(body.TenantId, body.Config, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func deleteTenantHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId    string                 `json:"tenantId"`
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

	response, err := multitenancy.DeleteTenant(body.TenantId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func getTenantHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId    string                 `json:"tenantId"`
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

	tenant, err := multitenancy.GetTenant(body.TenantId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}

func listAllTenantsHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
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

	response, err := multitenancy.ListAllTenants(userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func createOrUpdateThirdPartyConfigHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId       string                  `json:"tenantId"`
		Config         tpmodels.ProviderConfig `json:"config"`
		SkipValidation *bool                   `json:"skipValidation"`
		UserContext    map[string]interface{}  `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	response, err := multitenancy.CreateOrUpdateThirdPartyConfig(body.TenantId, body.Config, body.SkipValidation, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func deleteThirdPartyConfigHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId     string                 `json:"tenantId"`
		ThirdPartyId string                 `json:"thirdPartyId"`
		UserContext  map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	response, err := multitenancy.DeleteThirdPartyConfig(body.TenantId, body.ThirdPartyId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func associateUserToTenantHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId    string                 `json:"tenantId"`
		UserId      string                 `json:"userId"`
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

	response, err := multitenancy.AssociateUserToTenant(body.TenantId, body.UserId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func disassociateUserFromTenantHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId    string                 `json:"tenantId"`
		UserId      string                 `json:"userId"`
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

	response, err := multitenancy.DisassociateUserFromTenant(body.TenantId, body.UserId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}
