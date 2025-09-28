package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

func cmdRoleList(c *cli.Context) error {
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	// TODO my_customer?
	roles, err := srv.Roles.List(myAccoutsCustomerID).MaxResults(int64(ifZero(c.Int("limit"), 100))).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve roles in domain: %v", err)
	}

	if optionJSON(c, roles.Items) {
		return nil
	}
	for _, u := range roles.Items {
		fmt.Println(u.RoleName)
	}
	return nil
}

func cmdRoleAssignment(c *cli.Context) error {
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	// Get all roles
	// TODO my_customer?
	roles, err := srv.Roles.List(myAccoutsCustomerID).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve roles in domain: %v", err)
	}

	// find role by name
	roleName := c.Args().Get(0)
	var roleID int64
	for _, each := range roles.Items {
		if each.RoleName == roleName {
			roleID = each.RoleId
			break
		}
	}

	if roleID == 0 {
		return fmt.Errorf("unable to find role '%s'", roleName)
	}
	// find all assigments per role
	ass, err := srv.RoleAssignments.List(myAccoutsCustomerID).RoleId(strconv.FormatInt(roleID, 10)).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve roles in domain: %v", err)
	}

	principals := []interface{}{}
	emails := []string{}
	for _, each := range ass.Items {
		switch each.AssigneeType {
		case "user":
			usr, err := srv.Users.Get(each.AssignedTo).Do()
			if err == nil {
				emails = append(emails, usr.PrimaryEmail)
				principals = append(principals, usr)
				continue
			}
			sa, err := getServiceAccountByID(c, each.AssignedTo)
			if err == nil {
				emails = append(emails, sa.Email)
				principals = append(principals, sa)
				continue
			}
			log.Printf("unable to retrieve user %s, %s", each.AssignedTo, err)
		case "group":
			grp, err := srv.Groups.Get(each.AssignedTo).Do()
			if err == nil {
				emails = append(emails, grp.Email)
				principals = append(principals, grp)
				continue
			}
			log.Printf("unable to retrieve group %s, %s", each.AssignedTo, err)
		default:
			log.Printf("unknown assignee type %s", each.AssigneeType)
		}
	}
	// find all the emails

	if optionJSON(c, principals) {
		return nil
	}
	for _, email := range emails {
		fmt.Println(email)
	}
	return nil
}

func getServiceAccountByID(c *cli.Context, uniqueId string) (*iam.ServiceAccount, error) {
	client, err := iam.NewService(context.Background(), option.WithHTTPClient(sharedAuthClient(c)))
	if err != nil {
		return nil, err
	}
	return client.Projects.ServiceAccounts.Get("projects/-/serviceAccounts/" + uniqueId).Do()
}
