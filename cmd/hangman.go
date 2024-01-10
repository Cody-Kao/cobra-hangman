/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var title string = `
 _   _
| | | | __ _ _ __   __ _ _ __ ___   __ _ _ __
| |_| |/ _` + "`" + ` | '_ \ / _` + "`" + ` | '_ ` + "`" + `_  \ / _` + "`" + ` | '_ \
|  _  | (_| | | | | (_| | | | | | | (_| | | | |
|_| |_|\__,_|_| |_|\__, |_| |_| |_|\__,_|_| |_|
				   ___/` + `_` + `|`

var hidden = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#555555")).
	Faint(true)

var warning = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#ff0000")).
	Bold(true).
	Align(lipgloss.Center)

var highLight = lipgloss.NewStyle().
	Background(lipgloss.Color("#09516B")).
	Foreground(lipgloss.Color("50")).
	Bold(true)

var smallHighLight = lipgloss.NewStyle().
	Foreground(lipgloss.Color("50")).
	Bold(true)

var answer = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#CEBE0D")).
	Bold(true)

var revealAnswer = lipgloss.NewStyle().
	Background(lipgloss.Color("#665e08")).
	Bold(true)

var bold = lipgloss.NewStyle().
	Bold(true)

var myDisplayBorder = lipgloss.Border{
	Top:         "-",
	Bottom:      "-",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "╰",
	BottomRight: "╯",
}

var border = lipgloss.NewStyle().
	Border(myDisplayBorder).
	Padding(1).
	PaddingBottom(0)

var notification = lipgloss.NewStyle().
	Background(lipgloss.Color("#849c1e")).
	Bold(true).
	Align(lipgloss.Center)

var errNotification = lipgloss.NewStyle().
	Background(lipgloss.Color("#ff0000")).
	Bold(true).
	Align(lipgloss.Center)

var question = lipgloss.NewStyle().
	Background(lipgloss.Color("#09516B")).
	Bold(true)

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("50")).
	Width(49).
	PaddingLeft(1).
	PaddingBottom(1)

// Make your own border
var myTitleBorder = lipgloss.Border{
	Top:         "=",
	Bottom:      "=",
	Left:        "||",
	Right:       "||",
	TopLeft:     "=",
	TopRight:    "=",
	BottomLeft:  "=",
	BottomRight: "=",
}

func winAnimation() {
	a := `
	 o
	/|\
	/ \
	   `
	b := `
	 \o/
	  |
	 / \
		 `
	frames := map[string]string{a: b, b: a}
	cur := a
	fmt.Println()
	for i := 0; i < 10; i++ {
		fmt.Print(cur)
		time.Sleep(500 * time.Millisecond)
		fmt.Print(strings.Repeat("\033[1A\x1b[2K", 4))
		cur = frames[cur]
	}
}

func lossAnimation() {
	frames := []string{
		`
 	 _____ 
	|     |
	|     |
	|     |
	|     |
	---   o
	     /|\
	     / \
		    `,
		`
	 _____ 
	|     |
	|     |
	|     |
	|     o
	---  /|\
	      \\  
		    `,
		`
	 _____ 
	|     |
	|     |
	|     o     
	|    /|\	 
	---  //   
			`,
		`
	 _____ 
	|     |
	|     o
	|    /|\
	|     \\
	---       
			`,
		`
	 _____ 
	|     0
	|    /|\
	|     |\
	|     
	---        
			`,
	}
	fmt.Println()
	fmt.Println()
	for _, f := range frames {
		fmt.Print(f)
		time.Sleep(500 * time.Millisecond)
		fmt.Print(strings.Repeat("\033[1A\x1b[2K", 8))
	}
}

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
		log.Fatal(errNotification.Render(err.Error()))
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
	h.numOfCorrect = 0
	h.ans = h.getQuestion()
	l := len(h.ans)
	fmt.Println(question.Render("How many lives you want (you will at least have the ") + highLight.Render("length of puzzle word plus 5 lives") + question.Render("): "))
	_, err := fmt.Scanln(&live)
	if err != nil {
		log.Fatal(errNotification.Render(err.Error() + "Please type in an integer"))
	}
	fmt.Println("========================")
	h.live = max(l+5, live)
	fmt.Printf("You start with %s lives\n", smallHighLight.Render(strconv.Itoa(h.live)))
	h.cur = make([]string, l, l)
	h.record = make(map[string]interface{})
	fmt.Printf("The length of word is: %s\n", smallHighLight.Render(strconv.Itoa(l)))
	fmt.Println(hidden.Render(h.ans))
}

func (h *Hangman) guess() string {
	var g string
	fmt.Println(question.Render("Type in a letter to guess or type ") + highLight.Render("hint ") + question.Render("to get a hint or ") + highLight.Render("quit ") + question.Render("to quit the game or ") + highLight.Render("restart ") + question.Render("to restart:"))
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
		fmt.Println(warning.Render("You ran out of hints!"))
	} else if h.numOfCorrect == len(h.ans)-1 {
		fmt.Println(warning.Render("You can't use hint now"))
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
		fmt.Println(bold.Render("Congrats! you solve the puzzel word: ") + revealAnswer.Render(h.ans))
		winAnimation()
		return true
	}
	if h.live == 0 {
		fmt.Println(bold.Render("Noooo! you fail to solve the puzzel word: ") + revealAnswer.Render(h.ans))
		lossAnimation()
		return true
	}
	text := fmt.Sprint("Your current answer: ")
	for i := 0; i < len(h.ans); i++ {
		if string(h.ans[i]) != h.cur[i] {
			text += fmt.Sprint(answer.Render("_ "))
		} else {
			t := fmt.Sprintf("%s ", string(h.ans[i]))
			t = fmt.Sprintf(answer.Render(t))
			text += t
		}
	}
	text += "\n"
	keys := make([]string, 0, len(h.record))
	for k := range h.record {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	text += fmt.Sprint("Letters have guessed: ")
	for _, l := range keys {
		text += fmt.Sprint(bold.Render(l + " "))
	}
	text += "\n"
	text += fmt.Sprintf("%s Lives left\n", smallHighLight.Render(strconv.Itoa(h.live)))
	text += fmt.Sprintf("%s hints left\n", smallHighLight.Render(strconv.Itoa(h.hint)))
	fmt.Println(border.Render(text))
	return false
}

func (h *Hangman) startGame() {
	var over bool
	var g string
	h.init()
	for {
		g = h.guess()
		if g == "quit" {
			fmt.Println(notification.Render("You quit the game!"))
			break
		}
		if g == "restart" {
			break
		}
		over = h.print()
		if over {
			break
		}
	}
	if g == "restart" {
		fmt.Println(notification.Render("Game is restarted!"))
		h.startGame()
	}
}

func getWord() (string, error) {
	numOfLine := getRandom(0, 25321) // 因為那個字母表的總行數是25322
	fmt.Printf("The random lines of word is: %d\n", numOfLine)
	words := strings.Split(Words, "\n")
	return words[numOfLine], nil
}

func start(cmd *cobra.Command, args []string) {
	titleStyle.BorderStyle(myTitleBorder)
	fmt.Println(titleStyle.Render(title))
	h := &Hangman{}
	h.startGame()
}
