# RestAPI Step

PDF Restapi 입니다.

## POST

| URI | Description | Attributes | Curl Example |
| --- | --- | --- | --- |
| /api/pdf-to-json | pdf를 전송하고 json을 반환합니다. | pdf | curl -X POST  -H "Authorization: Basic {TOKEN}" -F "project=projectname" -F "version=20230403" -F "part=1" -F 'file=@/path/to/pdf/file.pdf' http://localhost/api/pdf-to-json"
