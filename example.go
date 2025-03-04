package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func cmdShowExamples(_ *cli.Context) error {
	examples := `
- - - - - - - - - - - - - - - - - - - - 
List the email addresses of all users.
	
	gws user list

List the groups of which the user is a member.

	gws user membership john.doe
	gws user membership john.doe@company.com

Show details of a user.

	gws user info john.doe
	gws user info john.doe@company.com

Manage users

	gws user suspend martin "left the company"

List the email address of all groups

	gws group list

List the members of a group

	gws group members all
	gws group members all@company.com

Show details of a group.

	gws group info all
	gws group info all@company.com

Managing groups

	gws group create brand-new-group
	gws group delete my-old-group
	gws group add my-group john.doe
	gws group remove my-group john.doe

List the available roles to manage.

	gws role list

List the users who have the administration role

	gws role assignments _USER_MANAGEMENT_ADMIN_ROLE

List the (internet) domains that are managed

	gws domain list

See full documentation on https://github.com/emicklei/gws
- - - - - - - - - - - - - - - - - - - -
`
	fmt.Println(examples)
	return nil
}
