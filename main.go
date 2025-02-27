package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// If you want unidoc license
	// Please set unidoc license use env.
	// info https://cloud.unidoc.io
	// unidoc is Go library for read pdf file
	err := license.SetMeteredKey(os.Getenv(`UNIDOC_LICENSE_API_KEY`))
	if *flagDebug {
		if err != nil {
			fmt.Println("not load unidoc module")
		}
	}
}

var (
	DBIP               = "127.0.0.1"
	DBPORT             = ":27017"
	DBNAME             = "OpenPipelineIO"
	DNS                = "openpipeline.io"
	TEMPLATES          = template.New("")                           // init template
	SHA1VER            = "26b300a004abae553650c924514dc550e7385c9e" // first git commit SHA1
	BUILDTIME          = "2012-11-08T10:00:00"                      // first commit time
	CachedAdminSetting = Setting{}                                  // init adminsetting for cache

	// 주요서비스 인수
	flagDBIP           = flag.String("dbip", DBIP+DBPORT, "mongodb ip and port")                                                            // mgo용 mongoDB 주소
	flagMongoDBURI     = flag.String("mongodburi", fmt.Sprintf("mongodb://%s%s", DBIP, DBPORT), "mongoDB URI ex)mongodb://localhost:27017") //mongo-driver용 인수
	flagDBName         = flag.String("dbname", DBNAME, "mongodb db name")                                                                   // mongoDB DB이름
	flagDebug          = flag.Bool("debug", false, "debug mode")
	flagHTTPPort       = flag.String("http", "", "Web Service Port number.")          // 웹서버 포트
	flagVersion        = flag.Bool("version", false, "Print Version")                 // 버전
	flagCookieAge      = flag.Int64("cookieage", 168, "cookie age (hour)")            // 기본 일주일(168시간)로 설정한다. 참고: MPAA 기준 Cookie save time is 4H.
	flagThumbnailAge   = flag.Int("thumbnailage", 1, "thumbnail image age (seconds)") // 썸네일 업데이트 시간. 3600초 == 1시간
	flagAuthmode       = flag.Bool("authmode", false, "restAPI authorization active") // restAPI 이용시 authorization 활성화
	flagCertFullchanin = flag.String("certfullchanin", fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", DNS), "certification fullchain path")
	flagCertPrivkey    = flag.String("certprivkey", fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", DNS), "certification privkey path")
	// Process
	flagProcessBufferSize = flag.Int("processbuffersize", 100, "process buffer size") // 최대 대기 리스트
	flagMaxProcessNum     = flag.Int("maxprocessnum", 4, "max process number")        // 최대 연산 갯수
	flagProcessDuration   = flag.Int64("processduration", 10, "process duration")
	flagReviewRender      = flag.Bool("reviewrender", false, "review render mode")
	flagScanPlateRender   = flag.Bool("scanplaterender", false, "scanplate render mode")

	// Commandline Args
	flagHelp              = flag.Bool("help", false, "help")
	flagID                = flag.String("id", "", "user id")
	flagAccessLevel       = flag.Int("accesslevel", -1, "edit user Access Level")
	flagSignUpAccessLevel = flag.Int("signupaccesslevel", 0, "signup access level")
)

func initMongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(*flagMongoDBURI))
	if err != nil {
		return nil, fmt.Errorf("failed to create new MongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("OPIO: ")
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
			log.Fatal(errors.New("need Root permission"))
		}
		client, err := initMongoClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(context.Background())
		u, err := getUserV2(client, *flagID)
		if err != nil {
			log.Fatal(err)
		}
		err = rmTokenV2(client, u.ID)
		if err != nil {
			log.Fatal(err)
		}
		u.AccessLevel = AccessLevel(*flagAccessLevel)
		err = setUserV2(client, u)
		if err != nil {
			log.Fatal(err)
		}
		err = addTokenV2(client, u)
		if err != nil {
			log.Fatal(err)
		}
		return
	} else if *flagHTTPPort != "" {
		client, err := initMongoClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(context.Background())

		// load admin setting from DB
		// and send to Cach
		admin, err := GetAdminSettingV2(client)
		if err != nil {
			log.Fatal(err)
		}
		CachedAdminSetting = admin
		if admin.ThumbnailRootPath == "" {
			log.Println("warning. need admin setup for thumbnail path")
		}

		plist, err := ProjectlistV2(client)
		if err != nil {
			log.Fatal(err)
		}
		// Create "test" project when no project in DB.
		if len(plist) == 0 {
			p := *NewProject("test")
			err = addProjectV2(client, p)
			if err != nil {
				log.Fatal(err)
			}
		}

		if *flagHTTPPort == ":80" {
			fmt.Printf("Service running: http://%s\n", ip)
		} else if *flagHTTPPort == ":443" {
			fmt.Printf("Service running: https://%s\n", ip)
		} else {
			fmt.Printf("Service running: http://%s%s\n", ip, *flagHTTPPort)
		}
		vfsTempates, err := LoadTemplates()
		if err != nil {
			log.Fatal(err)
		}
		TEMPLATES = vfsTempates
		if *flagScanPlateRender {
			// If scan data is available, start processing.
			go ProcessScanPlateRender()
		}
		if *flagReviewRender {
			// If review data is available, start processing.
			go ProcessReviewRender()
		}
		webserver(*flagHTTPPort)

	}
	if *flagHelp {
		flag.Usage()
		return
	}
}
