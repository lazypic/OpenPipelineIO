# RestAPI Step

PDF Restapi 입니다.

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/pdf-to-json | pdf를 전송하고 json을 반환합니다. | pdf | curl -X POST -F 'pdf=@/path/to/pdf/file.pdf' -H "Authorization: Basic {TOKEN}" http://localhost/api/pdf-to-json
"

