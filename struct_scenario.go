package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Scenario struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Project     string             `json:"project"`     // 프로젝트명
	Order       float64            `json:"order"`       // 정렬순서
	Seq         int                `json:"seq"`         // Sequence
	Scene       int                `json:"scene"`       // Scene
	Cut         int                `json:"cut"`         // Cut
	Name        string             `json:"name"`        // 간단한 장면 설명
	IsPreviz    bool               `json:"ispreviz"`    // 프리비즈 존재여부
	IsTechviz   bool               `json:"istechviz"`   // 테크비즈 존재여부
	Thumbnails  []string           `json:"Thumbnails"`  // Thumbnails, AI로 이미지를 생성하기 때문에 여러개의 이미지중에 골라야 할 수 있다.
	SeedImage   string             `json:"seedimage"`   // Seed Image, AI를 이용해서 Image2Image로 이미지를 만들 때 사용하는 Seed 이미지
	Prompt      string             `json:"prompt"`      // 그림을 그릴때 사용되는 Prompt
	Script      string             `json:"script"`      // 스크립트
	VFXScript   string             `json:"vfxscript"`   // VFX 스크립트
	Time        string             `json:"time"`        // D,L,SS(Sunset)
	Location    string             `json:"location"`    // I,E
	Length      string             `json:"length"`      // L,S
	VFXSolution string             `json:"vfxsolution"` // VFX 솔루션
	Type        string             `json:"type"`        // 2D,3D
	Difficult   string             `json:"difficult"`   // 난이도
	EA          int                `json:"ea"`          // 갯수, 견적에 필요
	Manday      map[string]Task    `json:"manday"`      // Task 리스트
	Cost        int                `json:"cost"`        // 견적
}
