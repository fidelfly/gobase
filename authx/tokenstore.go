package authx

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/store"
)

func NewMemoryTokenStore() oauth2.TokenStore {
	tokenStore, _ := store.NewMemoryTokenStore()
	return tokenStore
}

func NewFileTokenStore(filename string) (oauth2.TokenStore, error) {
	return store.NewFileTokenStore(filename)
}
