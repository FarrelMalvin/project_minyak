package models

import "time"

type ParetoBatch struct {
	BatchID        uint      `gorm:"primaryKey;column:batch_id" json:"batch_id"`
	StartDate      time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate        time.Time `gorm:"column:end_date" json:"end_date"`
	Summary        string    `gorm:"column:summary" json:"summary"`
	Recommendation string    `gorm:"column:recommendation" json:"recommendation"`
}

func (ParetoBatch) TableName() string {
	return "pareto_analysis_batch"
}
