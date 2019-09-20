# gsuite - command line tool to use the Google G Suite Admin SDK

## features

- retrieve list of users
- retrieve list of groups
- retrieve membership of a user
- retrieve members of a group

## examples

    gsuite user list
    gsuite user list -limit 1 -format JSON
    gsuite group list
    gsuite user membership john.doe@company.com
    gsuite group members all@company.com

## permissions

Using the tool requires the following authentication scopes to be consent per user.

- https://www.googleapis.com/auth/admin.directory.group.member.readonly
- https://www.googleapis.com/auth/admin.directory.user.readonly

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing

## install from source

This installation requires the Go SDK (1.13+).

    go install github.com/emicklei/gsuite