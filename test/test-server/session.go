package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func addSessionRoutes(router *mux.Router) {
	router.HandleFunc("/test/session/createnewsessionwithoutrequestresponse", createNewSessionWithoutRequestResponse).Methods("POST")
	router.HandleFunc("/test/session/getsessionwithoutrequestresponse", getSessionWithoutRequestResponse).Methods("POST")
	router.HandleFunc("/test/session/getsessioninformation", getSessionInformation).Methods("POST")
	router.HandleFunc("/test/session/getallsessionhandlesforuser", getAllSessionHandlesForUser).Methods("POST")
	router.HandleFunc("/test/session/refreshsessionwithoutrequestresponse", refreshSessionWithoutRequestResponse).Methods("POST")
	router.HandleFunc("/test/session/revokeallsessionsforuser", revokeAllSessionsForUser).Methods("POST")
	router.HandleFunc("/test/session/mergeintoaccesspayload", mergeIntoAccessTokenPayload).Methods("POST")
	router.HandleFunc("/test/session/fetchandsetclaim", fetchAndSetClaim).Methods("POST")
	router.HandleFunc("/test/session/validateclaimsforsessionhandle", validateClaimsForSessionHandle).Methods("POST")

	router.HandleFunc("/test/session/sessionobject/revokesession", revokeSession).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/getsessiondatafromdatabase", getSessionDataFromDatabase).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/updatesessiondataindatabase", updateSessionDataInDatabase).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/getaccesstokenpayload", getAccessTokenPayload).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/gethandle", getHandle).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/getuserid", getUserId).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/gettenantid", getTenantId).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/getjwtpayload", getJWTPayload).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/getallsessiontokens", getAllSessionTokens).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/updateaccesstokenpayload", updateAccessTokenPayload).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/updatejwtpayload", updateJWTPayload).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/setclaimvalue", setClaimValue).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/getclaimvalue", getClaimValue).Methods("POST")
	// router.HandleFunc("/test/session/sessionobject/removeclaim", removeClaim).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/assertclaims", assertClaims).Methods("POST")
	router.HandleFunc("/test/session/sessionobject/mergeintoaccesstokenpayload", mergeIntoAccessTokenPayloadOnSessionObject).Methods("POST")
}

func createNewSessionWithoutRequestResponse(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TenantId              string                 `json:"tenantId"`
		UserId                string                 `json:"userId"`
		AccessTokenPayload    map[string]interface{} `json:"accessTokenPayload"`
		SessionDataInDatabase map[string]interface{} `json:"sessionDataInDatabase"`
		DisableAntiCSRF       *bool                  `json:"disableAntiCsrf"`
		UserContext           map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	sessionContainer, err := session.CreateNewSessionWithoutRequestResponse(
		body.TenantId,
		body.UserId,
		body.AccessTokenPayload,
		body.SessionDataInDatabase,
		body.DisableAntiCSRF,
		userContext,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessionHandle":         sessionContainer.GetHandle(),
		"userId":                sessionContainer.GetUserID(),
		"tenantId":              sessionContainer.GetTenantId(),
		"userDataInAccessToken": sessionContainer.GetAccessTokenPayload(),
		"accessToken":           sessionContainer.GetAccessToken(),
		"frontToken":            sessionContainer.GetAllSessionTokensDangerously().FrontToken,
		"refreshToken":          sessionContainer.GetAllSessionTokensDangerously().RefreshToken,
		"antiCsrfToken":         sessionContainer.GetAllSessionTokensDangerously().AntiCsrfToken,
		"accessTokenUpdated":    sessionContainer.GetAllSessionTokensDangerously().AccessAndFrontendTokenUpdated,
	})
}

func getSessionWithoutRequestResponse(w http.ResponseWriter, r *http.Request) {
	var body struct {
		AccessToken   string                           `json:"accessToken"`
		AntiCSRFToken *string                          `json:"antiCsrfToken"`
		Options       *sessmodels.VerifySessionOptions `json:"options"`
		UserContext   map[string]interface{}           `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	sessionContainer, err := session.GetSessionWithoutRequestResponse(
		body.AccessToken,
		body.AntiCSRFToken,
		body.Options,
		userContext,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessionContainer)
}

func getSessionInformation(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SessionHandle string                 `json:"sessionHandle"`
		UserContext   map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	sessionInfo, err := session.GetSessionInformation(body.SessionHandle, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessionInfo)
}

func getAllSessionHandlesForUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId                            string                 `json:"userId"`
		FetchSessionsForAllLinkedAccounts *bool                  `json:"fetchSessionsForAllLinkedAccounts"`
		TenantId                          *string                `json:"tenantId"`
		UserContext                       map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	sessionHandles, err := session.GetAllSessionHandlesForUser(body.UserId, body.TenantId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessionHandles)
}

func refreshSessionWithoutRequestResponse(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken    string                 `json:"refreshToken"`
		DisableAntiCSRF *bool                  `json:"disableAntiCsrf"`
		AntiCSRFToken   *string                `json:"antiCsrfToken"`
		UserContext     map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	sessionContainer, err := session.RefreshSessionWithoutRequestResponse(
		body.RefreshToken,
		body.DisableAntiCSRF,
		body.AntiCSRFToken,
		userContext,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(sessionContainer)
}

func revokeAllSessionsForUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId                          string                 `json:"userId"`
		RevokeSessionsForLinkedAccounts *bool                  `json:"revokeSessionsForLinkedAccounts"`
		TenantId                        *string                `json:"tenantId"`
		UserContext                     map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	revokedSessionHandles, err := session.RevokeAllSessionsForUser(body.UserId, body.TenantId, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(revokedSessionHandles)
}

func mergeIntoAccessTokenPayload(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SessionHandle            string                 `json:"sessionHandle"`
		AccessTokenPayloadUpdate map[string]interface{} `json:"accessTokenPayloadUpdate"`
		UserContext              map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	success, err := session.MergeIntoAccessTokenPayload(
		body.SessionHandle,
		body.AccessTokenPayloadUpdate,
		userContext,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(success)
}

func fetchAndSetClaim(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SessionHandle string                  `json:"sessionHandle"`
		Claim         claims.TypeSessionClaim `json:"claim"`
		UserContext   map[string]interface{}  `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	success, err := session.FetchAndSetClaim(
		body.SessionHandle,
		&body.Claim,
		userContext,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(success)
}

func validateClaimsForSessionHandle(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SessionHandle string                 `json:"sessionHandle"`
		UserContext   map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	response, err := session.ValidateClaimsForSessionHandle(
		body.SessionHandle,
		nil, // You might want to implement overrideGlobalClaimValidators if needed
		userContext,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func revokeSession(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session     map[string]interface{} `json:"session"`
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

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessionContainer.RevokeSessionWithContext(userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "OK",
	})
}

func getSessionDataFromDatabase(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session     map[string]interface{} `json:"session"`
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

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionData, err := sessionContainer.GetSessionDataInDatabaseWithContext(userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"sessionData": sessionData,
	})
}

func updateSessionDataInDatabase(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session        map[string]interface{} `json:"session"`
		NewSessionData map[string]interface{} `json:"newSessionData"`
		UserContext    map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessionContainer.UpdateSessionDataInDatabaseWithContext(body.NewSessionData, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "OK",
	})
}

func getAccessTokenPayload(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session     map[string]interface{} `json:"session"`
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

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload := sessionContainer.GetAccessTokenPayloadWithContext(userContext)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"payload": payload,
	})
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session     map[string]interface{} `json:"session"`
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

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handle := sessionContainer.GetHandleWithContext(userContext)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"handle": handle,
	})
}

func getUserId(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session     map[string]interface{} `json:"session"`
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

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := sessionContainer.GetUserIDWithContext(userContext)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId": userId,
	})
}

func getTenantId(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session     map[string]interface{} `json:"session"`
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

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenantId := sessionContainer.GetTenantIdWithContext(userContext)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tenantId": tenantId,
	})
}

func assertClaims(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session         map[string]interface{}   `json:"session"`
		ClaimValidators []map[string]interface{} `json:"claimValidators"`
		UserContext     map[string]interface{}   `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logOverrideEvent("sessionobject.assertclaims", "CALL", body)
	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	claimValidators := []claims.SessionClaimValidator{}
	for _, validator := range body.ClaimValidators {
		val := deserializeValidator(validator)
		if val != nil {
			args := validator["args"].([]interface{})
			claimValidators = append(claimValidators, val(args...))
		}
	}

	err = sessionContainer.AssertClaimsWithContext(claimValidators, userContext)

	if err == nil {
		logOverrideEvent("sessionobject.assertclaims", "RES", nil)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"updatedSession": map[string]interface{}{
				"sessionHandle":         sessionContainer.GetHandle(),
				"userId":                sessionContainer.GetUserID(),
				"tenantId":              sessionContainer.GetTenantId(),
				"userDataInAccessToken": sessionContainer.GetAccessTokenPayload(),
				"accessToken":           sessionContainer.GetAccessToken(),
				"frontToken":            sessionContainer.GetAllSessionTokensDangerously().FrontToken,
				"refreshToken":          sessionContainer.GetAllSessionTokensDangerously().RefreshToken,
				"antiCsrfToken":         sessionContainer.GetAllSessionTokensDangerously().AntiCsrfToken,
				"accessTokenUpdated":    sessionContainer.GetAllSessionTokensDangerously().AccessAndFrontendTokenUpdated,
			},
		})
	} else {
		logOverrideEvent("sessionobject.assertclaims", "REJ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mergeIntoAccessTokenPayloadOnSessionObject(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Session                  map[string]interface{} `json:"session"`
		AccessTokenPayloadUpdate map[string]interface{} `json:"accessTokenPayloadUpdate"`
		UserContext              map[string]interface{} `json:"userContext"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionContainer, err := convertMapToSessionContainer(body.Session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userContext supertokens.UserContext = nil
	if body.UserContext != nil {
		userContext = &body.UserContext
	}

	err = sessionContainer.MergeIntoAccessTokenPayloadWithContext(body.AccessTokenPayloadUpdate, userContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "OK",
		"updatedSession": map[string]interface{}{
			"sessionHandle":         sessionContainer.GetHandle(),
			"userId":                sessionContainer.GetUserID(),
			"tenantId":              sessionContainer.GetTenantId(),
			"userDataInAccessToken": sessionContainer.GetAccessTokenPayload(),
			"accessToken":           sessionContainer.GetAccessToken(),
			"frontToken":            sessionContainer.GetAllSessionTokensDangerously().FrontToken,
			"refreshToken":          sessionContainer.GetAllSessionTokensDangerously().RefreshToken,
			"antiCsrfToken":         sessionContainer.GetAllSessionTokensDangerously().AntiCsrfToken,
			"accessTokenUpdated":    sessionContainer.GetAllSessionTokensDangerously().AccessAndFrontendTokenUpdated,
		},
	})
}

func convertMapToSessionContainer(sessionMap map[string]interface{}) (sessmodels.SessionContainer, error) {
	return session.NewSessionContainerFromSessionContainerInputForTestServer(sessionMap)
}
