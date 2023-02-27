# API Specification

Documentation about endpoints and their usage

| Route                  | Description                                        |
| :--------------------- | -------------------------------------------------- |
| [/](#root)             | Contains all general query operations              |
| [/file](#file)         | Contains all file operations                       |
| [/feedback](#feedback) | Contains feedback operations                       |
| [/account](#account)   | Contains all account management related operations |



### Request authorization

In order to authorize a request, you need to set `Authorization` header to `Bearer <password>` with base64 url-encoded password.

Example of authorization header, where file is locked with password `usva` (this password will be used in this documentation's example requests, too):

| Name          | Value          |
| ------------- | -------------- |
| Authorization | Bearer dXN2YQo |

- If authorization fails (e.g. password is invalid), the API returns 403.
- If authorization header doesn't exist even though it's required, the API returns HTTP 401.



## <a name="root">General/informing operations</a>

Contains all informing endpoints

### Existing routes

- [GET /restrictions](#restrc)



### <a name="restrc">GET /restrictions</a>

Returns API restrictions. These can be for example shown on client.

#### Fields

| Field name           | Description                                               |
| :------------------- | --------------------------------------------------------- |
| maxDailyUploadSize   | Maximum total size of uploaded files in a single day      |
| filePersistDuration  | Describes how long a file is saved (days, hours, seconds) |
| maxEncryptedFileSize | Maximum size for an server-side encrypted file            |
| maxSingleUploadSize  | Maximum size for an non-server-side encrypted file        |

```sh
> curl -L "usva.local/restrictions" | jq
{
  "filePersistDuration": {
    "days": 1,
    "hours": 24,
    "seconds": 86400
  },
  "maxDailyUploadSize": {
    "bytes": 0,
    "gigabytes": 0,
    "kilobytes": 0,
    "megabytes": 0
  },
  "maxEncryptedFileSize": {
    "bytes": 100000000,
    "gigabytes": 0,
    "kilobytes": 100000,
    "megabytes": 100
  },
  "maxSingleUploadSize": {
    "bytes": 0,
    "gigabytes": 0,
    "kilobytes": 0,
    "megabytes": 0
  }
}
```



## <a name="file">File operations</a>

Contains all file operations

##### Existing routes:

- [GET /file](#get_file)

- [GET /file/info](#get_file_info)

- [POST /file/upload](#post_file_upload)

- [POST /file/report](#post-file-report)

  

### <a name="get_file"> GET /file </a>

```sh
> curl -L "usva.local/file?filename=5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp" \
	-H "Authorization: Bearer dGFwc2Fpc2Jlc3QK" \
	-o-
(file content)	
```



### <a name="get_file_info">GET /file/info</a>

```sh
> curl "http://usva.local/file/info?filename=5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp" \
	-H "Authorization: Bearer dGFwc2Fpc2Jlc3QK"
{
  "encrypted": false,
  "filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp",
  "locked": true,
  "size": 8389748, 
  "title": {
    "String": "",
    "Valid": false
  },
  "uploadDate": "2023-02-17T10:33:56.493197Z",
  "viewCount": 0
}
```

### 

### <a name="post_file_upload">POST /file/upload</a>

| Name          | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| `title`       | specifies a title for the upload                             |
| `password`    | specifies a password for the upload. this is used for server-side encryption. **note**: has to be in url-encoded base64 format |
| `can_encrypt` | specifies whether file should be encrypted on the server     |
| `file`        | file to upload                                               |



```sh
> curl -L 'localhost:8080/file/upload' \
    --form 'title="my-upload-title"' \
    --form 'password="some-base64-encoded-string"' \
    --form 'can_encrypt=yes' \
	--form 'file=@./file.pgp'
{
  "message": "file uploaded",
  "filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp"
}
```



### <a name="post-file-report">POST /file/report</a>

```sh
> curl -LX POST 'http://localhost:8080/file/report' \
	-H 'Content-Type: application/json' \
	-d '{"filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp", "reason": "somethings wrong"}' | jq
{
    "message": "thank you! your report has been sent."
}
```



## <a name="feedback">Feedback operations</a>

##### Existing routes:

- [POST /feedback](#feedback-add)
- [GET /feedback](#feedback-get)



### <a name="feedback-add">POST /feedback</a>

Send a new feedback

```sh
> curl -LX POST 'http://localhost:8080/feedback' \
	-H 'Content-Type: application/json' \
	-d '{ "message": "mycomment", "boxes": [ 1, 2, 3 ] }' | jq
{
    "message": "Feedback added"
}
```



### <a name="feedback-get">GET /feedback</a>

```sh
> curl -L 'http://localhost:8080/feedback' | jq
[
    {
        "comment": {
        	"String": "Such a great experience!",
        	"Valid": true
        },
        "boxes": [ 1, 2, 3 ]
    },
    {
        "comment": {
        	"String": "",
        	"Valid": false
        },
        "boxes": [ 4, 2, 3, 1 ]
    }
]
```



## <a name="account">Account operations</a>

Contains all file operations

##### Existing routes:

- [GET /account](#get_account)
- [POST /login](#account_login)
- [POST /register](#account_register)
- [GET /account/files](#get_account_files)
- [GET /account/files/all](#get_all_account_files)

- [GET /sessions](#get_sessions)
- [DELETE /sessions](#account_delete_session)
- [DELETE /sessions/all](#account_delete_sessions)



### <a id="get_account">GET /account</a>

Get your current session profile. Session is read from cookies.

```sh
> curl -b mycookiefile -L "usva.local/account" | jq 		# jq: command line json parsing tool
{
  "token": "KULgs34BJjis_jRBysts84afxyg", 					# this is the same as your session token 
  "account": { 												# this includes your profile information
    "account_id": "7e08549f-44e1-4b97-b9dc-864f9f8fc5ca",
    "username": "toke",
    "register_date": "2023-02-21T01:00:48.617262Z",
    "last_login": "2023-02-21T01:00:48.617262Z",
    "activity_points": 0									# how many files are linked to your account 
  }
}
```



### <a id="account_login">POST /account/login</a>

Create a new session. Session is read from cookies.

```sh
> curl -c mycookiefile -d '{"username": "myuser", "password": "mypassword"}' -L "usva.local/account" | jq
{
  "sessionId": "KULgs34BJjis_jRBysts84afxyg", # session token. this is also saved to cookies. 
}
```



### <a id="account_login">POST /account/register</a>

Takes in the same parameters as [/account/login](#account_login).
This path does not implement existing account authentication logic (e.g. you can not use this for logging in to existing account).



### <a id="get_account_files">GET /account/files</a>

Get files that are linked to your profile, you will probably find this the only reason to use accounts

```sh
> curl -b mycookiefile -L "usva.local/account/files/?limit=1" # limit <= 10 (not required)
{
  "files": [
    {
      "filename": "<random uuid>",
      "title": {
        "String": "",
        "Valid": false
      },
      "file_size": 420,
      "viewcount": 69,
      "encrypted": true,
      "upload_date": "2023-02-27T13:25:39.778869Z",
      "last_seen": "2023-02-27T13:25:39.778869Z"
    }
  ]
}
```



### <a id="get_all_account_files">GET /account/files/all</a>

Get all owned files. Does not include anything else than filename. 

```sh
> curl -b mycookiefile -L "usva.local/account/files/all"
{
  "files": [
  	"<file uuid>",
  	"<some other file uuid>"
  ]
}
```



### <a id="get_sessions">GET /account/sessions</a>

Show all user sessions

```sh
> curl -b mycookiefile -L "usva.local/account/sessions"
{
  "sessions": [
    {
      "session_id": "KULgs34BJjis_jRBysts84afxyg",
      "start_date": "2023-02-27T13:15:37.053267+02:00"
    }
  ]
}
```



### <a id="account_delete_session">DELETE /account/sessions</a>

Delete a single session

```sh
> curl \
	-b cookiejar \
	-X DELETE \
	-d '{"token": "KULgs34BJjis_jRBysts84afxyg"}' \
	-L "localhost:8080/account/sessions" | jq
{ 
	"message": "ok"
}
```



### <a id="account_delete_sessions">DELETE /account/sessions/all</a>

Delete a single session

```sh
> curl -b cookiejar -X DELETE -L "localhost:8080/account/sessions/all" | jq
{
  "message": "ok",
  "removed": [
  	"KULgs34BJjis_jRBysts84afxyg",
  ]
}
```

