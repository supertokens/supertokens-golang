package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	apiPort = 3030
)
var router = mux.NewRouter()

type customHandler struct{}

func (h customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	supertokens.Middleware(router).ServeHTTP(w, r)
}

func main() {

	// Initialize default SuperTokens configuration
	defaultSTInit()

	// Middleware
	router.Use(loggingMiddleware)

	// Routes
	router.HandleFunc("/test/ping", pingHandler).Methods("GET")
	router.HandleFunc("/test/init", initHandler).Methods("POST")
	router.HandleFunc("/test/overrideparams", overrideParamsHandler).Methods("GET")
	router.HandleFunc("/test/featureflag", featureFlagHandler).Methods("GET")
	router.HandleFunc("/test/resetoverrideparams", resetOverrideParamsHandler).Methods("POST")
	router.HandleFunc("/test/mockexternalapi", mockExternalAPIHandler).Methods("POST")
	router.HandleFunc("/test/getoverridelogs", getOverrideLogsHandler).Methods("GET")

	// Add routes for each recipe
	addEmailPasswordRoutes(router)
	addSessionRoutes(router)
	// addAccountLinkingRoutes(router)
	// addEmailVerificationRoutes(router)
	addMultitenancyRoutes(router)
	// addPasswordlessRoutes(router)
	// addMultiFactorAuthRoutes(router)
	addThirdPartyRoutes(router)
	// addTOTPRoutes(router)
	// addUserMetadataRoutes(router)

	// Custom routes for session tests
	router.HandleFunc("/create", createSessionHandler).Methods("POST")
	router.HandleFunc("/getsession", getSessionHandler).Methods("POST")
	router.HandleFunc("/refreshsession", refreshSessionHandler).Methods("POST")
	router.HandleFunc("/verify", verifySessionHandler).Methods("GET")

	// Error handling
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.Use(errorMiddleware)

	// Start server
	if envPort := os.Getenv("API_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			apiPort = p
		}
	}
	fmt.Printf("Starting server on port %d\n", apiPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", apiPort), customHandler{}))
}

func loggingOverrideFuncSync[T any](
	name string,
	originalImpl func(args ...any) (T, error),
) func(args ...any) (T, error) {
	return func(args ...any) (T, error) {
		logOverrideEvent(name, "CALL", args)
		res, err := originalImpl(args...)
		if err != nil {
			logOverrideEvent(name, "REJ", err)
		} else {
			logOverrideEvent(name, "RES", res)
		}
		return res, err
	}
}

func defaultSTInit() {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:3567",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			APIDomain:     "http://api.supertokens.io",
			WebsiteDomain: "http://localhost:3000",
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func STReset() {
	// resetOverrideParams()
	resetOverrideLogs()

	supertokens.ResetForTest()
}

func initST(config map[string]interface{}) {
	STReset()

	recipeList := recipeListFromRecipeConfigs(config["recipeList"].([]interface{}))

	supertokens.LogDebugMessage(fmt.Sprintf("initST: %v", config))

	var interceptor func(*http.Request, supertokens.UserContext) (*http.Request, error) = nil

	if st, ok := config["supertokens"].(map[string]interface{}); ok {
		if interceptorStr, ok := st["networkInterceptor"].(string); ok {
			interceptorFunc, err := GetFunc(interceptorStr)
			if err != nil {
				log.Fatal(err)
			}
			interceptor = interceptorFunc.(func(*http.Request, supertokens.UserContext) (*http.Request, error))
		}
	}

	finalInterceptor := func(request *http.Request, userContext supertokens.UserContext) (*http.Request, error) {
		return loggingOverrideFuncSync("networkInterceptor", func(args ...any) (*http.Request, error) {
			if interceptor != nil {
				return interceptor(args[0].(*http.Request), args[1].(supertokens.UserContext))
			}
			return args[0].(*http.Request), nil
		})(request, userContext)
	}

	parsedConfig := supertokens.TypeInput{
		AppInfo: supertokens.AppInfo{
			AppName:       config["appInfo"].(map[string]interface{})["appName"].(string),
			APIDomain:     config["appInfo"].(map[string]interface{})["apiDomain"].(string),
			WebsiteDomain: config["appInfo"].(map[string]interface{})["websiteDomain"].(string),
		},
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI:      config["supertokens"].(map[string]interface{})["connectionURI"].(string),
			NetworkInterceptor: finalInterceptor,
		},
		RecipeList: recipeList,
	}

	err := supertokens.Init(parsedConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func recipeListFromRecipeConfigs(recipeListMaps []interface{}) []supertokens.Recipe {
	recipeList := []supertokens.Recipe{}
	for _, recipeItem := range recipeListMaps {
		recipeItemMap := recipeItem.(map[string]interface{})
		if recipeItemMap["recipeId"] == "emailpassword" {
			var recipeConfigMap map[string]interface{}
			err := json.Unmarshal([]byte(recipeItemMap["config"].(string)), &recipeConfigMap)
			if err != nil {
				log.Printf("Error unmarshaling recipe config: %v", err)
				continue
			}
			recipeConfig := epmodels.TypeInput{}
			if signUpFeature, ok := recipeConfigMap["signUpFeature"].(map[string]interface{}); ok {
				recipeConfig.SignUpFeature = &epmodels.TypeInputSignUp{}
				if formFields, ok := signUpFeature["formFields"].([]interface{}); ok {
					for _, field := range formFields {
						if fieldMap, ok := field.(map[string]interface{}); ok {
							formField := epmodels.TypeInputFormField{
								ID: fieldMap["id"].(string),
							}
							recipeConfig.SignUpFeature.FormFields = append(recipeConfig.SignUpFeature.FormFields, formField)
						}
					}
				}
			}

			recipeList = append(recipeList, emailpassword.Init(&recipeConfig))
		} else if recipeItemMap["recipeId"] == "session" {
			var recipeConfigMap map[string]interface{}
			err := json.Unmarshal([]byte(recipeItemMap["config"].(string)), &recipeConfigMap)
			if err != nil {
				log.Printf("Error unmarshaling recipe config: %v", err)
				continue
			}
			recipeConfig := sessmodels.TypeInput{}

			// Populate recipeConfig from recipeConfigMap
			if cookieSecure, ok := recipeConfigMap["cookieSecure"].(bool); ok {
				recipeConfig.CookieSecure = &cookieSecure
			}
			if cookieSameSite, ok := recipeConfigMap["cookieSameSite"].(string); ok {
				recipeConfig.CookieSameSite = &cookieSameSite
			}
			if sessionExpiredStatusCode, ok := recipeConfigMap["sessionExpiredStatusCode"].(float64); ok {
				code := int(sessionExpiredStatusCode)
				recipeConfig.SessionExpiredStatusCode = &code
			}
			if invalidClaimStatusCode, ok := recipeConfigMap["invalidClaimStatusCode"].(float64); ok {
				code := int(invalidClaimStatusCode)
				recipeConfig.InvalidClaimStatusCode = &code
			}
			if cookieDomain, ok := recipeConfigMap["cookieDomain"].(string); ok {
				recipeConfig.CookieDomain = &cookieDomain
			}
			if olderCookieDomain, ok := recipeConfigMap["olderCookieDomain"].(string); ok {
				recipeConfig.OlderCookieDomain = &olderCookieDomain
			}
			if antiCsrf, ok := recipeConfigMap["antiCsrf"].(string); ok {
				recipeConfig.AntiCsrf = &antiCsrf
			}
			if exposeAccessTokenToFrontendInCookieBasedAuth, ok := recipeConfigMap["exposeAccessTokenToFrontendInCookieBasedAuth"].(bool); ok {
				recipeConfig.ExposeAccessTokenToFrontendInCookieBasedAuth = exposeAccessTokenToFrontendInCookieBasedAuth
			}
			if useDynamicAccessTokenSigningKey, ok := recipeConfigMap["useDynamicAccessTokenSigningKey"].(bool); ok {
				recipeConfig.UseDynamicAccessTokenSigningKey = &useDynamicAccessTokenSigningKey
			}

			recipeList = append(recipeList, session.Init(&recipeConfig))
		} else if recipeItemMap["recipeId"] == "emailverification" {
			var recipeConfigMap map[string]interface{}
			err := json.Unmarshal([]byte(recipeItemMap["config"].(string)), &recipeConfigMap)
			if err != nil {
				log.Printf("Error unmarshaling recipe config: %v", err)
				continue
			}
			recipeConfig := evmodels.TypeInput{
				Mode: evmodels.ModeOptional,
			}

			if mode, ok := recipeConfigMap["mode"].(string); ok {
				recipeConfig.Mode = evmodels.TypeMode(mode)
			}

			recipeConfig.Override = &evmodels.OverrideStruct{
				Functions: func(originalImplementation evmodels.RecipeInterface) evmodels.RecipeInterface {
					prefix := "EmailVerification.override.functions"

					oCreateEmailVerificationToken := *originalImplementation.CreateEmailVerificationToken
					*originalImplementation.CreateEmailVerificationToken = func(userID, email, tenantId string, userContext supertokens.UserContext) (evmodels.CreateEmailVerificationTokenResponse, error) {
						logOverrideEvent(prefix+".createEmailVerificationToken", "CALL", map[string]interface{}{
							"userId":   userID,
							"email":    email,
							"tenantId": tenantId,
						})
						res, err := oCreateEmailVerificationToken(userID, email, tenantId, userContext)
						if err != nil {
							logOverrideEvent(prefix, "REJ", err)
						} else {
							logOverrideEvent(prefix, "RES", res)
						}
						return res, err
					}

					oIsEmailVerified := *originalImplementation.IsEmailVerified
					*originalImplementation.IsEmailVerified = func(userID, email string, userContext supertokens.UserContext) (bool, error) {
						logOverrideEvent(prefix+".isEmailVerified", "CALL", map[string]interface{}{
							"userId": userID,
							"email":  email,
						})
						res, err := oIsEmailVerified(userID, email, userContext)
						if err != nil {
							logOverrideEvent(prefix, "REJ", err)
						} else {
							logOverrideEvent(prefix, "RES", res)
						}
						return res, err
					}

					oVerifyEmailUsingToken := *originalImplementation.VerifyEmailUsingToken
					*originalImplementation.VerifyEmailUsingToken = func(token string, tenantId string, userContext *map[string]interface{}) (evmodels.VerifyEmailUsingTokenResponse, error) {
						logOverrideEvent(prefix+".verifyEmailUsingToken", "CALL", map[string]interface{}{
							"token":    token,
							"tenantId": tenantId,
						})
						res, err := oVerifyEmailUsingToken(token, tenantId, userContext)
						if err != nil {
							logOverrideEvent(prefix, "REJ", err)
						} else {
							logOverrideEvent(prefix, "RES", res)
						}
						return res, err
					}

					return originalImplementation
				},
			}

			recipeList = append(recipeList, emailverification.Init(recipeConfig))
		} else if recipeItemMap["recipeId"] == "thirdparty" {
			var recipeConfigMap map[string]interface{}
			err := json.Unmarshal([]byte(recipeItemMap["config"].(string)), &recipeConfigMap)
			if err != nil {
				log.Printf("Error unmarshaling recipe config: %v", err)
				continue
			}
			recipeConfig := tpmodels.TypeInput{}

			if signInAndUpFeature, ok := recipeConfigMap["signInAndUpFeature"].(map[string]interface{}); ok {
				if providers, ok := signInAndUpFeature["providers"].([]interface{}); ok {
					for _, provider := range providers {
						providerInput := tpmodels.ProviderInput{}
						providerMap := provider.(map[string]interface{})

						if config, ok := providerMap["config"].(map[string]interface{}); ok {
							configBytes, err := json.Marshal(config)
							if err != nil {
								log.Printf("Error marshaling provider config: %v", err)
								continue
							}
							err = json.Unmarshal(configBytes, &providerInput.Config)
							if err != nil {
								log.Printf("Error unmarshaling provider config: %v", err)
								continue
							}
						}

						if includeInNonPublicTenantsByDefault, ok := providerMap["includeInNonPublicTenantsByDefault"].(bool); ok {
							providerInput.IncludeInNonPublicTenantsByDefault = includeInNonPublicTenantsByDefault
						}
						recipeConfig.SignInAndUpFeature.Providers = append(recipeConfig.SignInAndUpFeature.Providers, providerInput)
					}
				}
			}

			recipeList = append(recipeList, thirdparty.Init(&recipeConfig))
		}
	}
	return recipeList
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Config string `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var config map[string]interface{}
	err := json.Unmarshal([]byte(input.Config), &config)
	if err != nil {
		http.Error(w, "Failed to parse config: "+err.Error(), http.StatusBadRequest)
		return
	}

	initST(config)

	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func overrideParamsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("TODO")
}

func featureFlagHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{})
}

func resetOverrideParamsHandler(w http.ResponseWriter, r *http.Request) {
	resetOverrideLogs()
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func mockExternalAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Implement mock external API logic
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func getOverrideLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs": getOverrideLogs(),
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("go-test-server: route not found %s %s", r.Method, r.URL.Path), http.StatusNotFound)
}

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func createSessionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RecipeUserId string `json:"recipeUserId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recipeUserId := input.RecipeUserId

	_, err := session.CreateNewSession(r, w, "public", recipeUserId, map[string]interface{}{}, map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionContainer, err := session.GetSession(r, w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"userId":       sessionContainer.GetUserID(),
		"recipeUserId": sessionContainer.GetUserID(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func refreshSessionHandler(w http.ResponseWriter, r *http.Request) {
	sessionContainer, err := session.RefreshSession(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"userId":       sessionContainer.GetUserID(),
		"recipeUserId": sessionContainer.GetUserID(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func verifySessionHandler(w http.ResponseWriter, r *http.Request) {
	session.VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK"})
	})(w, r)
}

// Implement the remaining recipe-specific route handlers
