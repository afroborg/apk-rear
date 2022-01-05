package models

import "time"

type Status struct {
	Time           time.Time `json:"time"`
	SyncedProducts int       `json:"syncedProducts"`
}
