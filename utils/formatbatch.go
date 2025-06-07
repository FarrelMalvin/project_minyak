package utils

import (
	"fmt"
	"time"
)

func FormatBatchName(startDate, endDate time.Time) string {
	return fmt.Sprintf("Analisis Pareto %02d - %02d %s %d",
		startDate.Day(), endDate.Day(),
		endDate.Month().String(), endDate.Year())
}
