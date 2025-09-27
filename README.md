# gws - command line tool to use the Google Workspace Admin SDK

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

    gws user list
    gws user list -limit 4
    gws user membership john.doe
    gws user membership john.doe@company.com
    gws user info john.doe
    gws user info john.doe@company.com
    gws user aliases john.doe@company.com
    gws user suspend angelina "retired"

    gws group list    
    gws group members all
    gws group members all@company.com

    gws group info somegroup
    gws group info somegroup@company.com
    gws group delete my-old@company.com
    gws group delete my-old@company.com
    gws group add my-group this-person other-person@company.com
    gws group remove my-group this-person
    gws group export -json > all.json    
    gws --domain company.com group export -csv > company-only.csv

    gws role list
    gws role assignments _USER_MANAGEMENT_ADMIN_ROLE
   
    gws domain list

    gws examples

## requirements

- A Google Cloud Identity domain with API access enabled
- A Google account in that domain with enough administrator privileges
- A Google Cloud Platform project with Admin SDK enabled ( https://console.developers.google.com/apis/library/admin.googleapis.com?project=YOURPROJECT )

### primary domain access

If your Google Workspace (GSuite) account only has "Group Editor" role then you cannot use the short syntax for accounts that require the lookup of the primary domain. You can workaround this missing permission by setting an enviroment variable such as:

    export GWS_PRIMARY_DOMAIN=yourhost.com

## install

Installation requires the Go SDK.

    go install github.com/emicklei/gws@latest 

## tool authentication

There are two ways to authenticate yourself. One through is through the [installed apps](https://developers.google.com/identity/protocols/oauth2#installed) flow,
the other is through the [service accounts](https://developers.google.com/identity/protocols/oauth2/service-account) flow. The first
flow requires you to create a new OAuth 2.0 client ID credential in the project and store these and the user credentials files on your file system.
The latter does not keep any credentials stored locally.

## installed application authentication flow

- Using the Google Cloud Platform console, create a new OAuth 2.0 client ID credential in the project for which you enabled the Admin SDK.
- Download the JSON file from the list of Credentials (download button on the right).
- Save the file to *gws-credentials.json* in your *home* directory or a *local* directory if you need access to more organisations. *gws* will look for this file in the current directoy first.

## user permissions

*gws* requires the following authentication scopes to be consent per user.
You will be asked to accept those on the first time you use *gws*.
Note that accepting these scopes does not mean you as a user have access; this is controlled in Cloud Identity (or Google Workspace/GSuite) Admin Console.

- https://www.googleapis.com/auth/admin.directory.user
- https://www.googleapis.com/auth/admin.directory.group (for group management)
- https://www.googleapis.com/auth/admin.directory.rolemanagement.readonly
- https://www.googleapis.com/auth/admin.directory.domain.readonly
- https://www.googleapis.com/auth/iam (for service account lookup)

See also https://developers.google.com/admin-sdk/directory/v1/guides/authorizing


## service account authentication flow

- Create a service account in your Google Cloud project 
  - grant the roles iam.serviceAccountTokenCreator and iam.serviceAccountUser to the service account itself
  - grant the roles iam.serviceAccountTokenCreator and iam.serviceAccountUser to the users you want to use *gws* with
- Grant the service account the role iam.serviceAccountViewer in your Google organization
- Delegate domain-wide authority to the service account
  - follow the instructions in the [documentation](https://developers.google.com/identity/protocols/oauth2/service-account#delegatingauthority)
  - using the following scopes for read-only access:
    - https://www.googleapis.com/auth/admin.directory.domain.readonly
    - https://www.googleapis.com/auth/admin.directory.rolemanagement.readonly 
    - https://www.googleapis.com/auth/admin.directory.group.readonly
    - https://www.googleapis.com/auth/admin.directory.user.readonly
    - https://www.googleapis.com/auth/iam

  - or the following scopes for read-write access:
      - https://www.googleapis.com/auth/admin.directory.domain.readonly
      - https://www.googleapis.com/auth/admin.directory.rolemanagement
      - https://www.googleapis.com/auth/admin.directory.group
      - https://www.googleapis.com/auth/admin.directory.user
      - https://www.googleapis.com/auth/iam
  
- set the environment variable GWS_SERVICE_ACCOUNT to the email address of the created service account
- set the environment variable GWS_ADMIN_USER to the email address of the user you want to impersonate

## help

Having problems using *gws* ? Read about [known errors](/errors.md)

&copy; 2023+, ernestmicklei.com. MIT License. Contributions welcome.
