package randx

import (
	uuid "github.com/satori/go.uuid"
)

//export
func GenUUID(name string) string {
	uid := uuid.NewV4()
	key := uuid.Must(uuid.NewV5(uid, name), nil)
	return key.String()
}

//export
func GetUUID(name string) string {
	uid := uuid.NewV4()
	key := uuid.Must(uuid.NewV5(uid, name), nil)
	return key.String()
}
