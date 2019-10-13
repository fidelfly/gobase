package gosrvx

import (
	"gopkg.in/oauth2.v3"

	"github.com/fidelfly/gox/authx"
)

//export
// nolint[lll]
func SetupPasswordAuthorizeServer(client authx.ClientInfo, pwdHandler func(username, password string) (string, error), storeFile string) *authx.Server {
	server := authx.NewOAuthServer()
	var tokenStore oauth2.TokenStore
	if len(storeFile) > 0 {
		fileStore, err := authx.NewFileTokenStore(storeFile)
		if err != nil {

		} else {
			tokenStore = fileStore
		}
	}
	if tokenStore == nil {
		tokenStore = authx.NewMemoryTokenStore()
	}
	server.SetTokenStorage(tokenStore)
	server.SetClients(client)
	server.SetPasswordAuthorizationHandler(pwdHandler)
	server.SetExtensionFieldsHandler(authx.NewTokenExtension(authx.UserExtension))
	return server
}
