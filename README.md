# gsuite - command line tool to use the Google G Suite Admin SDK

[![Build Status](https://travis-ci.org/emicklei/gsuite.png)](https://travis-ci.org/emicklei/gsuite)

## features

- list of users
- details of a user
- membership of a user
- management of groups
- list of roles
- user assignments of a role
- list of domains

Any command can produce JSON format using `-json` at the end of the command.

## examples

    gsuite user list
    gsuite user list -limit 4
    gsuite user membership john.doe
    gsuite user membership john.doe@company.com
    gsuite user info john.doe
    gsuite user info john.doe@company.com
    gsuite user aliases john.doe@company.com
    gsuite user suspend angelina "retired"

    gsuite group list    
    gsuite group members all
    gsuite group members all@company.com

    gsuite group info somegroup
    gsuite group info somegroup@company.com
    gsuite group delete my-old@company.com
    gsuite group delete my-old@company.com
    gsuite group add my-group this-person other-person@company.com
    gsuite group remove my-group this-person
    gsuite group export -json > all.json    
    gsuite --domain company.com group export -csv > company-only.csv

    gsuite role list
    gsuite role assignments _USER_MANAGEMENT_ADMIN_ROLE
   
    gsuite domain list

    gsuite examples

## requirements

- A Google Cloud Identity domain with API access enabled
- A Google account in that domain with enough administrator privileges
- A Google Cloud Platform project with Admin SDK enabled ( https://console.developers.google.com/apis/library/admin.googleapis.com?project=YOURPROJECT )

## install

Installation requires the Go SDK.

    go install github.com/emicklei/gsuite@latest 

## tool authentication

- Using the Google Cloud Platform console, create a new OAuth 2.0 client ID credential in the project for which you enabled the Admin SDK.
- Download the JSON file from the list of Credentials (download button on the right).
- Save the file to *gsuite-credentials.json* in your *home* directory or a *local* directory if you need access to more organisations. *gsuite* will look for this file in the current directoy first.

## user permissions

*gsuite* requires the following authentication scopes to be consent per user.
You will be asked to accept those on the first time you use *gsuite*.
Note that accepting these scopes does not mean you as a user have access ; this is controlled in Cloud Identity (or GSuite) Admin Console.

- https://www.googleapis.com/auth/admin.directory.user
- https://www.googleapis.com/auth/admin.directory.group ( for group management )
- https://www.googleapis.com/auth/admin.directory.rolemanagement.readonly
- https://www.googleapis.com/auth/admin.directory.domain.readonly

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing

## install from source

This installation requires the Go SDK (1.13+).

    go install github.com/emicklei/gsuite

## help

Having problems using *gsuite* ? Read about [known errors](/errors.md)

&copy; 2019+, ernestmicklei.com. MIT License. Contributions welcome.
