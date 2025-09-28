package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

// Get delegated token for the user, through impersonation the service account if specified
func getDelegatedUserToken(ctx context.Context, user string, serviceAccount string, scopes []string) (oauth2.TokenSource, error) {
	defaultCredentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, err
	}

	caller, err := getUserInfo(ctx, defaultCredentials)
	if err != nil {
		return nil, err
	}

	if serviceAccount == "" {
		if !strings.HasSuffix(caller.Email, ".m.gserviceaccount.com") {
			return nil, fmt.Errorf("service account is required for impersonating %s", user)
		}
		serviceAccount = caller.Email
	}

	iam, err := iamcredentials.NewService(ctx, option.WithCredentials(defaultCredentials))
	if err != nil {
		return nil, err
	}

	iat := time.Now().Unix()
	exp := iat + 3600

	claim, err := json.Marshal(map[string]interface{}{
		"iss":   serviceAccount,
		"scope": strings.Join(scopes, " "),
		"aud":   "https://accounts.google.com/o/oauth2/token",
		"sub":   user,
		"iat":   iat,
		"exp":   exp,
	})

	// Sign the JWT
	name := fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount)
	signReq := &iamcredentials.SignJwtRequest{
		Payload:   string(claim),
		Delegates: []string{},
	}
	signResp, err := iam.Projects.ServiceAccounts.SignJwt(name, signReq).Do()
	if err != nil {
		return nil, err
	}

	// Exchange JWT for OAuth2 token
	tokenResp, err := http.PostForm("https://accounts.google.com/o/oauth2/token", map[string][]string{
		"grant_type":     {"assertion"},
		"assertion_type": {"http://oauth.net/grant_type/jwt/1.0/bearer"},
		"assertion":      {signResp.SignedJwt},
	})
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(tokenResp.Body)

	body, _ := io.ReadAll(tokenResp.Body)
	if tokenResp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed: %s", body)
	}
	var tokenData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return nil, err
	}

	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: tokenData.AccessToken,
		TokenType:   tokenData.TokenType,
		Expiry:      time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second),
	}), nil
}

type userInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}

// get user information of the owner of the credentials
func getUserInfo(ctx context.Context, credentials *google.Credentials) (*userInfo, error) {
	token, err := credentials.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get userinfo: %s, %s", resp.Status, body)
	}

	result := userInfo{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
