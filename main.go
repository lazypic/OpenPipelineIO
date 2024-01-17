package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"gopkg.in/mgo.v2"
)

func init() {
	// unidoc 라이센스키 발급: https://cloud.unidoc.io
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if err != nil {
		fmt.Println("not load unidoc module")
	}
}

var (
	// DBIP 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	DBIP = "127.0.0.1"
	// DBPORT mongoDB 기본포트.
	DBPORT = ":27017"
	// DBNAME 값은 데이터베이스 이름이다.
	DBNAME = "OpenPipelineIO"
	// DNS 값은 서비스 DNS 값입니다.
	DNS = "openpipeline.io"
	// TEMPLATES 값은 웹서버 실행전 사용할 템플릿이다.
	TEMPLATES = template.New("")
	// SHA1VER  은 Git SHA1 값이다.
	SHA1VER = "26b300a004abae553650c924514dc550e7385c9e" // 첫번째 커밋
	// BUILDTIME 은 빌드타임 시간이다.
	BUILDTIME = "2012-11-08T10:00:00" // 최초로 만든 시간
	// CachedAdminSetting 은 서비스 시작전 어드민 셋팅값을 메모리에 넣어서 사용되는 변수이다.
	CachedAdminSetting = Setting{}

	// 주요서비스 인수
	flagDBIP           = flag.String("dbip", DBIP+DBPORT, "mongodb ip and port")                                                            // mgo용 mongoDB 주소
	flagMongoDBURI     = flag.String("mongodburi", fmt.Sprintf("mongodb://%s%s", DBIP, DBPORT), "mongoDB URI ex)mongodb://localhost:27017") //mongo-driver용 인수
	flagDBName         = flag.String("dbname", DBNAME, "mongodb db name")                                                                   // mongoDB DB이름
	flagDebug          = flag.Bool("debug", false, "디버그모드 활성화")
	flagHTTPPort       = flag.String("http", "", "Web Service Port number.")          // 웹서버 포트
	flagVersion        = flag.Bool("version", false, "Print Version")                 // 버전
	flagCookieAge      = flag.Int64("cookieage", 168, "cookie age (hour)")            // 기본 일주일(168시간)로 설정한다. 참고: MPAA 기준 4시간이다.
	flagThumbnailAge   = flag.Int("thumbnailage", 1, "thumbnail image age (seconds)") // 썸네일 업데이트 시간. 3600초 == 1시간
	flagAuthmode       = flag.Bool("authmode", false, "restAPI authorization active") // restAPI 이용시 authorization 활성화
	flagCertFullchanin = flag.String("certfullchanin", fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", DNS), "certification fullchain path")
	flagCertPrivkey    = flag.String("certprivkey", fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", DNS), "certification privkey path")
	// Process
	flagProcessBufferSize = flag.Int("processbuffersize", 100, "process buffer size") // 최대 대기 리스트
	flagMaxProcessNum     = flag.Int("maxprocessnum", 4, "max process number")        // 최대 연산 갯수
	flagProcessDuration   = flag.Int64("processduration", 10, "process duration")
	flagReviewRender      = flag.Bool("reviewrender", false, "리뷰 렌더링을 허용하는 옵션")
	flagScanPlateRender   = flag.Bool("scanplaterender", false, "ScanPlate 렌더링을 허용하는 옵션")

	// Commandline Args
	flagHelp              = flag.Bool("help", false, "자세한 도움말을 봅니다.")
	flagID                = flag.String("id", "", "user id")
	flagAccessLevel       = flag.Int("accesslevel", -1, "edit user Access Level")
	flagSignUpAccessLevel = flag.Int("signupaccesslevel", 0, "signup access level")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("OpenPipelineIO: ")
	flag.Usage = usage
	flag.Parse()
	if *flagVersion {
		fmt.Println("buildTime:", BUILDTIME)
		fmt.Println("git SHA1:", SHA1VER)
		os.Exit(0)
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	ip, err := serviceIP()
	if err != nil {
		log.Fatal(err)
	}
	if *flagAccessLevel != -1 && *flagID != "" {
		if user.Username != "root" {
			log.Fatal(errors.New("사용자의 레벨을 수정하기 위해서는 root 권한이 필요합니다"))
		}
		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		u, err := getUser(session, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		err = rmToken(session, u.ID)
		if err != nil {
			log.Fatal(err)
		}
		u.AccessLevel = AccessLevel(*flagAccessLevel)
		err = setUser(session, u)
		if err != nil {
			log.Fatal(err)
		}
		err = addToken(session, u)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagHTTPPort != "" {
		// 만약 프로젝트가 하나도 없다면 "TEMP" 프로젝트를 생성한다. 프로젝트가 있어야 템플릿이 작동하기 때문이다.
		session, err := mgo.DialWithTimeout(*flagDBIP, 2*time.Second)
		if err != nil {
			log.Fatal("DB가 실행되고 있지 않습니다.")
		}
		admin, err := GetAdminSetting(session) // V2 로 변경하기
		if err != nil {
			log.Fatal(err)
		}
		// 어드민설정을 한번 저장한다. CachedAdminSetting값은 매번 DB를 호출하면 안되는 작업에서 사용된다.
		CachedAdminSetting = admin
		// 만약 Admin설정에 ThumbnailRootPath가 잡혀있다면 그 값을 이용한다.
		if admin.ThumbnailRootPath == "" {
			log.Println("admin 설정창의 thumbnail 경로지정이 필요합니다.")
		}

		plist, err := Projectlist(session)
		if err != nil {
			log.Fatal(err)
		}
		// 프로젝트가 존재하지 않는다면 test 프로젝트를 추가한다.
		if len(plist) == 0 {
			p := *NewProject("test")
			err = addProject(session, p)
			if err != nil {
				log.Fatal(err)
			}
		}
		session.Close()
		if *flagHTTPPort == ":80" {
			fmt.Printf("Service start: http://%s\n", ip)
		} else if *flagHTTPPort == ":443" {
			fmt.Printf("Service start: https://%s\n", ip)
		} else {
			fmt.Printf("Service start: http://%s%s\n", ip, *flagHTTPPort)
		}
		vfsTempates, err := LoadTemplates()
		if err != nil {
			log.Fatal(err)
		}
		TEMPLATES = vfsTempates
		if *flagReviewRender {
			go ProcessReviewRender() // 연산(Review데이터 등등)이 필요한 것들이 있다면 연산을 시작한다.
		}
		if *flagScanPlateRender {
			go ProcessScanPlateRender() // 연산(Review데이터 등등)이 필요한 것들이 있다면 연산을 시작한다.
		}
		webserver(*flagHTTPPort)
	}
	if *flagHelp {
		flag.Usage()
		return
	}
	flag.Usage()
}
