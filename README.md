# OpenPipelineIO
![example workflow](https://github.com/lazypic/OpenPipelineIO/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/lazypic/OpenPipelineIO)](https://goreportcard.com/report/github.com/lazypic/OpenPipelineIO)

![screenshot](figures/screenshot.png)

![review](figures/review.png)

![statistics](figures/statistics.png)

OpenPipelineIO는 영화, 드라마, 전시영상, 애니메이션, 게임 등 콘텐츠 제작을 위한 프로젝트 매니징 솔루션, DATA IO 솔루션 입니다.
업무에 대해 파이프라인을 구분짓고 설정할 수 있다는 것은 한 조직에서 초기 설정된 정보를 활용하여, 프로젝트 진행, 부분 자동화, 완전 자동화, 빅데이터 단계, AI 단계를 준비할 수 있음을 의미합니다.

- 속도, 검색어 방식, 교육의 최소화, 단일파일 배포를 중점으로 개발되고 있습니다.
- 내부, 외부 서버에 설치가 가능합니다.
- 리뷰 시스템
- 사용자별 토큰키, 암호화키, 직급별 접근권한 사용이 가능합니다.
- [Google Site](https://sites.google.com/view/lazypic/openpipelineio)
- [Collaborate with other open sources](https://landscape.aswf.io/?category=aswf-member-company&grouping=category&fullscreen=yes)

## 설치 및 실행

데이터베이스, 파이프라인툴, 방화벽 순서대로 설정합니다.

### 데이터베이스(mongoDB) 설치 및 서비스 실행

- [RockyLinux, CentOS 에서 mongoDB 설정](https://github.com/lazypic/tdcourse/blob/master/docs/install_mongodb.md)
- [데비안 리눅스에서 설치하기](documents/install_debian.md)
- [macOS에서 설치하기](documents/install_macOS.md)
- [freeBSD에서 설치하기](documents/install_freebsd.md)
- [AWS EC2에 설치하기](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-on-amazon/)


### OpenPipelineIO 설치 및 실행

https://github.com/lazypic/OpenPipelineIO/releases 에서 최신 버전을 다운받아 압축을 풀어주세요.

#### FreeBSD 다운로드 예시

```bash
wget https://github.com/lazypic/OpenPipelineIO/releases/download/v3.13.19/OpenPipelineIO_freebsd_amd64.tgz
tar -xzvf ./OpenPipelineIO_freebsd_amd64.tgz
```

#### 실행

```bash
OpenPipelineIO -http :80 # 웹서버를 실행합니다.
OpenPipelineIO -http :80 -reviewrender # 웹서버 및 FFmpeg를 이용하여 리뷰를 렌더링하는 서버
```

1. 최초에 하단 Sign-up을 눌러서 관리자로 가입합니다.
1. Ctrl + C를 눌러서 웹서비스를 종료합니다.
1. 관리자 권한(sudo, su)을 이용해서 가입된 ID에 admin 권한을 줍니다.

```bash
./OpenPipelineIO -accesslevel 11 -id `최초가입ID명`
```

1. OpenPipelineIO 서비스를 다시 실행합니다.
1. 웹주소에 /adminsetting 값을 붙혀서 최초 admin 설정을 해줍니다.
1. 서비스를 다시 재시작해줍니다.

```bash
nohup OpenPipelineIO -http :80 &
```

> 여러분이 macOS 또는 리눅스에서 기본 웹서버가 켜진 사용한다면 기본적으로 80포트는 웹서버가 사용중일 수 있습니다. 80포트에 실행되는 아파치 서버를 종료하기 위해서 `$ sudo apachectl stop` 를 터미널에 입력해주세요.

### 방화벽 설정

다른 컴퓨터에서 접근하기 위해서는 해당 포트를 방화벽 해제합니다.

```bash
sudo firewall-cmd --zone=public --add-port=80/tcp --permanent
success
sudo firewall-cmd --reload
```

### 기타툴의 연동

OpenPipelineIO는 [wfs-웹파일시스템](https://github.com/digital-idea/wfs), [웹프로토콜](https://github.com/lazypic/opio)과 같이 연동됩니다. 아래 서비스 실행 및 프로토콜 설치도 같이 진행하면 더욱 편리한 OpenPipelineIO를 활용할 수 있습니다.

```bash
wfs -http :8081
```

### CommandLine

터미널에서 간단하게 명령어를 통해 관리를 할 수 있습니다.

- [User](documents/user.md)

### RestAPI

OpenPipelineIO는 RestAPI가 설계되어 있습니다.
Python, Go, Java, Javascript, node.JS, C++, C, C# 등 수많은 언어를 활용하여 OpenPipelineIO를 제어할 수 있습니다.

- [Project](documents/rest_project.md)
- [Item](documents/rest_item.md): Asset, Shot
- [User](documents/rest_user.md)
- [Organization](documents/rest_organization.md)
- [Tasksetting](documents/rest_tasksetting.md)
- [Status](documents/rest_status.md)
- [Review](documents/rest_review.md)
- [Statistics](documents/rest_statistics.md)
- [Partner](documents/rest_partner.md)
- [ProjectForPartner](documents/rest_projectforpartner.md)
- [Money](documents/rest_money.md)
- [Moneytype](documents/rest_moneytype.md)
- [Step](documents/rest_step.md)
- [Pipelinestep](documents/rest_pipelinestep.md)
- [FullCalendar Event](documents/rest_fcevent.md)
- [FullCalendar Resource](documents/rest_fcresource.md)
- [PDF](documents/rest_pdf.md)
- [Endpoint](documents/rest_endpoint.md)

### 썸네일 경로

위에서 생성된 thumbnail 폴더는 아래 구조를 띄고 있습니다.
썸네일은 사내 다른 응용프로그램에서도 사용될 수 있기 때문에 경로구조를 표기해둡니다.

- 썸네일주소 : `thumbnail/{projectname}/{id}.jpg`
- 사용자이미지 : `thumbnail/user/{id}.jpg`

### 프로젝트 Process

- [디자인 프로세스](documents/process_designer.md)
- [개발 프로세스](documents/process_developer.md)
- [Onset Setellite](documents/setellite.md)
- [DB관리](documents/dbbackup.md)

### Developer

- OpenPipelineIO: <https://openpipeline.io>
- Log서버: <https://openpipeline.io:8080>
- WFS서버: <https://openpipeline.io:8081>
- 회사 전용 빌드문의: hello@lazypic.org
- Maintainer: Jason / jason@lazypic.org
- Committer: Alex / alex@lazypic.org
- Contributors:
- 체험계정 ID/PW: guest
  - Guest 계정은 모든 메뉴가 보이지 않습니다.
  - Guest 계정은 일부 기능만 테스트 가능한 모드입니다.
  - 만약 많은 기능을 테스트하고 싶다면 가입한 ID와 함께 권한변경 요청메일을 hello@lazypic.org로 보내주세요.

### Infomation

- [OpenPipelineIO의 역사](documents/history.md)
- License: BSD 3-Clause License


### Support Companys

- DigitalIdea
- Magnon
- MMHUB
- D1TUS
- 75mm-Studio

### License

- OpenPipelineIO: BSD 3-Clause License
- [JScolor](http://jscolor.com/download/): GNU GPL license v3
- [Dropzone](https://www.dropzonejs.com): MIT License
- [JQuery](https://jquery.org/license/): MIT license
- [VFS](https://github.com/blang/vfs): MIT license
- [HttpFS](https://github.com/shurcooL/httpfs): MIT license
- [VFSgen](https://github.com/shurcooL/vfsgen): MIT license
- [Excelize](https://github.com/360EntSecGroup-Skylar/excelize): BSD 3-Clause License
- [Slack go webhook](https://github.com/ashwanthkumar/slack-go-webhook): Apache License, Version 2.0
- [Captcha](https://github.com/dchest/captcha): Apache License, Version 2.0
- [Mgo](https://github.com/go-mgo/mgo): <https://github.com/go-mgo/mgo/blob/v2-unstable/LICENSE>
- [JWT go](https://github.com/golang-jwt/jwt): MIT license
- [OpenColorIO](https://github.com/AcademySoftwareFoundation/OpenColorIO): BSD 3-Clause License
- [alfg/mp4](https://github.com/alfg/mp4): MIT license
- [amarburg/go-quicktime](https://github.com/amarburg/go-quicktime): MIT license
- [Gollia Mux](https://github.com/gorilla/mux): BSD 3-Clause License
