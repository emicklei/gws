package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
)

func cmdExportGroupMemberships(c *cli.Context) error {
	output := c.Args().Get(0)
	if len(output) == 0 {
		return errors.New("missing output file name, such as group-export-all.json")
	}

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}
	done := showSpinnerWhile(c)
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
	groupToMembersMap := map[string][]string{}
	for _, g := range r.Groups {
		usersList := []string{}
		// email is default
		r, err := srv.Members.List(g.Email).Do()
		if err != nil {
			fmt.Printf("unable to retrieve members of group [%s] : %v\n", g.Email, err)
		}
		for _, u := range r.Members {
			if len(u.Email) > 0 {
				usersList = append(usersList, u.Email)
			}
		}
		groupToMembersMap[g.Email] = usersList
	}
	done() // end spinner

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(groupToMembersMap)
}
