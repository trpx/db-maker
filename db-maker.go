package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/trpx/db-maker/utils"
)


func main() {

	opt, err := utils.ParseOptions()
	if err != nil {
		utils.Fatalf("%v", err)
	}

	fmt.Printf(
		"engine: %#v\n" +
		"host: %#v\n" +
		"port: %#v\n" +
		"adminUser: %#v\n" +
		"adminDB: %#v\n" +
		"adminPassFile: %#v\n" +
		"user: %#v\n" +
		"userDB: %#v\n" +
		"userPassFile: %#v\n",
		*opt.Engine,
		*opt.Host,
		*opt.Port,
		*opt.AdminUser,
		*opt.AdminDB,
		*opt.AdminPassFile,
		*opt.User,
		*opt.UserDB,
		*opt.UserPassFile,
	)
	fmt.Printf("tail: %#v\n", flag.Args())
}
