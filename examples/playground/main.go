package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const (
	ReqIDKey = "reqID"
)

func main() {
	newLogger := loggerImpl{
		log.New(os.Stdout, "", 0),
	}

	apiBasePath := "/"
	websiteBasePath := "/"

	supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// These are the connection details of the app you created on supertokens.com
			ConnectionURI: "https://try.supertokens.io",
			APIKey:        "9-o4QBtuWy65L1iZskNiug1DRFKWh6",
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "Auth Service",
			APIDomain:       "http://localhost:3001",
			WebsiteDomain:   "http://localhost:9000",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
		},
		Telemetry:             nil,
		OnSuperTokensAPIError: nil,
		Log:                   &newLogger,
		RequestIDKey:          middleware.RequestIDKey,
	},
	)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))
	r.Use(middleware.RequestID)
	r.Use(supertokens.Middleware)

	r.Get("/sessioninfo", session.VerifySession(nil, sessioninfo))

	http.ListenAndServe(":3001", r)
}

func sessioninfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		err = supertokens.ErrorHandler(err, r, w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}
		return
	}
	w.WriteHeader(200)
	w.Header().Add("content-type", "application/json")
	bytes, err := json.Marshal(map[string]interface{}{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error in converting to json"))
	} else {
		w.Write(bytes)
	}
}

type loggerImpl struct {
	*log.Logger
}

func (l *loggerImpl) Log(msg string) {
	l.Println(msg)
}
