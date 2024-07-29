package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func addEmailPasswordRoutes(router *mux.Router) {
	router.HandleFunc("/test/emailpassword/signup", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			TenantId    string                 `json:"tenantId"`
			Email       string                 `json:"email"`
			Password    string                 `json:"password"`
			UserContext map[string]interface{} `json:"userContext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if input.TenantId == "" {
			input.TenantId = "public"
		}

		var userContext supertokens.UserContext = nil
		if input.UserContext != nil {
			userContext = &input.UserContext
		}

		response, err := emailpassword.SignUp(input.TenantId, input.Email, input.Password, userContext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jsonResponse map[string]interface{}
		if response.OK != nil {
			jsonResponse = map[string]interface{}{
				"status": "OK",
				"user":   response.OK.User,
			}
		} else {
			jsonResponse = map[string]interface{}{
				"status": "EMAIL_ALREADY_EXISTS_ERROR",
			}
		}
		json.NewEncoder(w).Encode(jsonResponse)
	}).Methods("POST")

	router.HandleFunc("/test/emailpassword/signin", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			TenantId    string                 `json:"tenantId"`
			Email       string                 `json:"email"`
			Password    string                 `json:"password"`
			UserContext map[string]interface{} `json:"userContext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if input.TenantId == "" {
			input.TenantId = "public"
		}

		var userContext supertokens.UserContext = nil
		if input.UserContext != nil {
			userContext = &input.UserContext
		}

		response, err := emailpassword.SignIn(input.TenantId, input.Email, input.Password, userContext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jsonResponse map[string]interface{}
		if response.OK != nil {
			jsonResponse = map[string]interface{}{
				"status": "OK",
				"user":   response.OK.User,
			}
		} else {
			jsonResponse = map[string]interface{}{
				"status": "WRONG_CREDENTIALS_ERROR",
			}
		}
		json.NewEncoder(w).Encode(jsonResponse)
	}).Methods("POST")

	router.HandleFunc("/test/emailpassword/createresetpasswordlink", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			TenantId    string                 `json:"tenantId"`
			UserId      string                 `json:"userId"`
			UserContext map[string]interface{} `json:"userContext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if input.TenantId == "" {
			input.TenantId = "public"
		}

		var userContext supertokens.UserContext = nil
		if input.UserContext != nil {
			userContext = &input.UserContext
		}

		response, err := emailpassword.CreateResetPasswordLink(input.TenantId, input.UserId, userContext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jsonResponse map[string]interface{}
		if response.OK != nil {
			jsonResponse = map[string]interface{}{
				"status": "OK",
				"link":   response.OK.Link,
			}
		} else if response.UnknownUserIdError != nil {
			jsonResponse = map[string]interface{}{
				"status": "UNKNOWN_USER_ID_ERROR",
			}
		} else {
			jsonResponse = map[string]interface{}{
				"status": "UNKNOWN_ERROR",
			}
		}
		json.NewEncoder(w).Encode(jsonResponse)
	}).Methods("POST")

	router.HandleFunc("/test/emailpassword/updateemailorpassword", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserId                    string                 `json:"userId"`
			Email                     *string                `json:"email"`
			Password                  *string                `json:"password"`
			ApplyPasswordPolicy       *bool                  `json:"applyPasswordPolicy"`
			TenantIdForPasswordPolicy *string                `json:"tenantIdForPasswordPolicy"`
			UserContext               map[string]interface{} `json:"userContext"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var userContext supertokens.UserContext = nil
		if input.UserContext != nil {
			userContext = &input.UserContext
		}

		response, err := emailpassword.UpdateEmailOrPassword(input.UserId, input.Email, input.Password, input.ApplyPasswordPolicy, input.TenantIdForPasswordPolicy, userContext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jsonResponse map[string]interface{}
		if response.OK != nil {
			jsonResponse = map[string]interface{}{
				"status": "OK",
			}
		} else if response.UnknownUserIdError != nil {
			jsonResponse = map[string]interface{}{
				"status": "UNKNOWN_USER_ID_ERROR",
			}
		} else if response.EmailAlreadyExistsError != nil {
			jsonResponse = map[string]interface{}{
				"status": "EMAIL_ALREADY_EXISTS_ERROR",
			}
		} else if response.PasswordPolicyViolatedError != nil {
			jsonResponse = map[string]interface{}{
				"status":        "PASSWORD_POLICY_VIOLATED_ERROR",
				"failureReason": response.PasswordPolicyViolatedError.FailureReason,
			}
		} else {
			jsonResponse = map[string]interface{}{
				"status": "UNKNOWN_ERROR",
			}
		}
		json.NewEncoder(w).Encode(jsonResponse)
	}).Methods("POST")
}
