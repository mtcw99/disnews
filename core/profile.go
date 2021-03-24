package core

import (
	"time"
)

type Profile struct {
	Username     string
	DisplayName  string
	Info         string
	Link         string
	CreationDate time.Time
}

//func (p *Profile) CreationDateString() string {
//	return ""
//}
