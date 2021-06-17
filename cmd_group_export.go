package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func cmdExportGroupMemberships(c *cli.Context) error {
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
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
					// hide some internals
					u.Etag = ""
					u.Id = ""
					u.Kind = ""
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

	if optionJSON(c, groupToMembersMap) {
		return nil
	}

	if c.Bool("csv") {
		writeGroupToMembersCSV(groupToMembersMap)
		return nil
	}

	// simple print
	for k, v := range groupToMembersMap {
		fmt.Printf("[%s]\n", k)
		for _, o := range v {
			fmt.Printf("%s\n", o.Email)
		}
	}
	return nil
}

// rows with  {group.email},{role.name},{member.email}
func writeGroupToMembersCSV(groupToMembersMap map[string][]*admin.Member) {
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"group", "role", "member"}) // header
	for group, members := range groupToMembersMap {
		for _, member := range members {
			w.Write([]string{group, strings.ToLower(member.Role), member.Email})
		}
	}
	w.Flush()
}
