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
	Start        string               `json:"start"`        // 시작일 RFC3339
	End          string               `json:"end"`          // 마감일 RFC3339
	Mov          string               `json:"mov"`          // mov 경로
	Mdate        string               `json:"mdate"`        // mov 업데이트 날짜 RFC3339
	UserNote     string               `json:"usernote"`     // 아티스트와 관련된 엘리먼트등의 정보를 입력하기 위해 사용.
	Publishes    map[string][]Publish // 퍼블리쉬 정보, string값은 "Primary Key"가 된다.
}
