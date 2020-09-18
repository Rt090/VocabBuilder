package main

import (
	"bufio"
	"os"
	"fmt"
	"github.com/Rt090/VocabBuilder/vocab"
	"strconv"
	"strings"
)

const (
	burstSize = 5 // TODO should we move?
)

// entry point, times execution
func main() {
	filename,err := readEntries()
	if err != nil {
		panic(err)
	}
	streak,err := readRequiredStreak()
	if err != nil {
		panic(err)
	}
	v, err := vocab.NewVocabulary(filename,burstSize,streak)
	if err != nil {
		panic(err)
	}
	new,learned,tough,err := readGroupSizes()
	if err != nil {
		panic(err)
	}
	v.Distribute(new,learned,tough)
	v.Start()
	v.WriteOut("vocab.json")

}
// prompt user and return numbers we should use for new,learned, and tough
func readGroupSizes()(int,int,int,error){
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter how many new words to tackle today: ")
	newCount,err := readInt(reader)
	if err != nil {
		return -1,-1,1,err
	}
	fmt.Println("Enter how many learned words to tackle today: ")
	learnedCount,err := readInt(reader)
	if err != nil {
		return -1,-1,1,err
	}
	fmt.Println("Enter how many tough words to tackle today: ")
	toughCount,err := readInt(reader)
	if err != nil {
		return -1,-1,1,err
	}

	return newCount,learnedCount,toughCount,nil
}

func readInt(reader *bufio.Reader) (int,error) {
	text, err := reader.ReadString('\n')
	if err != nil {
		return -1,err
	}
	text = strings.TrimSpace(text)
	i,err := strconv.Atoi(text)
	if err != nil {
		return 0,err
	}
	return i,nil
}
// prompt user and return number of times we must correctly answer consecutively to pass
func readRequiredStreak()(int,error){
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter required streak: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return -1,err
	}
	text = strings.TrimSpace(text)
	num,err := strconv.Atoi(text)
	if err != nil {
		return 0,err
	}
	fmt.Println("requiring string of ",num)
	return num,nil
}
// return the group sizes for new,learned,tough,mastered and all
func countAllGroups(){}

// read all saved entries from file, deduping and removing pronunciation
func readEntries()(string,error){
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("--- Welcome to Vocabulary Learner ---")
	fmt.Println("Enter filename: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "",err
	}
	text = strings.TrimSpace(text)
	fmt.Println("Opening ",text)
	return text,nil
}
// transform an existing list of words into JSON to append to existing file
func csvToJSON(){}

// prompt user for if we should send the entry that is in tough back to learned
func shouldSendToLearned(){}
// prompt user for if we should send the entry that is in learned to tough
func shouldSendToTough(){}
// allow user to pull stats for single entry
func searchForEntry(){}
