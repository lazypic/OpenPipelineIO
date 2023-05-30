package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/digital-idea/dilog"
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
	// APPNAME 값은 어플리케이션 이름이다.
	APPNAME = "OpenPipelineIO"
	// DNS 값은 서비스 DNS 값입니다.
	DNS = "openpipeline.io"
	// MAILDNS 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	MAILDNS = "lazypic.org"
	// COMPANY 값은 컴파일 단계에서 회사에 따라 값이 바뀐다.
	COMPANY = "Lazypic"
	// TEMPLATES 값은 웹서버 실행전 사용할 템플릿이다.
	TEMPLATES = template.New("")
	// SHA1VER  은 Git SHA1 값이다.
	SHA1VER = "26b300a004abae553650c924514dc550e7385c9e" // 첫번째 커밋
	// BUILDTIME 은 빌드타임 시간이다.
	BUILDTIME = "2012-11-08T10:00:00" // 최초로 만든 시간
	// CachedAdminSetting 은 서비스 시작전 어드민 셋팅값을 메모리에 넣어서 사용되는 변수이다.
	CachedAdminSetting = Setting{}

	// 주요서비스 인수
	flagDBIP       = flag.String("dbip", DBIP+DBPORT, "mongodb ip and port")                                                            // mgo용 mongoDB 주소
	flagMongoDBURI = flag.String("mongodburi", fmt.Sprintf("mongodb://%s%s", DBIP, DBPORT), "mongoDB URI ex)mongodb://localhost:27017") //mongo-driver용 인수
	flagDBName     = flag.String("dbname", DBNAME, "mongodb db name")                                                                   // mongoDB DB이름
	flagAppName    = flag.String("appname", APPNAME, "app name")
	flagMailDNS    = flag.String("maildns", MAILDNS, "mail DNS name")

	flagDebug          = flag.Bool("debug", false, "디버그모드 활성화")
	flagHTTPPort       = flag.String("http", "", "Web Service Port number.")          // 웹서버 포트
	flagCompany        = flag.String("company", COMPANY, "Web Service Port number.")  // 회사이름
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
	flagAdd                = flag.String("add", "", "add project, add item(shot, asset)")
	flagRm                 = flag.String("rm", "", "remove project, shot, asset, user")
	flagProject            = flag.String("project", "", "project name")
	flagName               = flag.String("name", "", "name")
	flagSeason             = flag.String("season", "", "season")
	flagEpisode            = flag.String("episode", "", "episode")
	flagNetflixID          = flag.String("netflixid", "", "netflix id")
	flagType               = flag.String("type", "", "type: org,left,asset,org1,src,src1,lsrc,rsrc")
	flagAssettags          = flag.String("assettags", "", "asset tags, 입력예) prop,char,env,prop,comp,plant,vehicle,global,component,group,assembly 형태로 입력")
	flagAssettype          = flag.String("assettype", "", "assettype: char,env,global,prop,comp,plant,vehicle,global,group") // 추후 삭제예정.
	flagHelp               = flag.Bool("help", false, "자세한 도움말을 봅니다.")
	flagThumbnailImagePath = flag.String("thumbnailimagepath", "", "Thumbnail image 경로")
	flagThumbnailMovPath   = flag.String("thumbnailmovpath", "", "Thumbnail mov 경로")
	flagPlatePath          = flag.String("platepath", "", "Plate 경로")
	// Commandline Args: User
	flagID                = flag.String("id", "", "user id")
	flagAccessLevel       = flag.Int("accesslevel", -1, "edit user Access Level")
	flagSignUpAccessLevel = flag.Int("signupaccesslevel", 3, "signup access level")
	// scan정보 추가. plate scan tool에서 데이터를 등록할 때 활용되는 옵션
	flagPlatesize       = flag.String("platesize", "", "스캔 플레이트 사이즈")
	flagScantimecodein  = flag.String("scantimecodein", "00:00:00:00", "스캔 Timecode In")
	flagScantimecodeout = flag.String("scantimecodeout", "00:00:00:00", "스캔 Timecode Out")
	flagJusttimecodein  = flag.String("justtimecodein", "00:00:00:00", "Just구간 Timecode In")
	flagJusttimecodeout = flag.String("justtimecodeout", "00:00:00:00", "Just구간 Timecode Out")
	flagScanframe       = flag.Int("scanframe", 0, "스캔 총 프레임수")
	flagScanname        = flag.String("scanname", "", "스캔 폴더명")
	flagScanin          = flag.Int("scanin", -1, "스캔 In Frame")
	flagScanout         = flag.Int("scanout", -1, "스캔 Out Frame")
	flagJustin          = flag.Int("justin", -1, "Just In Frame")
	flagJustout         = flag.Int("justout", -1, "Just Out Frame")
	flagPlatein         = flag.Int("platein", -1, "플레이트 In Frame")
	flagPlateout        = flag.Int("plateout", -1, "플레이트 Out Frame")
	flagUpdateParent    = flag.Bool("updateparent", false, "org1,org2 형태의 재스캔 항목이라면 원본 org 정보를 업데이트 한다.")
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
	} else if *flagRm == "user" && *flagID != "" { // 사용자 삭제
		if user.Username != "root" {
			log.Fatal(errors.New("사용자를 삭제하기 위해서는 root 권한이 필요합니다"))
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
		err = rmUser(session, u.ID)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagAdd == "item" && *flagName != "" && *flagProject != "" && *flagType != "" { //아이템 추가
		switch *flagType {
		case "org", "left": // 일반영상은 org가 샷 타입이다. 입체프로젝트는 left가 샷타입이다.
			addShotItemCmd(*flagProject, *flagName, *flagType, *flagPlatesize, *flagScanname, *flagScantimecodein, *flagScantimecodeout, *flagJusttimecodein, *flagJusttimecodeout, *flagScanframe, *flagScanin, *flagScanout, *flagPlatein, *flagPlateout, *flagJustin, *flagJustout)
			dilog.Add(*flagDBIP, ip, "샷 생성되었습니다.", *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			dilog.Add(*flagDBIP, ip, "스캔이름 : "+*flagScanname, *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("스캔타임코드 : %s(%d) / %s(%d) (총%df)", *flagScantimecodein, *flagScanin, *flagScantimecodeout, *flagScanout, *flagScanframe), *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("플레이트 구간 : %d - %d", *flagPlatein, *flagPlateout), *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			dilog.Add(*flagDBIP, ip, "플레이트 사이즈 : "+*flagPlatesize, *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			return
		case "asset": //에셋 추가
			addAssetItemCmd(*flagProject, *flagName, *flagType, *flagAssettype, *flagAssettags)
			dilog.Add(*flagDBIP, ip, "에셋이 생성되었습니다.", *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			dilog.Add(*flagDBIP, ip, fmt.Sprintf("에셋타입 : %s, 에셋태그 : %s", *flagAssettype, *flagAssettags), *flagProject, *flagName+"_"+*flagType, *flagAppName, user.Username, 180)
			return
		default: //소스, 재스캔 추가
			addOtherItemCmd(*flagProject, *flagName, *flagType, *flagPlatesize, *flagScanname, *flagScantimecodein, *flagScantimecodeout, *flagJusttimecodein, *flagJusttimecodeout, *flagScanframe, *flagScanin, *flagScanout, *flagPlatein, *flagPlateout, *flagJustin, *flagJustout)
			logString := fmt.Sprintf("Create Item: %s_%s, Scanname: %s, ScanTimecode: %s(%d) / %s(%d) (Total:%df), Plate Range: %d - %d, Platesize: %s",
				*flagName,
				*flagType,
				*flagScanname,
				*flagScantimecodein, *flagScanin, *flagScantimecodeout, *flagScanout, *flagScanframe,
				*flagPlatein, *flagPlateout,
				*flagPlatesize,
			)
			dilog.Add(*flagDBIP, ip, logString, *flagProject, *flagName, *flagAppName, user.Username, 180)
			if *flagUpdateParent {
				// updateParent 옵션이 활성화되어있고, org, left가 재스캔이라면.. 원본플레이트의 정보를 업데이트한다.
				if (*flagType != "org" && strings.Contains(*flagType, "org")) || (*flagType != "left" && strings.Contains(*flagType, "left")) {
					session, err := mgo.Dial(*flagDBIP)
					if err != nil {
						log.Fatal(err)
					}
					defer session.Close()
					typ := "org"
					if strings.Contains(*flagType, "left") {
						typ = "left"
					}
					item, err := getItem(session, *flagProject, *flagName+"_"+typ)
					if err != nil {
						log.Fatal(err)
					}
					item.Platesize = *flagPlatesize
					item.ScanTimecodeIn = *flagScantimecodein
					item.ScanTimecodeOut = *flagScantimecodeout
					item.JustTimecodeIn = *flagJusttimecodein
					item.JustTimecodeOut = *flagJusttimecodeout
					item.ScanIn = *flagScanin
					item.ScanOut = *flagScanout
					item.ScanFrame = *flagScanframe
					item.Scanname = *flagScanname
					item.PlateIn = *flagPlatein
					item.PlateOut = *flagPlateout
					item.JustIn = *flagJustin
					item.JustOut = *flagJustout
					item.UseType = *flagType
					// adminsetting 값을 가지고와서 Thumbnailmov 값을 설정한다.
					admin, err := GetAdminSetting(session)
					if err != nil {
						log.Fatal(err)
					}
					var thumbnailMovPath bytes.Buffer
					thumbnailMovPathTmpl, err := template.New("thumbnailMovPath").Parse(admin.ThumbnailMovPath)
					if err != nil {
						log.Fatal(err)
					}
					err = thumbnailMovPathTmpl.Execute(&thumbnailMovPath, item)
					if err != nil {
						log.Fatal(err)
					}
					item.Thummov = thumbnailMovPath.String()

					err = setItem(session, *flagProject, item)
					if err != nil {
						log.Fatal(err)
					}
					// log
					logString := fmt.Sprintf("Scanname: %s, ScanTimecode: %s(%d) / %s(%d) Total: %df\nPlate Range: %d - %d\nPlatesize: %s",
						*flagScanname,
						*flagScantimecodein,
						*flagScanin,
						*flagScantimecodeout,
						*flagScanout,
						*flagScanframe,
						*flagPlatein,
						*flagPlateout,
						*flagPlatesize,
					)
					dilog.Add(*flagDBIP, ip, logString, *flagProject, *flagName, *flagAppName, user.Username, 180)
				}
			}
			return
		}
	} else if *flagRm == "item" && *flagName != "" && *flagProject != "" && *flagType != "" { //아이템 삭제
		rmItemCmd(*flagProject, *flagName, *flagType)
		return
	} else if *flagHTTPPort != "" {
		// 만약 프로젝트가 하나도 없다면 "TEMP" 프로젝트를 생성한다. 프로젝트가 있어야 템플릿이 작동하기 때문이다.
		session, err := mgo.DialWithTimeout(*flagDBIP, 2*time.Second)
		if err != nil {
			log.Fatal("DB가 실행되고 있지 않습니다.")
		}
		admin, err := GetAdminSetting(session)
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
