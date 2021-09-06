package model

type ShortLink struct {
	Id        string
	ShortUrl  string
	TargetUrl string
	Cnt       int32 `json:"Hits"`
	Expire    int64
}

type Blacklist struct {
	IsRegex bool
	Url     string
}
