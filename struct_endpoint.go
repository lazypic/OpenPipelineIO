package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Endpoint struct {
	ID            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Endpoint      string             `json:"endpoint"`      // 엔드포인트 주소 예) https://api.dns.com/api/v1/users
	Method        []string           `json:"method"`        // GET,PUT,DELETE,POST,OPTIONS
	Parameter     []string           `json:"parameter"`     // Endpoint 에서 지원하는 옵션
	CORS          string             `json:"cors"`          // Cross-Origin Resource Sharing: origin || *
	Authorization string             `json:"authorization"` // 인증방식
	ContentType   string             `json:"contenttype"`   // Content-Type
	Model         string             `json:"model"`         // Model 예) AI모델
	Description   string             `json:"description"`   // 설명
	IsWebhook     bool               `json:"iswebhook"`     // WebHook 서비스 인가?
	IsUser        bool               `json:"isuser"`        // 사용자 공개
	IsDeveloper   bool               `json:"isdeveloper"`   // 개발자 공개
	IsAdmin       bool               `json:"isadmin"`       // 관리자 공개
	IsSecurity    bool               `json:"issecurity"`    // 보안이 필요한가
	IsAsset       bool               `json:"isasset"`       // 회사의 자산이 될 수 있는 Endpoint 인가?
	IsPatent      bool               `json:"ispatent"`      // 특허와 관련된 기술이 처리되는 Endpoint 인가?
	IsUpload      bool               `json:"isupload"`      // 업로드가 되는 Endpoint인가?
	Token         string             `json:"token"`         // 보안토큰
	Tags          []string           `json:"tags"`          // 태그
	Apikey        string             `json:"apikey"`        // API KEY
	Curl          string             `json:"curl"`          // curl 예시
	Category      string             `json:"category"`      // 카테고리: 유저,결제,관리,CS(고객대응),메일,알림,메시징,내저징데이터,정산,검색,통계,마이페이지,공통사항,인증,웹소캣,공지사항,시스템 관리,메인페이지,모바일지원,광고,혜택,이벤트,친구관리,찜,장바구니,공유,채팅,기기관리,통화,음성,카메라,구매,보안,선물,뉴스,모니터링,구독,데이터분석,머신러닝,번역
}
