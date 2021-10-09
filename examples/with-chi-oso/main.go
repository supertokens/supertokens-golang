package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/osohq/go-oso"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

var (
	osoClient oso.Oso
)

func main() {
	initAuth()

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
			emailpassword.Init(nil),
			session.Init(nil),
		},
	})

	if err != nil {
		panic(err.Error())
	}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   append([]string{"Content-Type"}, supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	r.Use(supertokens.Middleware)

	r.Get("/sessioninfo", session.VerifySession(nil, sessioninfo))
	r.Get("/repo/{repoName}", session.VerifySession(nil, getRepo))

	http.ListenAndServe(":3001", r)
}

func sessioninfo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())

	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	sessionData, err := sessionContainer.GetSessionData()
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
		"sessionHandle": sessionContainer.GetHandle(),
		"userId":        sessionContainer.GetUserID(),
		"jwtPayload":    sessionContainer.GetJWTPayload(),
		"sessionData":   sessionData,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error in converting to json"))
	} else {
		w.Write(bytes)
	}
}

func getRepo(w http.ResponseWriter, r *http.Request) {
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	if sessionContainer == nil {
		w.WriteHeader(500)
		w.Write([]byte("no session found"))
		return
	}
	repoName := chi.URLParam(r, "repoName")
	if repoName == "" {
		w.WriteHeader(400)
		w.Write([]byte("repository name required"))
		return
	}
	repository := GetRepositoryByName(repoName)
	currentUser := GetCurrentUser(sessionContainer.GetUserID())
	err := osoClient.Authorize(currentUser, "read", repository)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Welcome to repo: %s", repository.Name)))
	} else {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized"))
	}
}

func initAuth() {
	var err error
	osoClient, err = oso.NewOso()
	if err != nil {
		fmt.Sprintf("Failed to set up Oso: %v", err)
		os.Exit(1)
	}
	osoClient.RegisterClass(reflect.TypeOf(Repository{}), nil)
	osoClient.RegisterClass(reflect.TypeOf(User{}), nil)
	if err = osoClient.LoadFiles([]string{"main.polar"}); err != nil {
		fmt.Sprintf("Failed to start: %s", err)
		os.Exit(1)
	}
}
