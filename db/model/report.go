package model

import (
	"time"
)

// ReportModel 每日报告模型
type ReportModel struct {
	Id          int64     `gorm:"column:id" json:"id"`
	ReportFrom  string    `gorm:"column:report_from" json:"reportFrom"`
	ReportType  string    `gorm:"column:report_type" json:"reportType"`
	ReportTitle string    `gorm:"column:report_title" json:"reportTitle"`
	ReportName  string    `gorm:"column:report_name" json:"reportName"`
	ReportDate  time.Time `gorm:"column:report_date" json:"reportDate"`
	ReportUrl   string    `gorm:"column:report_url" json:"reportUrl"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
