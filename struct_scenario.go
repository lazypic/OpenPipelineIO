package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Scenario struct {
	ID             primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Project        string             `json:"project"`        // 프로젝트명
	Order          float64            `json:"order"`          // 정렬순서
	Seq            int                `json:"seq"`            // Sequence
	Scene          int                `json:"scene"`          // Scene
	Cut            int                `json:"cut"`            // Cut
	Name           string             `json:"name"`           // 간단한 장면 설명
	IsPreviz       bool               `json:"ispreviz"`       // 프리비즈 존재여부
	IsTechviz      bool               `json:"istechviz"`      // 테크비즈 존재여부
	IsVisualLab    bool               `json:"isvisuallab"`    // 비쥬얼랩 개입여부
	GANImageIndex  int                `json:"ganimageindex"`  // 몇번째 이미지를 메인이미지로 사용할지에 대한 숫자. 기본값은 0 이다.
	GANImages      []GANImage         `json:"ganimages"`      // GANImage 리스트, AI로 이미지를 생성하기 때문에 여러개의 이미지중에 골라야 할 수 있다. /thumbpath/project/{id}/{seed}.png
	Prompt         string             `json:"prompt"`         // AI로 그림을 그릴때 사용되는 메인 Prompt
	NegativePrompt string             `json:"negativeprompt"` // AI로 그림을 그릴 때 적용되면 안되는 메인 NegativePrompt 정보
	Script         string             `json:"script"`         // 스크립트
	Time           string             `json:"time"`           // D,L,SS(Sunset)
	Location       string             `json:"location"`       // I,E
	Length         string             `json:"length"`         // L,S
	// 처음은 수기로, 나중엔 자동화로
	VFXScript   string          `json:"vfxscript"`   // VFX 스크립트
	VFXSolution string          `json:"vfxsolution"` // VFX 솔루션
	Type        string          `json:"type"`        // 2D,3D
	Difficult   string          `json:"difficult"`   // 난이도
	EA          int             `json:"ea"`          // 갯수, 견적에 필요
	Manday      map[string]Task `json:"manday"`      // Task 리스트
	Cost        int             `json:"cost"`        // 견적
	PageNum     int             `json:"pagenum"`     // 페이지수, 시나리오의 페이지수, 나중에 정보를 추적하기 좋다. PDFFormatScenario 자료구조에서 참고한다.
	LineNum     int             `json:"linenum"`     // 줄수, 줄수 정보는 나중에 정보를 역추적하기 좋다. PDFFormatScenario 자료구조에서 참고한다.
}

// GANImage 자료구조는 Generative Adversarial Network (GAN) image의 약자이다.
type GANImage struct {
	URL               string            `json:"url"`               // 생성이미지 URL
	SeedImage         string            `json:"seedimage"`         // Seed Image, AI를 이용해서 Image2Image로 이미지를 만들 때 사용하는 Seed 이미지
	SubPrompt         string            `json:"subprompt"`         // AI로 그림을 그릴때 사용되는 SubPrompt
	SubNegativePrompt string            `json:"subnegativeprompt"` // AI로 그림을 그릴 때 적용되면 안되는 SubNegativePrompt 정보
	Hyperparameter    map[string]string `json:"hyperparameter"`    // Hyperparamter 옵션
}

type PDFFormatScenario struct {
	Project string `json:"project"` // 프로젝트명
	Version string `json:"version"` // 문서 버전
	PageNum int    `json:"pagenum"` // 페이지수, 페이지 정보는 나중에 정보를 추적하기 좋다.
	LineNum int    `json:"linenum"` // 줄수, 줄수 정보는 나중에 정보를 역추적하기 좋다.
	Text    string `json:"text"`    // 글자
}
