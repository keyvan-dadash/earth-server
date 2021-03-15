package user

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func ParseUserFromMap(retrievedUserFieldsMap map[string]interface{}, user *User) error { //TODO: exactly how we need error?

	user.username = retrievedUserFieldsMap["username"].(string)

	user.password = retrievedUserFieldsMap["password"].(string)

	user.Email = retrievedUserFieldsMap["email"].(string)

	user.Nickname = retrievedUserFieldsMap["nickname"].(string)

	user.UUID = retrievedUserFieldsMap["uuid"].(uuid.UUID)

	user.IsOnline = retrievedUserFieldsMap["uuid"].(bool)

	user.JoinedDate = retrievedUserFieldsMap["uuid"].(time.Time)

	return nil

}
