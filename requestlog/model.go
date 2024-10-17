package requestlog

import "time"

type RequestLog struct {
	ID         uint   `gorm:"primaryKey"`
	Method     string `gorm:"size:8"`
	Host       string `gorm:"size:128"`
	Path       string `gorm:"size:256"`
	Query      string `gorm:"size:512"`
	Header     string `gorm:"type:text"`
	Body       string `gorm:"type:text"`
	Response   string `gorm:"type:text"`
	StatusCode int
	CreatedAt  time.Time
}
