# gsuite - command line tool to use the Google G Suite Admin SDK

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

    gsuite user list
    gsuite user list -limit 4
    gsuite user membership john.doe@company.com
    gsuite user info john.doe@company.com

    gsuite group list    
    gsuite group members all@company.com
    gsuite group info all@company.com

## permissions

Using the tool requires the following authentication scopes to be consent per user.

- https://www.googleapis.com/auth/admin.directory.group.member.readonly
- https://www.googleapis.com/auth/admin.directory.user.readonly

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing

## install from source

This installation requires the Go SDK (1.13+).

    go install github.com/emicklei/gsuite

&copy; 2019, ernestmicklei.com. MIT License. Contributions welcome.
