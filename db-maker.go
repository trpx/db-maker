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

	// create extensions if any
	var createdExtensions []string
	if len(opt.Extensions) > 0 {
		db := utils.DB{
			User:      *opt.AdminUser,
			Passwords: adminPasswords,
			Name:      *opt.UserDB,
			Host:      *opt.Host,
			Port:      *opt.Port,
		}
		db.Connect()
		defer func() {
			db.Disconnect()
		}()
		for _, extension := range opt.Extensions {
			created := db.CreateExtensionIfNotExists(extension)
			if created {
				createdExtensions = append(createdExtensions, extension)
			}
		}
	}

	msg := `updated admin '%s' pass:		%v
created user '%s':			%v
created db '%s':			%v
updated user pass:			%v
created extensions:			%v`
	fmt.Printf(
		msg,
		*opt.AdminUser, updatedAdminPass,
		*opt.User, createdUser,
		*opt.UserDB, createdDb,
		updatedUserPass, createdExtensions,
	)
}
