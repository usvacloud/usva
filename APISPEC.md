# API Specification

Documentation about endpoints and their usage

| Route                  | Description                  |
| :--------------------- | ---------------------------- |
| [/](#root)             | Contains all API operations  |
| [/file](#file)         | Contains all file operations |
| [/feedback](#feedback) | Contains feedback operations |

### Request authorization

In order to authorize a request, you need to set `Authorization` header to `Bearer <password>` with base64 url-encoded password.

Example of authorization header, where file is locked with password `usva` (this password will be used in following example requests, too):

| Name          | Value          |
| ------------- | -------------- |
| Authorization | Bearer dXN2YQo |

- If authorization fails (e.g. password is invalid), the API returns 403.
- If authorization header doesn't exist even though it's required, the API returns HTTP 401.



## <a name="root">API operations</a>

Contains all API operations

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

#### Examples

Example request

```sh
curl "http://usva.local/restrictions"
```

Example response

```json
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
- [POST /file](#post_file)
- [POST /file/upload](#post_file_upload)
- [POST /file/report](#post-file-report)

### <a name="get_file"> GET /file </a>

**Existing file operation: possibly requires authentication**

Request file's content. `filename` parameter is required.

#### Examples

##### Example request

```sh
curl "http://usva.local/file?filename=5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp" \
	--header "Authorization: Bearer dGFwc2Fpc2Jlc3QK"
```

##### Example response

```
(file content)
```

### <a name="get_file_info">GET /file/info</a>

**Existing file operation: possibly requires authentication**

Request file's information. `filename` parameter is required.

#### Examples

##### Example request

```sh
curl "http://usva.local/file/info?filename=5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp" \
	--header "Authorization: Bearer dGFwc2Fpc2Jlc3QK"
```

##### Example response

```json
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

### <a name="post_file">POST /file </a>

Uploads a file without any parameters and returns the perfect path for querying it's result

#### Examples

##### Example request

```sh
curl --form 'file=@./5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp' http://usva.local/file/
```

##### Example response

```
http://usva.local/file/?filename=bdbd3766-4cc1-46f3-ac28-296d958e848c
```





### <a name="post_file_upload">POST /file/upload </a>

Uploads a file and returns it's filename on server.

##### Required headers:

| Name         | Value     | Description            |
| ------------ | --------- | ---------------------- |
| Content-Type | form-data | Specifies content type |

##### Possible form parameters:

| Name          | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| `title`       | specifies a title for the upload                             |
| `password`    | specifies a password for the upload. this is used for server-side encryption. **note**: has to be in url-encoded base64 format |
| `can_encrypt` | specifies whether file should be encrypted on the server     |
| `file`        | file to upload                                               |



#### Examples

##### Example request:

```sh
curl -L -X POST 'localhost:8080/file/upload' \
    --form 'title="my-upload-title"' \
    --form 'password="some-base64-encoded-string"' \
    --form 'can_encrypt=yes' \
	--form 'file=@"./5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp"'
```

##### Example response:

```json
{
  "message": "file uploaded",
  "filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp"
}
```

### <a name="post-file-report">POST /file/report</a>

#### Examples

##### Example request:

```sh
curl --location --request POST 'http://localhost:8080/file/report' \
	--header 'Content-Type: application/json' \
	--data-raw '{
    	"filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp",
    	"reason": "this file includes copyrighted content!"
	}'
```

##### Example response:

```json
{
    "message": "thank you! your report has been sent."
}
```

## <a name="feedback">Feedback operations</a>

##### Existing routes:

- [POST /feedback](#feedback-add)
- [GET /feedback](#feedback-get)

### <a name="feedback-add">POST /feedback</a>

### Examples

##### Example request:

```sh
curl --location --request POST 'http://localhost:8080/feedback' \
--header 'Content-Type: application/json' \
--data-raw '{
    "comment": "Hello there. Your website is amazing. Great job.",
    "boxes": [
        1, 
        2,
        3
    ]
}'
```

#### Example response:

```json
{
    "message": "Feedback added"
}
```

### <a name="feedback-get">GET /feedback</a>

### Examples

##### Example request:

```sh
curl --location 'http://localhost:8080/feedback'
```

#### Example response:

```json
[
    {
        "comment": "Such a great experience!",
        "boxes": [ 1, 2, 3 ]
    },
    {
        "comment": "I would wish for more smooth workflow",
        "boxes": [ 4, 2, 3, 1 ]
    }
]
```

