package main

import (
	"encoding/json"
	"fmt"

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
		Customer("my_customer"). // ??
		MaxResults(int64(IfZero(c.Int("limit"), 100))).
		OrderBy("email").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve groups in domain: %v", err)
	}

	wantsJSON := c.String("format") == "JSON"
	if wantsJSON {
		data, _ := json.MarshalIndent(r.Groups, "", "\t")
		fmt.Println(string(data))
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
	r, err := srv.Members.List(groupKey).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve members of group: %v", err)
	}
	for _, g := range r.Members {
		// email is default
		fmt.Println(g.Email)
	}
	return nil
}
