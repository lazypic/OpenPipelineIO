package main

// Tasksetting 자료구조이다
type Tasksetting struct {
	ID           string            `json:"id"`           // Task ID
	Name         string            `json:"name"`         // Task 표기명
	Type         string            `json:"type"`         // Type: Asset, Shot, R&D, Development
	LinuxPath    string            `json:"linuxpath"`    // Task 클릭시 프로토콜에서 열리는 리눅스 경로
	WindowPath   string            `json:"windowpath"`   // Task 클릭시 프로토콜에서 열리는 윈도우즈 경로
	MacOSPath    string            `json:"macospath"`    // Task 클릭시 프로토콜에서 열리는 맥 경로
	WFSPath      string            `json:"wfspath"`      // Task 클릭시 wfs에서 열리는 경로
	Attributes   map[string]string `json:"attributes"`   // Task에 필요한 속성추가. 예) 특정 Task는 멀티 퍼브리쉬 경로가 발생할 수 있다.
	Order        float64           `json:"order"`        // Task 순서. 드로잉시 정렬되는 순서이다. 낮을수록 위쪽
	ExcelOrder   float64           `json:"excelorder"`   // Excel 문서에서 Task 순서. 낮을수록 앞쪽
	InitGenerate bool              `json:"initgenerate"` // Item이 만들어질 때 해당 테스크 생성여부
	Pipelinestep string            `json:"pipelinestep"` // 파이프라인스탭 이름
}
