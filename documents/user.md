# 사용자

사용자 관련 Commandline 명령어

## 최초 관리자 설정

최초 사용자를 가입 이후 관리자 권한의 Shell로 해당 아이디를 관리자로 만들어 줍니다.

```bash
$ sudo openpipelineio -accesslevel 11 -id [userid]
```

## 서비스 시작시 사용자 권한 설정

openpipelineio를 실행하면 기본적으로 accesslevel 0 으로 가입됩니다.
만약 웹서비스 시작시 accesslevel 3으로 서비스를 시작하고 싶다면, -initaccesslevel 옵션을 이용해서 웹서버를 실행할 수 있습니다.

```bash
$ sudo openpipelineio -initaccesslevel 3 -http :80
```