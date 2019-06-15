package authx

import (
	"time"

	"gopkg.in/oauth2.v3"
)

//export
//nolint[lll]
func SetupPasswordAuthServer(client ClientInfo, pwdHandler func(username, password string) (string, error), tokenStore oauth2.TokenStore, opts ...AuthOption) *Server {
	server := NewOAuthServer()
	if tokenStore != nil {
		server.SetTokenStorage(tokenStore)
	}
	server.SetClients(client)
	server.SetPasswordAuthorizationHandler(pwdHandler)
	server.SetExtensionFieldsHandler(NewTokenExtension(UserExtension))
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(server)
		}
	}
	return server
}

type AuthOption func(server *Server)

//export
func FileStore(file string) AuthOption {
	return func(server *Server) {
		store, err := NewFileTokenStore(file)
		if err == nil {
			server.SetTokenStorage(store)
		}
	}
}

//export
func MemeoryStore(server *Server) {
	server.SetTokenStorage(NewMemoryTokenStore())
}

//export
func WebTokenCfg(exps ...time.Duration) AuthOption {
	return func(server *Server) {
		accessTokenExp := time.Hour * 2
		refreshTokenExp := time.Hour * 3
		if len(exps) > 0 {
			accessTokenExp = exps[0]
		}
		if len(exps) > 1 {
			refreshTokenExp = exps[1]
		}

		server.SetPasswordTokenCfg(accessTokenExp, refreshTokenExp, true)
		server.SetRefreshTokenCfg(accessTokenExp, refreshTokenExp, true, true, true, true)
	}
}

//export
func AppTokenCfg(exps ...time.Duration) AuthOption {
	return func(server *Server) {
		accessTokenExp := time.Hour * 24
		refreshTokenExp := time.Hour * 24 * 7
		if len(exps) > 0 {
			accessTokenExp = exps[0]
		}
		if len(exps) > 1 {
			refreshTokenExp = exps[1]
		}

		server.SetPasswordTokenCfg(accessTokenExp, refreshTokenExp, true)
		server.SetRefreshTokenCfg(accessTokenExp, refreshTokenExp, true, true, true, true)
	}
}
