# gdom - command line tool to use the Google G Suite Admin SDK

[![Build Status](https://travis-ci.org/emicklei/gsuite.png)](https://travis-ci.org/emicklei/gsuite)

## features

- retrieve list of users
- retrieve list of groups
- retrieve membership of a user
- retrieve members of a group
- show details of user
- show details of group
- any command can produce JSON format using `-format JSON` at the end of the command

## examples

    gdom user list
    gdom user list -limit 4
    gdom user membership john.doe@company.com
    gdom user info john.doe@company.com

    gdom group list    
    gdom group members all@company.com
    gdom group info all@company.com
    
    gdom reset

## requirements (TOWRITE)

- A G Suite domain with API access enabled
- A Google account in that domain with administrator privileges
- A Google Cloud Platform project with Directory API enabled

[missinglink]

In resulting dialog click DOWNLOAD CLIENT CONFIGURATION and save the file *credentials.json* to your *home* directory.

## permissions

Using the tool requires the following authentication scopes to be consent per user.

- https://www.googleapis.com/auth/admin.directory.group.member.readonly
- https://www.googleapis.com/auth/admin.directory.user.readonly

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing

## install from source

This installation requires the Go SDK (1.13+).

    go install github.com/emicklei/gdom

&copy; 2019, ernestmicklei.com. MIT License. Contributions welcome.
