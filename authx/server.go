package authx

import (
	"net/http"
	"time"

	"gopkg.in/oauth2.v3/errors"

	"github.com/fidelfly/gox/errorx"

	"github.com/fidelfly/gox/logx"

	"gopkg.in/oauth2.v3/models"

	"gopkg.in/oauth2.v3/store"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
)

type Server struct {
	manager *manage.Manager
	server  *server.Server
}

type ClientInfo interface {
	oauth2.ClientInfo
}

type AuthClient models.Client

func (ac *AuthClient) GetID() string {
	return ac.ID
}
func (ac *AuthClient) GetSecret() string {
	return ac.Secret
}
func (ac *AuthClient) GetDomain() string {
	return ac.Domain
}
func (ac *AuthClient) GetUserID() string {
	return ac.UserID
}

func NewOAuthServer() *Server {
	m := manage.NewDefaultManager()
	s := server.NewDefaultServer(m)
	return &Server{m, s}
}

func (as *Server) SetTokenStorage(tokenStore oauth2.TokenStore) {
	as.manager.MapTokenStorage(tokenStore)
}

func (as *Server) SetPasswordTokenCfg(accessTokenExp, refreshTokenExp time.Duration, isGenerateRefresh bool) {
	as.manager.SetPasswordTokenCfg(
		&manage.Config{
			AccessTokenExp:    accessTokenExp,
			RefreshTokenExp:   refreshTokenExp,
			IsGenerateRefresh: isGenerateRefresh,
		})
}

// nolint[lll]
func (as *Server) SetRefreshTokenCfg(accessTokenExp, refreshTokenExp time.Duration, isGenerateRefresh, isRemoveAccess, isRemoveRefreshing bool, isResetRefreshTime bool) {
	as.manager.SetRefreshTokenCfg(
		&manage.RefreshingConfig{
			AccessTokenExp:     accessTokenExp,
			RefreshTokenExp:    refreshTokenExp,
			IsGenerateRefresh:  isGenerateRefresh,
			IsRemoveAccess:     isRemoveAccess,
			IsRemoveRefreshing: isRemoveRefreshing,
			IsResetRefreshTime: isResetRefreshTime,
		})
}

func (as *Server) SetClients(clients ...ClientInfo) {
	clientStore := store.NewClientStore()
	for _, client := range clients {
		logx.CaptureError(clientStore.Set(client.GetID(), client))
	}
	as.manager.MapClientStorage(clientStore)
}

func (as *Server) SetClientStore(clientStore oauth2.ClientStore) {
	as.manager.MapClientStorage(clientStore)
}

func (as *Server) SetPasswordAuthorizationHandler(handler func(username, password string) (string, error)) {
	as.server.PasswordAuthorizationHandler = handler
}

func (as *Server) SetExtensionFieldsHandler(handler func(ti oauth2.TokenInfo) map[string]interface{}) {
	as.server.ExtensionFieldsHandler = handler
}

func (as *Server) HandleTokenRequest(w http.ResponseWriter, r *http.Request) {
	logx.CaptureError(as.server.HandleTokenRequest(w, r))
}

func (as *Server) ValidateToken(w http.ResponseWriter, r *http.Request) (ti oauth2.TokenInfo, err error) {
	ti, err = as.server.ValidationBearerToken(r)
	if err != nil {
		switch err {
		case errors.ErrInvalidAccessToken:
		case errors.ErrExpiredRefreshToken:
			err = errorx.NewCodeError(err, UnauthorizedErrorCode)
		case errors.ErrExpiredAccessToken:
			err = errorx.NewCodeError(err, TokenExpiredErrorCode)
		default:
			err = errorx.NewCodeError(err, UnauthorizedErrorCode)
		}
		return
	}

	/*if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).After(time.Now()) {
		return
	}

	if ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
		err = errorx.NewError(UnauthorizedErrorCode, "refresh token is expired")
		return
	}
	err = errorx.NewError(TokenExpiredErrorCode, "token is expired")*/
	return
}

func (as *Server) RemoveAccessToken(access string) (err error) {
	return as.manager.RemoveAccessToken(access)
}

func (as *Server) RemoveRefreshToken(refresh string) (err error) {
	return as.manager.RemoveRefreshToken(refresh)
}
