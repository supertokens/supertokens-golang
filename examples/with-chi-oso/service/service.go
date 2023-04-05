package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/osohq/go-oso"
	"github.com/supertokens/supertokens-golang/examples/with-chi-oso/models"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type UserReader interface {
	GetCurrentUser(userID string) (models.User, error)
}

type RepoReader interface {
	GetRepositoryByName(name string) (models.Repository, error)
}

type DbReader interface {
	UserReader
	RepoReader
}

type Opts struct {
	Database DbReader
	Auth     oso.Oso
}

type service struct {
	db   DbReader
	auth oso.Oso
}

func NewService(opts Opts) *service {
	return &service{db: opts.Database, auth: opts.Auth}
}

func (s *service) Sessioninfo(w http.ResponseWriter, r *http.Request) {
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

func (s *service) Repo(w http.ResponseWriter, r *http.Request) {
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
	repository, err := s.db.GetRepositoryByName(repoName)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	userByID, err := thirdpartyemailpassword.GetUserById(sessionContainer.GetUserID())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	currentUser, err := s.db.GetCurrentUser(userByID.Email)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.auth.Authorize(currentUser, "read", repository)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Welcome to repo: %s", repository.Name)))
	} else {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized"))
	}
}
