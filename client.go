package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
)

var once sync.Once
var sharedClient *http.Client

func sharedAuthClient(c *cli.Context) *http.Client {
	once.Do(func() {
		credfile := c.GlobalString("credentials")
		if len(credfile) == 0 {
			if info, err := os.Stat("gws-credentials.json"); err == nil && !info.IsDir() {
				credfile = "gws-credentials.json"
			}
		}
		if len(credfile) == 0 {
			credfile = filepath.Join(os.Getenv("HOME"), "gws-credentials.json")
		}
		sharedClient = newAuthClient(credfile)
	})
	return sharedClient
}

func newAuthClient(credentialsFilename string) *http.Client {
	b, err := os.ReadFile(credentialsFilename)
	if err != nil {
		abs, _ := filepath.Abs(credentialsFilename)
		log.Fatalf("Unable to read client secret file [%s]: %v", abs, err)
	}

	// After modifying these scopes, delete your previously saved token.json.
	// TODO verify saved token against this list
	// https://developers.google.com/identity/protocols/googlescopes
	// https://developers.google.com/admin-sdk/directory/v1/guides/authorizing
	config, err := google.ConfigFromJSON(b,
		admin.AdminDirectoryRolemanagementReadonlyScope,
		admin.AdminDirectoryDomainReadonlyScope,
		// Create,Delete,Add,Remove
		admin.AdminDirectoryGroupScope,
		admin.AdminDirectoryUserScope,
	)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return getClient(config)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := filepath.Join(os.Getenv("HOME"), "gws-token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving authorisation file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
