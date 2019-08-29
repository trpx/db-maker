package utils

import (
	"errors"
	"flag"
	"fmt"
)

type Options struct {
	Engine        *string
	Host          *string
	Port          *string
	AdminUser     *string
	AdminDB       *string
	AdminPassFile *string
	User          *string
	UserDB        *string
	UserPassFile  *string
}

func ParseOptions() (opt Options, parseErr error) {
	defer func() {
		if r := recover(); r != nil {
			parseErr = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	opt.Engine = flag.String("engine", "postgresql", "db engine")

	opt.Host = flag.String("host", "localhost", "db host")
	opt.Port = flag.String("port", "5432", "db port")

	opt.AdminUser = flag.String("admin-user", "postgres", "db port")
	opt.AdminDB = flag.String("admin-db", "postgres", "db port")

	opt.AdminPassFile = flag.String("admin-pass-file", "", "file with admin passwords, one per line")
	opt.User = flag.String("user", "", "db user to create (if not exists)")
	opt.UserDB = flag.String("user-db", "", "user db to create (if not exists)")
	opt.UserPassFile = flag.String("user-pass-file", "", "file with user passwords, one per line")

	flag.Parse()

	// validate engine
	if *opt.Engine != "postgresql" {
		Panicf("engine '%s' is not supported, the only engine supported for now is 'postgresql'", *opt.Engine)
	}

	adminPasswords, userPasswords := ReadPasswordsFromEnv()

	requiredOptions := [][]string{
		{"user", *opt.User},
		{"user-db", *opt.UserDB},
	}

	msgPostfix := ""
	if len(adminPasswords) == 0 && len(*opt.AdminPassFile) == 0 {
		requiredOptions = append(requiredOptions, []string{"admin-pass-file", *opt.AdminPassFile})
		msgPostfix += "\nnote: --admin-pass-file may be substituted by setting DB_MAKER_ADMIN_PASSWORDS env variable"
	}
	if len(userPasswords) == 0 && len(*opt.UserPassFile) == 0 {
		requiredOptions = append(requiredOptions, []string{"user-pass-file", *opt.UserPassFile})
		msgPostfix += "\nnote: --user-pass-file may be substituted by setting DB_MAKER_USER_PASSWORDS env variable"
	}

	// validate required args
	var needOptions []string
	for _, i := range requiredOptions {
		if len(i[1]) == 0 {
			needOptions = append(needOptions, "--"+i[0])
		}
	}
	if len(needOptions) > 0 {
		Panicf("options %s are required"+msgPostfix, needOptions)
	}

	// validate unexpected args
	if len(flag.Args()) > 0 {
		Panicf("unexpected args: %s", flag.Args())
	}

	return opt, parseErr
}
