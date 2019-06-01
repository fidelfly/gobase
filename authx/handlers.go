package authx

import (
	"strconv"

	"gopkg.in/oauth2.v3"
)

type TokenExtension func(map[string]interface{}, oauth2.TokenInfo)

func UserExtension(fieldsValue map[string]interface{}, ti oauth2.TokenInfo) {
	fieldsValue["user_id"], _ = strconv.Atoi(ti.GetUserID())
}

func NewTokenExtension(tes ...TokenExtension) func(oauth2.TokenInfo) map[string]interface{} {
	return func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{}) {
		data := make(map[string]interface{}, 1)
		for _, te := range tes {
			te(data, ti)
		}
		return data
	}
}
