package authx

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/store"

	"github.com/fidelfly/gox/logx"
)

func NewMemoryTokenStore() oauth2.TokenStore {
	tokenStore, _ := store.NewMemoryTokenStore()
	return tokenStore
}

func NewFileTokenStore(filename string) (oauth2.TokenStore, error) {
	return store.NewFileTokenStore(filename)
}

type multiLevelTokenStore struct {
	stores []oauth2.TokenStore
}

func NewMultiLevelTokenStore(stores ...oauth2.TokenStore) oauth2.TokenStore {
	if len(stores) == 0 {
		return nil
	}
	return &multiLevelTokenStore{stores}
}

// create and store the new token information
func (s *multiLevelTokenStore) Create(info oauth2.TokenInfo) (err error) {
	for i := len(s.stores) - 1; i >= 0; i-- {
		err = s.stores[i].Create(info)
		if err != nil {
			logx.Error(err)
			if i == len(s.stores)-1 {
				return
			}
		}
	}
	return
}

// delete the authorization code
func (s *multiLevelTokenStore) RemoveByCode(code string) (err error) {
	for i := 0; i < len(s.stores); i++ {
		err = s.stores[i].RemoveByCode(code)
		if err != nil {
			logx.Error(err)
			return
		}
	}
	return
}

// use the access token to delete the token information
func (s *multiLevelTokenStore) RemoveByAccess(access string) (err error) {
	for i := 0; i < len(s.stores); i++ {
		err = s.stores[i].RemoveByAccess(access)
		if err != nil {
			logx.Error(err)
			return
		}
	}
	return
}

// use the refresh token to delete the token information
func (s *multiLevelTokenStore) RemoveByRefresh(refresh string) (err error) {
	for i := 0; i < len(s.stores); i++ {
		err = s.stores[i].RemoveByRefresh(refresh)
		if err != nil {
			logx.Error(err)
			return
		}
	}
	return
}

// use the authorization code for token information data
func (s *multiLevelTokenStore) GetByCode(code string) (info oauth2.TokenInfo, err error) {
	for i := 0; i < len(s.stores); i++ {
		info, err = s.stores[i].GetByCode(code)
		if err == nil {
			if info != nil && i > 0 {
				for i--; i >= 0; i-- {
					logx.CaptureError(s.stores[i].Create(info))
				}
			}
			return
		}
	}
	return
}

// use the access token for token information data
func (s *multiLevelTokenStore) GetByAccess(access string) (info oauth2.TokenInfo, err error) {
	for i := 0; i < len(s.stores); i++ {
		info, err = s.stores[i].GetByAccess(access)
		if err == nil {
			if info != nil && i > 0 {
				for i--; i >= 0; i-- {
					logx.CaptureError(s.stores[i].Create(info))
				}
			}
			return
		}
	}
	return
}

// use the refresh token for token information data
func (s *multiLevelTokenStore) GetByRefresh(refresh string) (info oauth2.TokenInfo, err error) {
	for i := 0; i < len(s.stores); i++ {
		info, err = s.stores[i].GetByRefresh(refresh)
		if err == nil {
			if info != nil && i > 0 {
				for i--; i >= 0; i-- {
					logx.CaptureError(s.stores[i].Create(info))
				}
			}
			return
		}
	}
	return
}
