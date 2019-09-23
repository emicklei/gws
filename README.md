# gdom - command line tool to use the Google G Suite Admin SDK

[![Build Status](https://travis-ci.org/emicklei/gdom.png)](https://travis-ci.org/emicklei/gdom)

## features

- retrieve list of users
- retrieve list of groups
- retrieve membership of a user
- retrieve members of a group
- show details of user
- show details of group
- any command can produce JSON format using `-json` at the end of the command

## examples

    gdom user list
    gdom user list -limit 4
    gdom user membership john.doe@company.com
    gdom user info john.doe@company.com

    gdom group list    
    gdom group members all@company.com
    gdom group info all@company.com
    
    gdom reset

## requirements

- A G Suite domain with API access enabled
- A Google account in that domain with administrator privileges
- A Google Cloud Platform project with Admin SDK enabled ( https://console.developers.google.com/apis/library/admin.googleapis.com?project=YOURPROJECT )


## tool authentication

- Using the Google Cloud Platform console, create a new OAuth 2.0 client ID credential in the project for which you enabled the Admin SDK.
- Download the JSON file from the list of Credentials (download button on the right).
- Save the file to *gdom-credentials.json* in your *home* directory.

## user permissions

*gdom* requires the following authentication scopes to be consent per user.
You will be asked to accept those on the first time you use *gdom*.

- https://www.googleapis.com/auth/admin.directory.group.member.readonly
- https://www.googleapis.com/auth/admin.directory.user.readonly

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing

## install from source

This installation requires the Go SDK (1.13+).

    go install github.com/emicklei/gdom

&copy; 2019, ernestmicklei.com. MIT License. Contributions welcome.
