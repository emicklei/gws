package main

import (
	"errors"
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
		return fmt.Errorf("unable to retrieve directory client %w", err)
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
		return fmt.Errorf("unable to retrieve group [%s] because: %v (%T)", groupKey, err, err)
	}
	if optionJSON(c, r) {
		return nil
	}
	fmt.Printf("%s (%s) [members=%d]\n", r.Email, r.Name, r.DirectMembersCount)
	return nil
}

func cmdGroupDelete(c *cli.Context) error {
	srv, err := admin.New(sharedAuthClient())
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
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
	// prompt
	if !promptForYes(c, fmt.Sprintf("Are you sure to delete group [%s] (y/N)? ", groupKey)) {
		return errors.New("group delete aborted")
	}
	// doit
	err = srv.Groups.Delete(groupKey).Do()
	if err != nil {
		return fmt.Errorf("unable to delete group [%s] because: %w (%T)", groupKey, err, err)
	}
	fmt.Printf("deleted group [%s]/n", groupKey)
	return nil
}

func cmdGroupAddMember(c *cli.Context) error {
	srv, err := admin.New(sharedAuthClient())
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
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
	// user argument
	userKey := c.Args().Get(1)
	if len(userKey) == 0 {
		return fmt.Errorf("missing user email in command")
	}
	if strings.Index(userKey, "@") == -1 {
		domain, err := primaryDomain()
		if err != nil {
			return err
		}
		userKey = fmt.Sprintf("%s@%s", userKey, domain)
	}
	// prompt
	if !promptForYes(c, fmt.Sprintf("Are you sure to add member [%s] to group [%s] (y/N)? ", userKey, groupKey)) {
		return errors.New("group add aborted")
	}
	// doit
	member := &admin.Member{
		Email: userKey,
	}
	_, err = srv.Members.Insert(groupKey, member).Do()
	if err != nil {
		return fmt.Errorf("unable to add member [%s] to group [%s] because: %w (%T)", userKey, groupKey, err, err)
	}
	fmt.Printf("added member [%s] to group [%s]/n", userKey, groupKey)
	return nil
}

func cmdGroupRemoveMember(c *cli.Context) error {
	srv, err := admin.New(sharedAuthClient())
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
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
	// member argument
	userKey := c.Args().Get(1)
	if len(userKey) == 0 {
		return fmt.Errorf("missing user email in command")
	}
	if strings.Index(userKey, "@") == -1 {
		domain, err := primaryDomain()
		if err != nil {
			return err
		}
		userKey = fmt.Sprintf("%s@%s", userKey, domain)
	}
	// prompt
	if !promptForYes(c, fmt.Sprintf("Are you sure to remove member [%s] from group [%s] (y/N)? ", userKey, groupKey)) {
		return errors.New("group remove aborted")
	}
	// doit
	err = srv.Members.Delete(groupKey, userKey).Do()
	if err != nil {
		return fmt.Errorf("unable to remove member [%s] from group [%s] because: %w (%T)", userKey, groupKey, err, err)
	}
	fmt.Printf("removed member [%s] from group [%s]/n", userKey, groupKey)
	return nil
}

func cmdGroupCreate(c *cli.Context) error {
	srv, err := admin.New(sharedAuthClient())
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
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
	// doit
	grp := &admin.Group{
		Email: groupKey,
	}
	r, err := srv.Groups.Insert(grp).Do()
	if err != nil {
		return fmt.Errorf("unable to create group [%s] because: %w (%T)", groupKey, err, err)
	}
	if optionJSON(c, r) {
		return nil
	}
	fmt.Printf("created group [%s]/n", groupKey)
	return nil
}
