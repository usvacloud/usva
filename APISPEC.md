# API Specification

Documentation about endpoints and their usage

| Route          | Description                  |
| -------------- | ---------------------------- |
| [/file](#file) | Contains all file operations |

### Request authorization

In order to authorize a request, you need to set `Authorization` header to `Bearer <password>` with base64 encoded password.  

Example of authorization header, where file is locked with password `usvaisbest` (this example will be used in sample requests):

| Name          | Value                   | Description |
| ------------- | ----------------------- | ----------- |
| Authorization | Bearer dGFwc2Fpc2Jlc3QK |             |

- If authorization fails (e.g. password is invalid), the API returns 403.
- If authorization header doesn't exist even though it's required, the API returns HTTP 401.



## <a name="file">Files</a>

### Existing routes

- [GET /file](#get_file)
- [GET /file/info](#get_file_info)
- [POST /file/upload](#post_file_upload)



### <a name="get_file"> GET /file </a>

**Existing file operation: possibly requires authentication**

Request file's content. `filename` param is required.

#### Examples

##### Sample request

```sh
curl "http://usva.local/file?filename=5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp" \
	--header "Authorization: Bearer dGFwc2Fpc2Jlc3QK"
```

##### Sample response

```json
<file content>
```



### <a name="get_file_info">GET /file/info</a>

Request file's information. `filename` param is required.

##### Required headers: none

#### Examples

##### Sample request

```sh
curl "http://usva.local/file/info?filename=5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp" \
	--header "Authorization: Bearer dGFwc2Fpc2Jlc3QK"
```

##### Sample response

```json
{
    "filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp",
    "locked": true,
    "size": 6757,
    "uploadDate": "1970-01-01T00:00:00+03:00",
    "viewCount": 10
}
```



### <a name="post_file_upload">POST /file/upload </a>

Uploads a file and returns it's filename on server.

##### Required headers:

| Name         | Valuew    | Description            |
| ------------ | --------- | ---------------------- |
| Content-Type | form-data | Specifies content type |

#### Examples

##### Sample request:

```sh
curl --location --request POST 'localhost:8080/file/upload' \
	--form 'file=@"./5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp"' \
	--data-raw '{
		"password": "anypassword1234",
	}'
```

##### Sample response:

```json
{
    "message": "file uploaded",
    "filename": "5cf42bdf-aa14-4b33-8534-ea214fbd1c8f.pgp"
}
```
