package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
)

func cmdGroupList(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Groups.List().
		Customer(myAccoutsCustomerId).
		MaxResults(int64(IfZero(c.Int("limit"), 100))).
		OrderBy("email").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve groups in domain: %v", err)
	}
	if optionJSON(c, r.Groups) {
		return nil
	}
	for _, g := range r.Groups {
		// email is default
		fmt.Println(g.Email)
	}
	return nil
}

func cmdGroupMembers(c *cli.Context) error {
	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if strings.Index(groupKey, "@") == -1 {
		domain, err := primaryDomain()
		if err != nil {
			return err
		}
		groupKey = fmt.Sprintf("%s@%s", groupKey, domain)
	}

	r, err := srv.Members.List(groupKey).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve members of group [%s] : %v", groupKey, err)
	}
	if optionJSON(c, r.Members) {
		return nil
	}
	for _, g := range r.Members {
		// email is default
		fmt.Println(g.Email)
	}
	return nil
}

func cmdGroupInfo(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %v", err)
	}

	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if strings.Index(groupKey, "@") == -1 {
		domain, err := primaryDomain()
		if err != nil {
			return err
		}
		groupKey = fmt.Sprintf("%s@%s", groupKey, domain)
	}

	r, err := srv.Groups.Get(groupKey).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve group [%s] because: %v", groupKey, err)
	}
	if optionJSON(c, r) {
		return nil
	}
	fmt.Printf("%s (%s) [members=%d]\n", r.Email, r.Name, r.DirectMembersCount)
	return nil
}
