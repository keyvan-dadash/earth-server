package user

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func ParseUserFromMap(retrievedUserFieldsMap map[string]interface{}, user *User) error { //TODO: exactly how we need error?

	user.Username = retrievedUserFieldsMap["username"].(string)

	user.Password = retrievedUserFieldsMap["password"].(string)

	user.Email = retrievedUserFieldsMap["email"].(string)

	user.Nickname = retrievedUserFieldsMap["nickname"].(string)

	user.UUID = retrievedUserFieldsMap["uuid"].(uuid.UUID)

	// user.Isonline = retrievedUserFieldsMap["uuid"].(bool)

	user.JoinedDate = retrievedUserFieldsMap["uuid"].(time.Time)

	return nil

}
