package models

import "time"

type Remind struct {
	Text           string
	ActivationTime time.Time
	ChatID         int64
	IDForChat      int
}
