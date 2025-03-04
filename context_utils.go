package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func userKey(c *cli.Context) (string, error) {
	userKey := c.Args().Get(0)
	if len(userKey) == 0 {
		return "", fmt.Errorf("missing user email in command")
	}
	if !strings.Contains(userKey, "@") {
		domain, err := primaryDomain(c)
		if err != nil {
			return "", err
		}
		userKey = fmt.Sprintf("%s@%s", userKey, domain)
	}
	return userKey, nil
}
