package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type empty struct {
	English string `json:"eng"`
	Korean string `json:"kor"`
	Attempts []*attempt `json:"attempts"`
}

type attempt struct{
	Ts *time.Time `json:"ts"`
	Misses int `json:"misses"`
}

func main() {
	file, err := os.Open("new.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	all := []*empty{}
	scanner := bufio.NewScanner(file)
	now := time.Now()
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.Replace(line," ,",",",-1)
		line = strings.Replace(line,", ",",",-1)
		line = strings.Replace(line," /","/",-1)
		line = strings.Replace(line,"/ ","/",-1)
		tup := strings.Split(line,",")
		kor,eng := tup[0],tup[1]
		buffer, err := strconv.Atoi(tup[2])
		if err != nil {
			fmt.Println(count)
			panic(err)
		}
		attempts := make([]*attempt,0,buffer)
		mock := &attempt{Ts:&now,Misses: 0}
		for i := 0; i < buffer; i++ {
			attempts = append(attempts,mock)
		}
		e := &empty{English: eng,Korean: kor,Attempts: attempts}
		all = append(all,e)
		count += 1
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	j, err := json.Marshal(all)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("new.json", j, 0644)
}