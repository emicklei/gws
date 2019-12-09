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

	gsuite user membership john.doe
	gsuite user membership john.doe@company.com

Show details of a user.

	gsuite user info john.doe
	gsuite user info john.doe@company.com

List the email address of all groups

	gsuite group list

List the members of a group

	gsuite group members all
	gsuite group members all@company.com

Show details of a group.

	gsuite group info all
	gsuite group info all@company.com

Managing groups

	gsuite group delete my-old-group
	gsuite group add my-group important-person
	gsuite group remove my-group important-person 

List the available roles to manage.

	gsuite role list

List the users who have the administration role

	gsuite role assignments _USER_MANAGEMENT_ADMIN_ROLE

List the (internet) domains that are managed

	gsuite domain list

See full documentation on https://github.com/emicklei/gsuite
- - - - - - - - - - - - - - - - - - - -
`
	fmt.Println(examples)
	return nil
}
