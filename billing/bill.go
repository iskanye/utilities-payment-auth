package billing

import "time"

type Bill struct {
	BillId int64
	UserId int64
	Sum    int
	DueTo  time.Time
}

func (b *Bill) Pay() {
}
