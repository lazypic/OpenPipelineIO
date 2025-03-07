package main

// Setting 자료구조는 관리자 설정 자료구조이다.
type Setting struct {
	ID                             string `json:"id"`                             // 셋팅ID
	AppName                        string `json:"appname"`                        // App 이름
	EmailDNS                       string `json:"emaildns"`                       // Email DNS 서버이름
	OCIOConfig                     string `json:"ocioconfig"`                     // OpenColorIO Config Path 설정
	FFmpeg                         string `json:"ffmpeg"`                         // FFmpeg 경로 셋팅
	FFprobe                        string `json:"ffprobe"`                        // FFprobe 경로 셋팅
	FFmpegThreads                  int    `json:"ffmpegthreads"`                  // FFmpeg 연산 Thread 셋팅
	OpenImageIO                    string `json:"openimageio"`                    // OpenImageIO 경로
	Iinfo                          string `json:"iinfo"`                          // iinfo 경로
	Curl                           string `json:"curl"`                           // curl 경로
	RVPath                         string `json:"rvpath"`                         // RV 경로 셋팅
	RootPath                       string `json:"rootpath"`                       // Root경로 예) /show
	ProjectPath                    string `json:"projectpath"`                    // Project경로 예) /show/{{.Project}}
	ProjectPathPermission          string `json:"projectpathpermission"`          // Project경로의 권한
	ProjectPathUID                 string `json:"projectpathuid"`                 // Project경로의 User ID
	ProjectPathGID                 string `json:"projectpathgid"`                 // Project경로의 Group ID
	ScanPlateUploadPath            string `json:"scanplateuploadpath"`            // ScanPlate 업로드 경로
	ShotRootPath                   string `json:"shotrootpath"`                   // Shot Root 경로 예) /show/{{.Project}}/seq/
	ShotRootPathPermission         string `json:"shotrootpathpermission"`         // Shot Root 경로의 권한
	ShotRootPathUID                string `json:"shotrootpathuid"`                // Shot Root 경로의 User ID
	ShotRootPathGID                string `json:"shotrootpathgid"`                // Shot Root 경로의 Group ID
	SeqPath                        string `json:"seqpath"`                        // Seq 경로 예) /show/{{.Project}}/seq/{{.Seq}}
	SeqPathPermission              string `json:"seqpathpermission"`              // Seq 경로의 권한
	SeqPathUID                     string `json:"seqpathuid"`                     // Seq 경로의 User ID
	SeqPathGID                     string `json:"seqpathgid"`                     // Seq 경로의 Group ID
	ShotPath                       string `json:"shotpath"`                       // 개별 Shot 경로 예) /show/{{.Project}}/seq/{{.Seq}}/{{.Name}}
	ShotPathPermission             string `json:"shotpathpermission"`             // 개별 Shot 경로의 권한
	ShotPathUID                    string `json:"shotpathuid"`                    // 개별 Shot 경로의 User ID
	ShotPathGID                    string `json:"shotpathgid"`                    // 개별 Shot 경로의 Group ID
	AssetRootPath                  string `json:"assetrootpath"`                  // Asset Root 경로 예) /show/{{.Project}}/assets/
	AssetRootPathPermission        string `json:"assetrootpathpermission"`        // Asset Root 경로의 권한
	AssetRootPathUID               string `json:"assetrootpathuid"`               // Asset Root 경로의 User ID
	AssetRootPathGID               string `json:"assetrootpathgid"`               // Asset Root 경로의 Group ID
	AssetTypePath                  string `json:"assettypepath"`                  // Asset Type 경로 예) /show/{{.Project}}/assets/{{.Assettype}}
	AssetTypePathPermission        string `json:"assettypepathpermission"`        // Asset Type 경로의 권한
	AssetTypePathUID               string `json:"assettypepathuid"`               // Asset Type 경로의 User ID
	AssetTypePathGID               string `json:"assettypepathgid"`               // Asset Type 경로의 Group ID
	AssetPath                      string `json:"assetpath"`                      // 개별 Asset 예) /show/{{.Project}}/assets/{{.Assettype}}/{{.Name}}
	AssetPathPermission            string `json:"assetpathpermission"`            // 개별 Asset 경로의 권한
	AssetPathUID                   string `json:"assetpathuid"`                   // 개별 Asset 경로의 User ID
	AssetPathGID                   string `json:"assetpathgid"`                   // 개별 Asset 경로의 Group ID
	ThumbnailRootPath              string `json:"thumbnailrootpath"`              // 썸네일 이미지 Root 경로
	ThumbnailRootPathPermission    string `json:"thumbnailrootpathpermission"`    // 썸네일 이미지 Root 경로의 권한
	ThumbnailRootPathUID           string `json:"thumbnailrootpathuid"`           // 썸네일 이미지 Root 경로의 User ID
	ThumbnailRootPathGID           string `json:"thumbnailrootpathgid"`           // 썸네일 이미지 Root 경로의 Group ID
	ThumbnailImagePath             string `json:"thumbnailimagepath"`             // 썸네일 이미지 경로 /thumbnail/{{.Project}}/{{.Name}}_{{.Type}}
	ThumbnailImagePathPermission   string `json:"thumbnailimagepathpermission"`   // 썸네일 이미지 경로의 권한
	ThumbnailImagePathUID          string `json:"thumbnailimagepathuid"`          // 썸네일 이미지 경로의 User ID
	ThumbnailImagePathGID          string `json:"thumbnailimagepathgid"`          // 썸네일 이미지 경로의 Group ID
	ThumbnailMovPath               string `json:"thumbnailmovpath"`               // 썸네일 동영상 경로
	ThumbnailMovPathPermission     string `json:"thumbnailmovpathpermission"`     // 썸네일 동영상 경로의 권한
	ThumbnailMovPathUID            string `json:"thumbnailmovpathuid"`            // 썸네일 동영상 경로의 User ID
	ThumbnailMovPathGID            string `json:"thumbnailmovpathgid"`            // 썸네일 동영상 경로의 Group ID
	PlatePath                      string `json:"platepath"`                      // Plate 경로
	PlatePathPermission            string `json:"platepathpermission"`            // Plate 경로의 권한
	PlatePathUID                   string `json:"platepathuid"`                   // Plate 경로의 User ID
	PlatePathGID                   string `json:"platepathgid"`                   // Plate 경로의 Group ID
	ReviewDataPath                 string `json:"reviewdatapath"`                 // 리뷰 데이터가 저장되는 경로
	ReviewDataPathPermission       string `json:"reviewdatapathpermission"`       // 리뷰 데이터가 저장되는 경로의 권한
	ReviewDataPathUID              string `json:"reviewdatapathuid"`              // 리뷰 데이터가 저장되는 경로의 User ID
	ReviewDataPathGID              string `json:"reviewdatapathgid"`              // 리뷰 데이터가 저장되는 경로의 Group ID
	ReviewUploadPath               string `json:"reviewuploadpath"`               // 리뷰 업로드 파일이 저장되는 경로
	ReviewUploadPathPermission     string `json:"reviewuploadpathpermission"`     // 리뷰 업로드 파일이 저장되는 경로의 권한
	ReviewUploadPathUID            string `json:"reviewuploadpathuid"`            // 리뷰 업로드 파일이 저장되는 경로의 User ID
	ReviewUploadPathGID            string `json:"reviewuploadpathgid"`            // 리뷰 업로드 파일이 저장되는 경로의 Group ID
	DirectUploadPath               string `json:"directuploadpath"`               // Direct 업로드 파일이 저장되는 경로
	DirectUploadPathPermission     string `json:"directuploadpathpermission"`     // Direct 업로드 파일이 저장되는 경로의 권한
	DirectUploadPathUID            string `json:"directuploadpathuid"`            // Direct 업로드 파일이 저장되는 경로의 User ID
	DirectUploadPathGID            string `json:"directuploadpathgid"`            // Direct 업로드 파일이 저장되는 경로의 Group ID
	ProductionStartFrame           int    `json:"prodcutionstartframe"`           // 프로덕션의 시작프레임
	ProductionPaddingVersionNumber int    `json:"productionpaddingversionnumber"` // 프로덕션의 버전 자리수

	DefaultScaleRatioOfUndistortionPlate float64 `json:"defaultscaleratioofundistortionplate"` // 언디스토션 플레이트의 기본 스케일비율 1.1배
	ItemNumberOfPage                     int     `json:"itemnumberofpage"`                     // 한 페이지에 보이는 아이템 갯수
	MultipartFormBufferSize              int     `json:"multipartformbuffersize"`              // Multipart form buffer size
	ThumbnailImageWidth                  int     `json:"thumbnailimagewidth"`                  // Thumbnail Image 가로사이즈
	ThumbnailImageHeight                 int     `json:"thumbnailimageheight"`                 // Thumbnail Image 세로사이즈
	InitPassword                         string  `json:"initpassword"`                         // 사용자가 패스워드를 잃어버렸을 때 사용하는 패스워드
	EnableEndpoint                       bool    `json:"enableendpoint"`                       // Endpoint 활성화
	EnableDirectupload                   bool    `json:"enabledirectupload"`                   // Direct Upload 활성화
	EnableDirectuploadWithDate           bool    `json:"enabledirectuploadwithdate"`           // Enable Direct Upload with Date path.
	EnableDirectuploadWithProject        bool    `json:"enabledirectuploadwithproject"`        // Enable Direct Upload with Project path.
	EnableDirectuploadWithCompanyID      bool    `json:"enabledirectuploadwithcompanyid"`      // Enable Direct Upload with CompanyID path.
	FullcalendarSchedulerLicenseKey      string  `json:"fullcalendarschedulerlicensekey"`      // Fullcalendar License Key
	SlateFontPath                        string  `json:"slatefontpath"`                        // Slate에 사용되는 폰트 경로
	Protocol                             string  `json:"protocol"`                             // 프로토콜 이름
	WFS                                  string  `json:"wfs"`                                  // Web File system URL
	AudioCodec                           string  `json:"audiocodec"`                           // 오디오 코덱
}
