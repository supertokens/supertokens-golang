package main

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func main() {
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://try.supertokens.io",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens Demo App",
			APIDomain:     "http://localhost:3001",
			WebsiteDomain: "http://localhost:3000",
		},
		RecipeList: []supertokens.Recipe{
			emailverification.Init(evmodels.TypeInput{
				Mode: evmodels.ModeRequired,
			}),
			thirdpartyemailpassword.Init(&tpepmodels.TypeInput{
				/*
				   We use different credentials for different platforms when required. For example the redirect URI for Github
				   is different for Web and mobile. In such a case we can provide multiple providers with different client Ids.
				   When the frontend makes a request and wants to use a specific clientId, it needs to send the clientId to use in the
				   request. In the absence of a clientId in the request the SDK uses the default provider, indicated by `isDefault: true`.
				   When adding multiple providers for the same type (Google, Github etc), make sure to set `isDefault: true`.
				*/
				Providers: []tpmodels.ProviderInput{
					// We have provided you with development keys which you can use for testsing.
					// IMPORTANT: Please replace them with your own OAuth keys for production use.
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientType:   "web",
									ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
									ClientSecret: "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",
								},
								{
									// we use this for mobile apps
									ClientType:   "mobile",
									ClientID:     "1060725074195-c7mgk8p0h27c4428prfuo3lg7ould5o7.apps.googleusercontent.com",
									ClientSecret: "", // this is empty because we follow Authorization code grant flow via PKCE for mobile apps (Google doesn't issue a client secret for mobile apps).
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "github",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientType:   "web",
									ClientID:     "467101b197249757c71f",
									ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
								},
								{
									// We use this for mobile apps
									ClientType:   "mobile",
									ClientID:     "8a9152860ce869b64c44",
									ClientSecret: "00e841f10f288363cd3786b1b1f538f05cfdbda2",
								},
							},
						},
					},
					/*
					   For Apple signin, iOS apps always use the bundle identifier as the client ID when communicating with Apple. Android, Web and other platforms
					   need to configure a Service ID on the Apple developer dashboard and use that as client ID.
					   In the example below 4398792-io.supertokens.example.service is the client ID for Web. Android etc and thus we mark it as default. For iOS
					   the frontend for the demo app sends the clientId in the request which is then used by the SDK.
					*/
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "apple",
							Clients: []tpmodels.ProviderClientConfig{
								{
									// For Android and website apps
									ClientType: "web",
									ClientID:   "4398792-io.supertokens.example.service",
									AdditionalConfig: map[string]interface{}{
										"keyId":      "7M48Y4RYDL",
										"privateKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
										"teamId":     "YWQCXGJRJL",
									},
								},
								{
									// For iOS Apps
									ClientType: "ios",
									ClientID:   "4398792-io.supertokens.example",
									AdditionalConfig: map[string]interface{}{
										"keyId":      "7M48Y4RYDL",
										"privateKey": "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
										"teamId":     "YWQCXGJRJL",
									},
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "discord",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientType:   "web",
									ClientID:     "4398792-907871294886928395",
									ClientSecret: "His4yXGEovVp5TZkZhEAt0ZXGh8uOVDm",
								},
								{
									// We use this for mobile apps
									ClientType:   "mobile",
									ClientID:     "4398792-907871294886928395",
									ClientSecret: "His4yXGEovVp5TZkZhEAt0ZXGh8uOVDm",
								},
							},
						},
					},
					{
						Config: tpmodels.ProviderConfig{
							ThirdPartyId: "google-workspaces",
							Clients: []tpmodels.ProviderClientConfig{
								{
									ClientType:   "web",
									ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
									ClientSecret: "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",
								},
								{
									// We use this for mobile apps
									ClientType:   "mobile",
									ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
									ClientSecret: "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",
								},
							},
						},
					},
				},
			}),
			session.Init(nil),
			dashboard.Init(nil),
		},
	})
	if err != nil {
		log.Fatal("Something went wrong while starting up supertokens: ", err.Error())
	}
	app := fiber.New()

	allowedHeaders := append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...)
	allowedHeadersInCommaSeparetedStringFormat := stringArrayToStringConvertor(allowedHeaders)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST, PUT, HEAD, OPTIONS",
		AllowHeaders:     allowedHeadersInCommaSeparetedStringFormat,
		AllowCredentials: true,
	}))

	//adding the supertokens middleware
	app.Use(adaptor.HTTPMiddleware(supertokens.Middleware))

	app.Get("/sessionInfo", verifySession(nil), sessioninfo)
	log.Fatal(app.Listen(":3001"))
}

// wrapper of the original implementation of verify session to match the required function signature
func verifySession(options *sessmodels.VerifySessionOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var errFromNextHandler error
		err := adaptor.HTTPHandlerFunc(session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.SetUserContext(r.Context())
			errFromNextHandler = c.Next()

			if errFromNextHandler != nil {
				// just in case a supertokens error was returned, we call the supertokens error handler
				// also, if supertokens error was handled, we don't want to return it, hence updating errFromNextHandler
				errFromNextHandler = supertokens.ErrorHandler(errFromNextHandler, r, rw)
			}
		}))(c)

		if err != nil {
			return err
		}
		return errFromNextHandler
	}
}

func sessioninfo(c *fiber.Ctx) error {
	sessionContainer := session.GetSessionFromRequestContext(c.UserContext())
	if sessionContainer == nil {
		return c.Status(500).JSON("no session found")
	}
	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	c.Response().Header.Add("content-type", "application/json")

	currAccessTokenPayload := sessionContainer.GetAccessTokenPayload()
	counter, ok := currAccessTokenPayload["counter"]
	if !ok {
		counter = 1
	} else {
		counter = int(counter.(float64) + 1)
	}
	err = sessionContainer.MergeIntoAccessTokenPayload(map[string]interface{}{
		"counter": counter.(int),
	})
	if err != nil {
		return err
	}
	return c.Status(200).JSON(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
}

// utility funtion to help convert an array of string to convert to comma separeted string format
func stringArrayToStringConvertor(stringArray []string) string {
	var stringToBeReturned string
	for _, val := range stringArray {
		stringToBeReturned += (val + ", ")
	}
	return stringToBeReturned
}
