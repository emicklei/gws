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
