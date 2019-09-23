# Troubleshooting


### 1. Access Not Configured

    unable to retrieve users in domain: googleapi: Error 403: Access Not Configured. Admin Directory API has not been used in project xxxxxxxx before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/admin.googleapis.com/overview?project=xxxxxxxx then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry., accessNotConfigured

In your current Google Cloud Platform project, enable the Admin SDK api. Then create an OAuth 2.0 client id credential and download the secret to your home folder.

### 2. Not Authorized

    unable to retrieve users in domain: googleapi: Error 403: Not Authorized to access this resource/api, forbidden

In G Suite, make sure the you have the permissions to read users and groups.


### 3. Insufficient Permission

    2019/09/23 12:08:01 unable to retrieve roles in domain: googleapi: Error 403: Insufficient Permission: Request had insufficient authentication scopes., insufficientPermissions

You do not have an (updated) authorisation file in your home directory. Remove the file *gsuite-token.json* and retry.