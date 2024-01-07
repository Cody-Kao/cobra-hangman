/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// hangmanCmd represents the hangman command
var hangmanCmd = &cobra.Command{
	Use:   "start",
	Short: "It's a start command",
	Long:  `The words contain only LOWERCASE English letters. Please enjoy a small but delicate game!`,
	Run:   start,
}

func init() {
	rootCmd.AddCommand(hangmanCmd)

	// Here you will define your flags and configuration settings.

	// hangmanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Hangman struct {
	live         int
	hint         int
	cur          []string
	numOfCorrect int
	record       map[string]interface{}
	ans          string
}

func (h *Hangman) getQuestion() string {
	word, err := getWord()
	if err != nil {
		log.Fatal(err)
	}
	return word
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func (h *Hangman) init() {
	var live int // 先說好要得到的是什麼type，就不會都拿到string type，但也有可能使用者輸入的是不能轉成int的東西，所以要偵錯
	h.hint = 2
	h.ans = h.getQuestion()
	l := len(h.ans)
	fmt.Println("How many lives you want (you will at least have the length of puzzle word plus 5 lives): ")
	_, err := fmt.Scanln(&live)
	if err != nil {
		log.Fatal(err, "Please type in an integer")
	}
	h.live = max(l+5, live)
	fmt.Printf("You start with %d lives\n", h.live)
	h.cur = make([]string, l, l)
	h.record = make(map[string]interface{})
	fmt.Printf("The length of word is: %d\n", l)
	fmt.Println(h.ans)
}

func (h *Hangman) guess() string {
	var g string
	fmt.Println("Type in a letter to guess or type 'hint' to get a hint or 'quit' to quit the game or 'restart' to restart: ")
	fmt.Scanln(&g)
	if g == "hint" {
		h.getHint()
		return ""
	}
	if g == "quit" || g == "restart" {
		return g
	}
	h.record[g] = nil
	h.live -= 1
	h.check(g)
	return ""
}

func (h *Hangman) getHint() {
	if h.hint == 0 {
		fmt.Println("You ran out of hints!")
	} else if h.numOfCorrect == len(h.ans)-1 {
		fmt.Println("You can't use hint now")
	} else {
		h.hint -= 1
		for {
			index := getRandom(0, len(h.ans)-1)
			if h.cur[index] == "" {
				h.cur[index] = string(h.ans[index])
				h.numOfCorrect += 1
				break
			}
		}
	}
}

func getRandom(min int, max int) int {
	src := rand.NewSource(time.Now().UnixNano()) // return rand.Source to generate a random source based on the given seed
	rng := rand.New(src)                         // return *rand.Rand that can generate random number based on the given source
	// Intn returns, as an int, a non-negative pseudo-random number in the half-open interval [0,n). It panics if n <= 0.
	// 因為他預設是從0 ~ max，所以我們只要加上min就會變成從min ~ max
	index := rng.Intn(max) + min
	return index
}

func (h *Hangman) check(g string) {
	for i := 0; i < len(h.ans); i++ {
		if string(h.ans[i]) == g && h.cur[i] == "" {
			h.cur[i] = g
			h.numOfCorrect += 1
		}
	}
}

func (h *Hangman) print() bool {
	if h.numOfCorrect == len(h.ans) {
		fmt.Printf("Congrats! you solve the puzzel word: %s", h.ans)
		return true
	}
	if h.live == 0 {
		fmt.Printf("Noooo! you fail to solve the puzzel word: %s", h.ans)
		return true
	}
	fmt.Print("Your current answer: ")
	for i := 0; i < len(h.ans); i++ {
		if string(h.ans[i]) != h.cur[i] {
			fmt.Print("_ ")
		} else {
			fmt.Printf("%s ", string(h.ans[i]))
		}
	}
	fmt.Println()
	fmt.Print("Letters have guessed: ")
	for k := range h.record {
		fmt.Print(k, " ")
	}
	fmt.Println()
	fmt.Printf("%d Lives left\n", h.live)
	fmt.Printf("%d hints left\n", h.hint)
	return false
}

func (h *Hangman) startGame() {
	var over bool
	var g string
	h.init()
	for {
		g = h.guess()
		if g == "quit" {
			fmt.Println("You quit the game!")
			break
		}
		if g == "restart" {
			break
		}
		fmt.Println("----------------")
		over = h.print()
		if over {
			break
		}
	}
	if g == "restart" {
		fmt.Println("Game is restarted!")
		h.startGame()
	}
}

func getWord() (string, error) {
	f, err := os.Open("./cmd/static/words.txt")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	numOfLine := getRandom(0, 25321) // 因為那個字母表的總行數是25322
	fmt.Printf("The random lines of word is: %d\n", numOfLine)
	var i int
	for scanner.Scan() && i < numOfLine { // Read line by line
		i++
	}
	return scanner.Text(), nil
}

func start(cmd *cobra.Command, args []string) {
	h := &Hangman{}
	h.startGame()
}
