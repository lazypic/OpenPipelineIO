package main

// Task 자료구조는 태크스 정보를 담는 자료구조이다.
type Task struct {
	Title        string               `json:"title"`        // 테스크 제목
	UserID       string               `json:"userid"`       // 아티스트ID
	User         string               `json:"user"`         // 아티스트명
	UserComment  string               `json:"usercomment"`  // 아티스트 코멘트
	StatusV2     string               `json:"statusv2"`     // 샷 상태.
	ReviewStatus string               `json:"reviewstatus"` // 리뷰상태
	BeforeStatus string               `json:"beforestatus"` // 이전상태
	Startdate    string               `json:"startdate"`    // 1차 시작일 RFC3339 <- 최초시작일(툴이 커지면서 이 부분은 지져분해졌다.)
	Predate      string               `json:"predate"`      // 1차 마감일 RFC3339
	Startdate2nd string               `json:"startdate2nd"` // 2차 시작일 RFC3339
	Date         string               `json:"date"`         // 2차 마감일 RFC3339
	Mov          string               `json:"mov"`          // mov 경로
	Mdate        string               `json:"mdate"`        // mov 업데이트 날짜 RFC3339
	ExpectDay    int                  `json:"expectday"`    // 예측 맨데이
	ResultDay    int                  `json:"resultday"`    // 실제 맨데이
	UserNote     string               `json:"usernote"`     // 아티스트와 관련된 엘리먼트등의 정보를 입력하기 위해 사용.
	TaskLevel    `json:"tasklevel"`   // 샷 레벨
	Publishes    map[string][]Publish // 퍼블리쉬 정보, string값은 "Primary Key"가 된다.
}

// TaskLevel 은 태스크 난이도이다.
type TaskLevel int

// TaskLevel list
const (
	TaskLevel0 = TaskLevel(iota) // 0 쉬운난이도
	TaskLevel1                   // 1
	TaskLevel2                   // 2
	TaskLevel3                   // 3
	TaskLevel4                   // 4
	TaskLevel5                   // 5 높은난이도
)
