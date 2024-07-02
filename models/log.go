package models

import "gorm.io/gorm"

type RequestLog struct {
	gorm.Model
	Method       string `json:"method"`
	Path         string `json:"path"`
	StatusCode   int    `json:"status_code"`
	Latency      int64  `json:"latency"` // in milliseconds
	ClientIP     string `json:"client_ip"`
	UserAgent    string `json:"user_agent"`
	RequestBody  string `json:"request_body" gorm:"type:text"`
	ResponseBody string `json:"response_body" gorm:"type:text"`
}
