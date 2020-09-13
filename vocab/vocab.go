package vocab

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

type Vocabulary struct {
	Entries []*entry `json:"entries"`
	burst map[string]int
	new []*entry
	learned []*entry
	tough []*entry
	consolidated map[string]*entry
	burstSize int
	successReq int
}

func NewVocabulary(filepath string,burstSize,new,learned,tough,success int) (*Vocabulary,error){
	f, err := os.Open(filepath)
	if err != nil {
		return nil,err
	}
	defer f.Close()

	entries := []*entry{}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil,err
	}
	if err = json.Unmarshal(b,&entries); err != nil {
		return nil, err
	}

	v := &Vocabulary{Entries: entries,successReq: success,burstSize: burstSize}
	v.distribute(new,learned,tough)
	return v, nil
}

func (v *Vocabulary) distribute(newestCount,learnedCount,toughCount int) {
	newestList := make([]*entry,newestCount) // TODO move to size of entries
	learned := make([]*entry,learnedCount)
	tough := make([]*entry,toughCount)

	all := make(map[string]*entry,len(v.Entries))

	for _, e := range v.Entries{
		if _, ok := all[e.Eng]; !ok {
			all[e.Eng] = e
		}else {
			k1 := all[e.Eng].Kor
			k2 := e.Kor

			if k1 != k2 {
				k1 = strings.Join([]string{k1,k2},",")
			}
			all[e.Eng].Kor = k1
		}
	}
	v.consolidated = all

	for _, e := range all {
		switch e.state() {
		case STATE_NEW:
			newestList = append(newestList,e)
		case STATE_LEARNED:
			learned = append(learned,e)
		case STATE_TOUGH:
			tough = append(tough,e)
		case STATE_MASTERED:

		}

	}

	var cur []*entry
	cmp := func (a,b int) bool{
		return cur[a].LastAttempted.Before(*cur[b].LastAttempted)
	}

	cur = newestList
	sort.Slice(cur,cmp)
	newestList = cur[:newestCount]

	cur = tough
	sort.Slice(cur,cmp)
	tough = cur[:toughCount]

	cur = learned
	sort.Slice(cur,cmp)
	learned = cur[:learnedCount]

	v.new = newestList
	v.tough = tough
	v.learned = learned
}

func (v *Vocabulary) writeOutBurst() {
	for eng, tries := range v.burst {
		cur := v.consolidated[eng]
		misses := tries - v.successReq
		ts := time.Now()
		attempt := &attempt{Misses: misses,Ts: &ts}
		cur.Attempts = append(cur.Attempts,attempt)
	}
}

func (v *Vocabulary) loadNextBurst()int{
	b := make(map[string]int,v.burstSize)
	for i := 0; i < v.burstSize; i++ {
		if len(v.new) > 0 {
			b[v.new[0].Eng] = 0
			v.new = v.new[1:]
		} else if len(v.learned) > 0 {
			b[v.learned[0].Eng] = 0
			v.learned = v.learned[1:]
		} else if len(v.tough) > 0 {
			b[v.tough[0].Eng] = 0
			v.tough = v.tough[1:]
		}
	}
	return len(b)
}

// write out the
func (v *Vocabulary) WriteOut(filepath string)error{
	l := make([]*entry,len(v.consolidated))
	for _, e := range v.consolidated {
		l = append(l,e)
	}
	v.Entries = l

	j, err := json.Marshal(v)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filepath, j, 0644)
	return nil
}