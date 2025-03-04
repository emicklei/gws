package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func cmdGroupList(c *cli.Context) error {

	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	call := srv.Groups.List().
		Customer(myAccoutsCustomerID).
		MaxResults(int64(ifZero(c.Int("limit"), 100))).
		OrderBy("email")
	if domain := c.GlobalString("domain"); len(domain) > 0 {
		call = call.Domain(domain)
	}
	r, err := call.Do()
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
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if !strings.Contains(groupKey, "@") {
		domain, err := primaryDomain(c)
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

	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w", err)
	}

	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if !strings.Contains(groupKey, "@") {
		domain, err := primaryDomain(c)
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
	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(sharedAuthClient(c)))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if !strings.Contains(groupKey, "@") {
		domain, err := primaryDomain(c)
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
	fmt.Printf("deleted group [%s]\n", groupKey)
	return nil
}

func cmdGroupAddMembers(c *cli.Context) error {
	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(sharedAuthClient(c)))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if !strings.Contains(groupKey, "@") {
		domain, err := primaryDomain(c)
		if err != nil {
			return err
		}
		groupKey = fmt.Sprintf("%s@%s", groupKey, domain)
	}
	// users argument
	userArgument := c.Args().Get(1)
	if len(userArgument) == 0 {
		return fmt.Errorf("missing (space separated) user(s) email in command")
	}
	userKeys := []string{}
	for _, each := range c.Args()[1:] {
		userKey := each
		if !strings.Contains(each, "@") {
			domain, err := primaryDomain(c)
			if err != nil {
				return err
			}
			userKey = fmt.Sprintf("%s@%s", each, domain)
		}
		userKeys = append(userKeys, userKey)
	}
	// prompt
	if !promptForYes(c, fmt.Sprintf("Are you sure to add member(s) %v to group [%s] (y/N)? ", userKeys, groupKey)) {
		return errors.New("group add aborted")
	}
	for _, each := range userKeys {
		// doit
		member := &admin.Member{
			Email: each,
		}
		_, err = srv.Members.Insert(groupKey, member).Do()
		if err != nil {
			// TODO: show different log when *googleapi.Error=&{409 Member already exists.
			fmt.Printf("unable to add member [%s] to group [%s] because: %v (%T)\n", each, groupKey, err, err)
			continue
		}
		fmt.Printf("added member [%s] to group [%s]\n", each, groupKey)
	}
	return nil
}

func cmdGroupRemoveMember(c *cli.Context) error {
	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(sharedAuthClient(c)))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if !strings.Contains(groupKey, "@") {
		domain, err := primaryDomain(c)
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
	if !strings.Contains(userKey, "@") {
		domain, err := primaryDomain(c)
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
	fmt.Printf("removed member [%s] from group [%s]\n", userKey, groupKey)
	return nil
}

func cmdGroupCreate(c *cli.Context) error {
	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(sharedAuthClient(c)))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %w (%T)", err, err)
	}
	// group argument
	groupKey := c.Args().Get(0)
	if len(groupKey) == 0 {
		return fmt.Errorf("missing group email in command")
	}
	if !strings.Contains(groupKey, "@") {
		domain, err := primaryDomain(c)
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
	fmt.Printf("created group [%s]\n", groupKey)
	return nil
}
