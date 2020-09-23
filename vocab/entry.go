package vocab

import "time"

const (
	STATE_NEW = iota
	STATE_LEARNED = iota
	STATE_TOUGH = iota
	STATE_MASTERED = iota
)

type entry struct{
	Eng string `json:"eng"`
	Kor string `json:"kor"`
	Attempts []*attempt `json:"attempts"`
	LastAttempted *time.Time `json:"lastAttempted"`
	Remedial *remedial `json:"remedial"`
	Mastered *mastered `json:"mastered"`
}

type remedial struct{
	Attempts []*attempt `json:"attempts"`
	InProgress bool `json:"inProgress"`
}
type mastered struct{
	Added []*time.Time `json:"added"`
	Removed []*time.Time `json:"removed"`
}
type attempt struct{
	Ts *time.Time `json:"ts"`
	Misses int `json:"misses"`
	Required int `json:"required"`
}

func (e *entry) state() int{
	if len(e.Attempts) < 7 {
		return STATE_NEW
	}
	if e.Mastered != nil && len(e.Mastered.Added) != len(e.Mastered.Removed) {
		return STATE_MASTERED
	}
	if e.Remedial != nil && e.Remedial.InProgress {
		return STATE_TOUGH
	}
	return STATE_LEARNED
}