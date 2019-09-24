package main

import (
	"fmt"
	"log"
	"sync"

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

	if optionJSON(c, r.Users) {
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

	// visit all groups and check membership
	memberKey := c.Args().Get(0)
	if len(memberKey) == 0 {
		return fmt.Errorf("missing user email in command")
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

	checks := make(chan memberCheck)
	wg := new(sync.WaitGroup)

	for _, g := range r.Groups {
		if c.GlobalBool("v") {
			log.Printf("[gsuite] is %s member of group %s ?\n", memberKey, g.Email)
		}
		wg.Add(1)
		go func(check memberCheck) {
			// Use Email or immutable ID of the group
			hasResult, err := srv.Members.HasMember(check.group.Id, check.memberKey).Do()
			if err != nil {
				check.callError = fmt.Errorf("unable to check membership of [%s] in [%s] because [%v]", memberKey, check.group.Id, err)
			} else {
				check.isMember = hasResult.IsMember
			}
			checks <- check
			wg.Done()
		}(memberCheck{
			memberKey: memberKey,
			group:     g,
		})
	}
	// collect memberships
	membership := []*admin.Group{}
	go func() {
		for each := range checks {
			if each.callError != nil {
				log.Println(each.callError)
			} else {
				if each.isMember {
					membership = append(membership, each.group)
				}
			}
		}
	}()
	wg.Wait()
	close(checks)

	done() // end spinner

	if optionJSON(c, membership) {
		return nil
	}
	for _, g := range membership {
		// email is default
		fmt.Println(g.Email)
	}
	return nil
}

type memberCheck struct {
	memberKey string
	isMember  bool
	group     *admin.Group
	callError error
}

func cmdUserInfo(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %v", err)
	}

	userKey := c.Args().Get(0)
	if len(userKey) == 0 {
		return fmt.Errorf("missing user email in command")
	}

	r, err := srv.Users.Get(userKey).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve user [%s] because: %v", userKey, err)
	}
	if optionJSON(c, r) {
		return nil
	}
	fmt.Printf("%s (%s) [tel: %s, 2nd: %s, suspended: %v]\n", r.PrimaryEmail, r.Name.FullName, r.RecoveryPhone, r.RecoveryEmail, r.Suspended)
	return nil
}
