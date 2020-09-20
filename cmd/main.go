package main

import (
	"bufio"
	"fmt"
	"github.com/Rt090/VocabBuilder/vocab"
	"golang.org/x/net/html"

	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	burstSize = 5 // TODO should we move?
)

// entry point, times execution
func main() {
	//filename,err := readEntries()
	//if err != nil {
	//	panic(err)
	//}
	//streak,err := readRequiredStreak()
	//if err != nil {
	//	panic(err)
	//}
	//v, err := vocab.NewVocabulary(filename,burstSize,streak)
	//if err != nil {
	//	panic(err)
	//}
	//new,learned,tough,err := readGroupSizes()
	//if err != nil {
	//	panic(err)
	//}
	//v.Distribute(new,learned,tough)
	//v.Start()
	//v.WriteOut("vocab.json")

	v := &vocab.Vocabulary{}

	http.HandleFunc("/home",func(w http.ResponseWriter,r *http.Request){
		f, err := os.Open("./html/home.html")
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		w.Write(data)
	})
	http.HandleFunc("/submit",func(w http.ResponseWriter,r *http.Request){
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		fmt.Println(r.Form)
		newWords,err := strconv.Atoi(r.Form.Get("newWords"))
		if err != nil {
			panic(err)
		}
		learnedWords,err := strconv.Atoi(r.Form.Get("learnedWords"))
		if err != nil {
			panic(err)
		}
		toughWords,err := strconv.Atoi(r.Form.Get("toughWords"))
		if err != nil {
			panic(err)
		}
		streak,err := strconv.Atoi(r.Form.Get("requiredStreak"))
		if err != nil {
			panic(err)
		}
		filepath := r.Form.Get("filepath")


		f, err := os.Open("./html/list.html")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		doc, err := html.Parse(f)
		if err != nil {
			panic(err)
		}

		v,err = vocab.NewVocabulary(filepath,burstSize,streak)
		if err != nil {
			panic(err)
		}
		v.Distribute(newWords,learnedWords,toughWords)

		rand.Seed(time.Now().UnixNano())

		v.StartWeb()
		words,rem,err := v.NextBatch()

		fmt.Println("rem",rem)
		addRem(doc,rem)
		addForm(doc,words)


		if err != nil {
			panic(err)
		}



		if err = html.Render(w,doc); err != nil {
			panic(err)
		}

	})

	http.HandleFunc("/run",func(w http.ResponseWriter,r *http.Request){

		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		checkMap := map[string]string{}
		for k, v := range r.Form{
			checkMap[k] = v[0]
		}
		correct,answerKey := v.SubmitBatch(checkMap)

		f, err := os.Open("./html/list.html")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		doc, err := html.Parse(f)
		if err != nil {
			panic(err)
		}

		addFeedback(doc,correct,answerKey)
		words,rem,err := v.NextBatch()
		if err != nil {
			panic(err) // really we are just done
		}
		addForm(doc,words)
		addRem(doc,rem)
		if err = html.Render(w,doc); err != nil {
			panic(err)
		}
	})

	http.ListenAndServe(":8080",nil)

}

func addRem(n *html.Node,rem map[string]int) {
	new := findByID(n,"newRem")
	learned := findByID(n,"learnedRem")
	tough := findByID(n,"toughRem")

	if new == nil {
		panic("new is nil")
	}
	if learned == nil {
		panic("learned is nil")
	}
	if tough == nil {
		panic("tough is nil")
	}

	text := &html.Node{}
	text.Type = 1
	text.Data = "New Remaining: " + strconv.Itoa(rem["new"])

	new.AppendChild(text)

	text = &html.Node{}
	text.Type = 1
	text.Data = "Learned Remaining: " + strconv.Itoa(rem["learned"])

	learned.AppendChild(text)

	text = &html.Node{}
	text.Type = 1
	text.Data = "Tough Remaining: " + strconv.Itoa(rem["tough"])

	tough.AppendChild(text)

}

func findByID(n *html.Node, id string) *html.Node {
	var ret *html.Node
	if n == nil {
		return nil
	}
	for _, attr := range n.Attr {
		if attr.Key == "id" && attr.Val == id {
			return n
		}
	}
	for at := n.FirstChild; at != nil; at = at.NextSibling {
		ret = findByID(at,id)
		if ret != nil {
			return ret
		}
	}
	return ret
}

func addFeedback(n *html.Node, correct map[string]bool,answerKey map[string]string) {
	if n.Data == "body" {
		ul := &html.Node{}
		ul.Type = 3
		ul.Data = "ul"

		for word,correct := range correct {
			li := &html.Node{}
			li.Type = 3
			li.Data = "li"
			text := &html.Node{}
			text.Type = 1
			text.Data = word+":"
			if correct{
				text.Data += "Correct!"
			}else {
				text.Data += fmt.Sprintf("Wrong! Wanted:%s",answerKey[word])
			}
			li.AppendChild(text)
			ul.AppendChild(li)
		}


		appendBr(n,2)

		n.AppendChild(ul)
	}
	for at := n.FirstChild; at != nil; at = at.NextSibling {
		addFeedback(at,correct,answerKey)
	}
}

func appendBr(node *html.Node,n int) {
	for i := 0; i < n; i++ {
		br := &html.Node{}
		br.Type = 3
		br.Data = "br"
		node.AppendChild(br)
	}
}

 func addForm (n *html.Node, words []string)  {
	i := 0
	for _, attr := range n.Attr {
		if attr.Key == "id" && attr.Val == "vocabList" {
			for _, word := range words {
				input := &html.Node{}
				// TODO: potential bug if we have the same word twice, shouldn't happen but is assumed unique here
				input.Attr = []html.Attribute{html.Attribute{Key:"name",Val:word},html.Attribute{Key:"type",Val: "text"}}
				input.Data = "input"
				input.Type = 3


				label := &html.Node{}
				label.Type = 3
				label.Data = "label"
				label.Attr = []html.Attribute{html.Attribute{Key: "for",Val:"word " + strconv.Itoa(i)}}

				text := &html.Node{}
				text.Type = 1
				text.Data = word
				label.AppendChild(text)

				br := &html.Node{}
				br.Type = 3
				br.Data = "br"

				n.AppendChild(label)
				n.AppendChild(input)
				n.AppendChild(br)

				i++
			}
			br := &html.Node{}
			br.Type = 3
			br.Data = "br"
			n.AppendChild(br)
			submit := &html.Node{}
			submit.Attr = []html.Attribute{html.Attribute{Key:"type",Val:"submit"},html.Attribute{Key:"value",Val: "Submit"}}
			submit.Data = "input"
			submit.Type = 3
			n.AppendChild(submit)
			return
		}
	}

	for at := n.FirstChild; at != nil; at = at.NextSibling {
		addForm(at,words)
	}
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
