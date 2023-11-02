package data

import "time"

type DiffLogItem struct {
	Timestamp time.Time
	Revision  string
	Content   string
}
