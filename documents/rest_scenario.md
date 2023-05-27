# RestAPI Scenario

Scenario Restapi 입니다.

## GET

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/scenario/{id} | 시나리오 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/scenario/{id}"
| /api/ganimages/{id} | 시나리오의 GANImage 정보를 가져옵니다|id|curl -X GET -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/ganimages/{id}"

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/scenario | 새로운 시나리오 정보를 추가합니다 | Option 참고 | curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"project":"test","script":"씬 내용"}' "https://openpipeline.io/api/scenario"
| /api/scenarios | 새로운 시나리오 리스트 정보를 추가합니다 | Option 참고 | curl -X POST -H 'Authorization: Basic {TOKEN}' -d '[{"project":"test","script":"씬 내용"}]' "https://openpipeline.io/api/scenarios"
| /api/ganimage/{id} | GANImage를 추가합니다 | Option 참고 | curl -X POST -H 'Authorization: Basic {TOKEN}' -d '{"url":"https://server/test.jpg","prompt":"프롬프트"}' "https://openpipeline.io/api/ganimage/{id}"

- Option: https://github.com/lazypic/OpenPipelineIO/blob/master/struct_scenario.go

## PUT

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/scenario/{id} | 기존 시나리오 정보를 수정합니다 | Option 참고 |curl -X PUT -H "Authorization: Basic {TOKEN}" -d '{"project":"test","script":"씬 내용"}' "https://openpipeline.io/api/scenario/{id}"
| /api/ganimage/{id} | GANImage를 수정합니다(url기준) | Option 참고 | curl -X PUT -H 'Authorization: Basic {TOKEN}' -d '{"url":"https://server/test.jpg","prompt":"프롬프트"}' "https://openpipeline.io/api/ganimage/{id}"

## DELETE

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/scenario/{id} | 값을 삭제합니다.|id|curl -X DELETE -H "Authorization: Basic {TOKEN}" "https://openpipeline.io/api/scenario/{id}"
| /api/ganimage/{id} | GANImage를 삭제합니다(url기준) | Option 참고 | curl -X DELETE -H 'Authorization: Basic {TOKEN}' -d '{"url":"https://server/test.jpg"}' "https://openpipeline.io/api/ganimage/{id}"

## Option 체크

```bash
curl https://openpipeline.io/api/scenario -v
```

```bash
HTTP/1.1 200 OK
< Access-Control-Allow-Methods: GET,PUT,DELETE,OPTIONS,POST
< Access-Control-Allow-Origin: *
< Date: Tue, 17 May 2022 02:10:41 GMT
< Content-Length: 0
```
