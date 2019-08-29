package main

import (
	"errors"
	"fmt"
	"github.com/trpx/db-maker/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fatalError := errors.New(fmt.Sprintf("%v", r))
			utils.Fatalf("%v", fatalError)
		}
	}()

	opt, err := utils.ParseOptions()
	if err != nil {
		utils.Panicf("%v", err)
	}

	adminPasswords, userPasswords := utils.ReadPasswordsFromEnv()
	if len(adminPasswords) == 0 {
		adminPasswords = utils.ReadPassFile(opt.AdminPassFile)
	}
	if len(userPasswords) == 0 {
		userPasswords = utils.ReadPassFile(opt.UserPassFile)
	}

	db := utils.DB{
		User:      *opt.AdminUser,
		Passwords: adminPasswords,
		Name:      *opt.AdminDB,
		Host:      *opt.Host,
		Port:      *opt.Port,
	}
	db.Connect()
	defer func() {
		db.Disconnect()
	}()
	updatedAdminPass := db.UpdateUserPassword(*opt.AdminUser, adminPasswords[0])
	createdUser := db.CreateUserIfNotExists(*opt.User, userPasswords[0])
	createdDb := db.CreateDBIfNotExistsWithOwner(*opt.UserDB, *opt.User)
	updatedUserPass := db.UpdateUserPassword(*opt.User, userPasswords[0])

	msg := `updated admin pass:		%v
created user:		%v
created db:		%v
updated user pass:		%v`
	fmt.Printf(msg, updatedAdminPass, createdUser, createdDb, updatedUserPass)
}
