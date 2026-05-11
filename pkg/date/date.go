package date

import (
	"time"
)

func ExpireDateString(notafter time.Time) string {
	return notafter.Format("2006-01-02")
}

func ExpireDiffDays(notafter time.Time) int {
	days := int(time.Until(notafter).Hours() / 24)
	if days < 1 {
		return 1
	}
	return days
}
