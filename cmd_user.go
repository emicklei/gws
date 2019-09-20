package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
)

func cmdUserList(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Users.List().
		Customer("my_customer"). // ??
		MaxResults(int64(IfZero(c.Int("limit"), 100))).
		OrderBy("email").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve users in domain: %v", err)
	}

	wantsJSON := c.String("format") == "JSON"
	if wantsJSON {
		data, _ := json.MarshalIndent(r.Users, "", "\t")
		fmt.Println(string(data))
		return nil
	}
	for _, u := range r.Users {
		// email is default
		fmt.Println(u.PrimaryEmail)
	}
	return nil
}

func cmdUserMembershipList(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %v", err)
	}

	done := showSpinnerWhile(c)

	// get all groups
	if c.GlobalBool("v") {
		log.Println("[gsuite] fetching all groups")
	}
	r, err := srv.Groups.List().
		Customer("my_customer"). // ??
		MaxResults(int64(IfZero(c.Int("limit"), 100))).
		OrderBy("email").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve groups in domain: %v", err)
	}

	// visit all groups and check membership
	memberKey := c.Args().Get(0)
	if len(memberKey) == 0 {
		return fmt.Errorf("missing user email in command")
	}

	membership := []*admin.Group{}
	for _, g := range r.Groups {
		// Email or immutable ID of the group
		groupKey := g.Id
		if c.GlobalBool("v") {
			log.Printf("[gsuite] is %s member of group %s ?\n", memberKey, g.Email)
		}
		hasResult, err := srv.Members.HasMember(groupKey, memberKey).Do()
		if err != nil {
			return fmt.Errorf("unable to check membership of [%s] in [%s] because [%v]", memberKey, groupKey, err)
		}
		if hasResult.IsMember {
			membership = append(membership, g)
		}
	}
	done()

	for _, g := range membership {
		// email is default
		fmt.Println(g.Email)
	}
	return nil
}
