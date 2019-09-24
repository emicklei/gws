package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func cmdShowExamples(c *cli.Context) error {
	examples := `
- - - - - - - - - - - - - - - - - - - - 
List the email addresses of all users.
	
	gsuite user list

List the groups of which the user is a member.

	gsuite user membership john.doe@company.com

Show details of a user.

	gsuite user info john.doe@company.com

List the email address of all groups

	gsuite group list

List the members of a group

	gsuite group members all@company.com

Show details of a group.

	gsuite group info all@company.com

List the available roles to manage.

	gsuite role list

List the users who have the administration role

	gsuite role assignments _USER_MANAGEMENT_ADMIN_ROLE

See full documentation on https://github.com/emicklei/gsuite
- - - - - - - - - - - - - - - - - - - -
`
	fmt.Println(examples)
	return nil
}
