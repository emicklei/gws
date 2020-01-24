package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
)

func cmdExportGroupMemberships(c *cli.Context) error {
	output := c.Args().Get(0)
	if len(output) == 0 {
		output = fmt.Sprintf("gsuite-group-export-%s.json", time.Now().Format("2006-01-02"))
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

	type entry struct {
		groupKey string
		members  []*admin.Member
	}
	entries := make(chan entry)
	wg := new(sync.WaitGroup)
	for _, g := range r.Groups {
		wg.Add(1)
		// concurrently fetch members of a group
		go func(groupKey string) {
			if c.GlobalBool("v") {
				fmt.Println("fetching members of", groupKey, "...")
			}
			usersList := []*admin.Member{}
			// email is default
			r, err := srv.Members.List(groupKey).Do()
			if err != nil {
				fmt.Printf("unable to retrieve members of group [%s] : %v\n", groupKey, err)
			}
			for _, u := range r.Members {
				if len(u.Email) > 0 {
					// hide internals
					u.Etag = ""
					u.Id = ""
					u.Kind = ""
					u.Role = ""
					u.Type = ""
					usersList = append(usersList, u)
				}
			}
			entries <- entry{groupKey: groupKey, members: usersList}
			wg.Done()
		}(g.Email)
	}
	groupToMembersMap := map[string][]*admin.Member{}
	// collect results
	go func() {
		for each := range entries {
			if c.GlobalBool("v") {
				fmt.Println("... collect members of", each.groupKey)
			}
			groupToMembersMap[each.groupKey] = each.members
		}
	}()

	// wait for all group queries
	wg.Wait()

	// no more results
	close(entries)

	// end spinner
	done()

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(groupToMembersMap)
}
