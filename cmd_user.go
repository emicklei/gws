package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
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
		Customer(myAccoutsCustomerID).
		MaxResults(int64(ifZero(c.Int("limit"), 100))).
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
	if strings.Index(memberKey, "@") == -1 {
		domain, err := primaryDomain()
		if err != nil {
			return err
		}
		memberKey = fmt.Sprintf("%s@%s", memberKey, domain)
	}

	done := showSpinnerWhile(c)

	// get all groups
	if c.GlobalBool("v") {
		log.Println("[gsuite] fetching all groups")
	}
	r, err := srv.Groups.List().
		Customer(myAccoutsCustomerID).
		MaxResults(int64(ifZero(c.Int("limit"), 100))).
		OrderBy("email").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve groups in domain: %v", err)
	}

	checks := make(chan memberCheck, len(r.Groups))
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
				check.callError = fmt.Errorf("aborting! unable to check membership of [email:%s] in [group:%s] because [%v]", memberKey, check.group.Email, err)
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

	wg.Wait()
	close(checks)
	// collect memberships
	membership := []*admin.Group{}
	var abortError error = nil
	for each := range checks {
		if abortError == nil {
			if each.callError != nil {
				abortError = each.callError
			} else {
				if each.isMember {
					membership = append(membership, each.group)
				}
			}
		}
	}
	done() // end spinner
	if abortError != nil {
		return abortError
	}

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
	if strings.Index(userKey, "@") == -1 {
		domain, err := primaryDomain()
		if err != nil {
			return err
		}
		userKey = fmt.Sprintf("%s@%s", userKey, domain)
	}

	r, err := srv.Users.Get(userKey).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve [user:%s] because: %v", userKey, err)
	}
	if optionJSON(c, r) {
		return nil
	}
	fmt.Printf("%s (%s) [tel: %s, 2nd: %s, suspended: %v]\n", r.PrimaryEmail, r.Name.FullName, r.RecoveryPhone, r.RecoveryEmail, r.Suspended)
	return nil
}

func cmdUserAlias(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %v", err)
	}

	userKey := c.Args().Get(0)
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

	r, err := srv.Users.Aliases.List(userKey).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve [user:%s] because: %v", userKey, err)
	}
	if optionJSON(c, r) {
		return nil
	}
	fmt.Printf("%v\n", r)
	return nil
}

func cmdUserSuspend(c *cli.Context) error {
	srv, err := admin.New(sharedAuthClient())
	if err != nil {
		return fmt.Errorf("unable to retrieve directory client %v", err)
	}
	// user
	userKey, err := userKey(c)
	if err != nil {
		return err
	}
	reason := c.Args().Get(1)

	// prompt
	if !promptForYes(c, fmt.Sprintf("Are you sure to suspend user [%s] because [%s] (y/N)? ", userKey, reason)) {
		return errors.New("user suspend aborted")
	}

	user := &admin.User{
		Suspended:        true,
		SuspensionReason: reason,
	}
	_, err = srv.Users.Patch(userKey, user).Do()
	if err != nil {
		return fmt.Errorf("unable to suspend [user:%s] because: %v", userKey, err)
	}
	return nil
}
