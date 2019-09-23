# gsuite - command line tool to use the Google G Suite Admin SDK

[![Build Status](https://travis-ci.org/emicklei/gsuite.png)](https://travis-ci.org/emicklei/gsuite)

## features

- retrieve list of users
- show details of user
- retrieve membership of a user
- retrieve list of groups
- retrieve members of a group
- show details of group
- retrieve list of roles
- any command can produce JSON format using `-json` at the end of the command

## examples

    gsuite user list
    gsuite user list -limit 4
    gsuite user membership john.doe@company.com
    gsuite user info john.doe@company.com

    gsuite group list    
    gsuite group members all@company.com
    gsuite group info all@company.com
    
    gsuite role assignments _SEED_ADMIN_ROLE

## requirements

- A G Suite domain with API access enabled
- A Google account in that domain with administrator privileges
- A Google Cloud Platform project with Admin SDK enabled ( https://console.developers.google.com/apis/library/admin.googleapis.com?project=YOURPROJECT )


## tool authentication

- Using the Google Cloud Platform console, create a new OAuth 2.0 client ID credential in the project for which you enabled the Admin SDK.
- Download the JSON file from the list of Credentials (download button on the right).
- Save the file to *gsuite-credentials.json* in your *home* directory.

## user permissions

*gsuite* requires the following authentication scopes to be consent per user.
You will be asked to accept those on the first time you use *gsuite*.

- https://www.googleapis.com/auth/admin.directory.group.member.readonly
- https://www.googleapis.com/auth/admin.directory.user.readonly
- https://www.googleapis.com/auth/admin.directory.rolemanagement.readonly

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing

## install from source

This installation requires the Go SDK (1.13+).

    go install github.com/emicklei/gsuite

## help

Have problems using *gsuite* ? Read about [known errors](/errors.md)

&copy; 2019, ernestmicklei.com. MIT License. Contributions welcome.
