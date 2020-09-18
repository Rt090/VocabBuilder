package vocab

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
	"fmt"
)

type Vocabulary struct {
	burst map[string]*task
	new []*entry
	learned []*entry
	tough []*entry
	consolidated map[string]*entry
	burstSize int
	successReq int
}
// TODO for consolidated use a new thing besides entry, representing korean in map[string]

type task struct {
	attempts int
	correctSequential int
}

func NewVocabulary(filepath string,burstSize,success int) (*Vocabulary,error){
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

	v := &Vocabulary{successReq: success,burstSize: burstSize}
	v.loadAll(entries)
	return v, nil
}

func (v *Vocabulary) loadAll(entries []*entry) {
	newestList := make([]*entry,0,len(entries))
	learned := make([]*entry,0,len(entries))
	tough := make([]*entry,0,len(entries))

	all := make(map[string]*entry,len(entries))

	for _, e := range entries{
		if _, ok := all[e.Eng]; !ok {
			all[e.Eng] = e
		}else {
			k1 := all[e.Eng].Kor
			k2 := e.Kor

			if k1 != k2 {
				// we may have already joined with /, so break apart again
				all := strings.Split(k1,"/")
				for _, option := range all {
					if k2 == option { // if one entry matches, just move on
						continue
					}
				}
				k1 = strings.Join([]string{k1,k2},"/")
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
	v.new = newestList
	v.tough = tough
	v.learned = learned

	fmt.Printf("New Words:%d, Learned Words:%d, Tough Words:%d\n",len(newestList),len(learned),len(tough))
}

func (v *Vocabulary) Distribute(newestCount,learnedCount,toughCount int) {

	cur := []*entry{}
	cmp := func (a,b int) bool{
		if cur[a].LastAttempted == nil {
			return true
		}
		if cur[b].LastAttempted == nil {
			return false
		}
		return cur[a].LastAttempted.Before(*cur[b].LastAttempted)
	}

	// TODO make this cleaner with the count comp
	cur = v.new
	sort.Slice(cur,cmp)
	if len(cur) > newestCount {
		v.new = cur[:newestCount]
	}else {
		v.new = cur
	}


	cur = v.tough
	sort.Slice(cur,cmp)
	if len(v.tough) > toughCount{
		v.tough = cur[:toughCount]
	}else {
		v.tough = cur
	}

	cur = v.learned
	sort.Slice(cur,cmp)
	if len(v.learned) > learnedCount{
		v.learned = cur[:learnedCount]
	}else {
		v.learned = cur
	}


}

// TODO should this talk directly to user
func (v *Vocabulary) Start() {
	rand.Seed(time.Now().UnixNano())

	for rem := v.loadNextBurst(); rem > 0;rem = v.loadNextBurst() {
		err := v.completeBurst()
		if err != nil {
			v.WriteOut("vocab.json")
			panic(err)
		}
		v.writeOutBurst()
	}
}

func (v *Vocabulary) completeBurst() error{
	skipped := 0
	for skipped < v.burstSize{
		l := []string{}
		for eng := range v.burst{
			l = append(l,eng)
		}
		rand.Shuffle(len(l),func(i,j int){l[i],l[j] = l[j], l[i]})
		skipped = 0
		for _,eng := range l {
			if v.burst[eng].correctSequential >= v.successReq {
				skipped++
				if skipped == len(l) {
					return nil
				}
				continue
			}
			entry := v.consolidated[eng]
			answer,err := getAnswer(entry.Eng)
			if err != nil {
				return err
			}
			v.burst[eng].attempts ++

			// TODO should we just build this off the rip? makes entry object ugly, but what's stored shouldnt
			// be equivalent to what is running

			allAnswers := strings.Split(entry.Kor,"/")
			correctMap := map[string]struct{}{}
			for _, a := range allAnswers {
				correctMap[a] = struct{}{}
			}

			if _,ok := correctMap[answer]; !ok {
				v.burst[eng].correctSequential = 0
				fmt.Printf("Incorrect: looking for %s\n",entry.Kor)
			}else {
				v.burst[eng].correctSequential++
				fmt.Printf("Correct! Attempts:%d, Correct in a row:%d\n",v.burst[eng].attempts,v.burst[eng].correctSequential)
			}
		}
	}
	return nil
}

func getAnswer(input string) (string,error) {
	reader := bufio.NewReader(os.Stdin) // TODO don't make a new one every time
	fmt.Printf("Enter the korean for the english '%s'\n",input)
	output,err := reader.ReadString('\n')
	if err != nil {
		return "",err
	}
	output = strings.TrimSpace(output)
	return output,err
}

func (v *Vocabulary) writeOutBurst() {
	for eng, task := range v.burst {
		cur := v.consolidated[eng]
		misses := task.attempts - v.successReq
		ts := time.Now()
		attempt := &attempt{Misses: misses,Ts: &ts}
		cur.Attempts = append(cur.Attempts,attempt)
		cur.LastAttempted = &ts
	}
}

func (v *Vocabulary) loadNextBurst() int{
	b := make(map[string]*task,v.burstSize)
	for i := 0; i < v.burstSize; i++ {
		if len(v.new) > 0 {
			b[v.new[0].Eng] = &task{}
			v.new = v.new[1:]
		} else if len(v.learned) > 0 {
			b[v.learned[0].Eng] = &task{}
			v.learned = v.learned[1:]
		} else if len(v.tough) > 0 {
			b[v.tough[0].Eng] = &task{}
			v.tough = v.tough[1:]
		}
	}
	v.burst = b
	return len(b)
}

// write out the
func (v *Vocabulary) WriteOut(filepath string)error{
	l := make([]*entry,0,len(v.consolidated))
	for _, e := range v.consolidated {
		l = append(l,e)
	}
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}
	ioutil.WriteFile(filepath, j, 0644)
	return nil
}