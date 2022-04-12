package dashboard

import (
	"fmt"
	"time"
)

func nice(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%d", n)
	}
	return fmt.Sprintf("%d", n)
}

func GetTimePassedSince(createdAt time.Time, withSeconds bool) string {
	seconds := int(time.Now().Sub(createdAt).Seconds()) % 60
	minutes := int(time.Now().Sub(createdAt).Minutes())
	hours := int(minutes / 60)
	if !withSeconds {
		return fmt.Sprintf("%s:%s", nice(hours), nice(minutes))
	}
	return fmt.Sprintf("[%s:%s:%s]", nice(hours), nice(minutes), nice(seconds))
}
