package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func cmdDomainList(c *cli.Context) error {

	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Domains.List(myAccoutsCustomerID).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve domains: %v", err)
	}
	if optionJSON(c, r.Domains) {
		return nil
	}
	for _, each := range r.Domains {
		fmt.Println(each.DomainName)
	}
	return nil
}

const primaryDomainEnvironmentKey = "GWS_PRIMARY_DOMAIN"

var cachedPrimaryDomain string

func primaryDomain(c *cli.Context) (string, error) {
	if len(cachedPrimaryDomain) > 0 {
		return cachedPrimaryDomain, nil
	}
	if p := os.Getenv(primaryDomainEnvironmentKey); len(p) > 0 {
		return p, nil
	}
	client := sharedAuthClient(c)

	srv, err := admin.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Domains.List(myAccoutsCustomerID).Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve domains: %v", err)
	}

	for _, each := range r.Domains {
		if each.IsPrimary {
			cachedPrimaryDomain = each.DomainName
			return each.DomainName, nil
		}
	}
	return "", errors.New("no primary domain found")
}
