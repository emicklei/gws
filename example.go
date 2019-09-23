package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func cmdShowExamples(c *cli.Context) error {
	examples := `
- - - - - - - - - - - - - - - - - - - - 
List the email addresses of all users.
	
	gdom user list

List the groups of which the user is a member.

	gdom user membership john.doe@company.com

Show details of a user.

	gdom user info john.doe@company.com

List the email address of all groups

	gdom group list

List the members of a group

	gdom group members all@company.com

Show details of a group.

	gdom group info all@company.com

See full documentation on https://github.com/emicklei/gdom
- - - - - - - - - - - - - - - - - - - -
`
	fmt.Println(examples)
	return nil
}
