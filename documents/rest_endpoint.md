# RestAPI Endpoint

Endpoint Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/endpoint/{id}|엔드포인트 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/endpoint/{id}"
|/api/endpoints|모든 엔드포인트 정보를 가지고 온다. | . |curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/endpoints"
|/api/endpoints?search=word| 엔드포인트를 검색한다. | . |curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/endpoints?search=word"


## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/endpoint|새로운 엔드포인트 정보를 추가합니다| 자료구조 참고 |curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"endpoint":"https://api.dns.com/endpoint"}' "https://openpipeline.io/api/endpoint"

- Option: https://github.com/lazypic/OpenPipelineIO/blob/master/struct_endpoint.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/endpoint/{id} | 기존 엔드포인트 정보를 수정합니다 | 자료구조 참고 | curl -X PUT -H "Authorization: Basic {TOKEN}“ -d '{"endpoint":"https://api.dns.com/endpoint"}' "https://openpipeline.io/api/endpoint/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
|/api/endpoint/{id}|값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/endpoint/{id}"

## Option 체크

```bash
curl https://openpipeline.io/api/endpoint -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,PUT,DELETE,OPTIONS,POST
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
