package data

import "time"

// TO DO : THink about the data or the dish like what to keep and what not to keep
type Dishses struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Type      string // veg or Non veg
	Version   int32
}
