package main

//基于游标的分页
import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Token string

// Page 分页
type Page struct {
	NextID        string `json:"next_id"`          //cursor
	NextTimeAtUTC int64  `json:"next_time_at_utc"` //记录分页发生的时间点
	PageSize      int64  `json:"page_size"`        //每一页元素的个数
}

// Isvalid 判断解析的token是否有效的函数
func (p Page) Isinvalid() bool {
	return p.NextID == "" ||
		p.NextTimeAtUTC == 0 ||
		p.NextTimeAtUTC > time.Now().Unix() ||
		p.PageSize <= 0
}

// Encode 返回分页token
func (p Page) Encode() Token {
	b, err := json.Marshal(p)
	if err != nil {
		return Token("")
	}
	return Token(base64.StdEncoding.EncodeToString(b))
}

// Decode 解析分页信息
func (t Token) Decode() Page {
	var result Page
	if len(t) == 0 {
		return result
	}

	bytes, err := base64.StdEncoding.DecodeString(string(t))
	if err != nil {
		return result
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result
	}

	return result
}
