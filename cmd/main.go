package main

import (
	"bufio"
	"os"
	"fmt"
	"github.com/Rt090/VocabBuilder/vocab"

)

// entry point, times execution
func main() {
	filename,err := readEntries()
	if err != nil {
		panic(err)
	}
	v, err := vocab.NewVocabular(filename)
}
// prompt user and return numbers we should use for new,learned, and tough
func readGroupSizes(){}
// prompt user and return number of times we must correctly answer consecutively to pass
func readRequiredStreak(){}
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
