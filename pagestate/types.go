package pagestate

import "time"

type Pagestate struct {
	Id        int
	Url       string
	ScrollPos int
	CreatedAt time.Time
}
