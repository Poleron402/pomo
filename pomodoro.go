package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const (
	_ = iota
	TimeSelect
	ContinueBreak
	ContinueWork
	StopWork
)

func prompt(prompt int) (int, int) {
	var option int
	var duration int
	var rest int
	if prompt == 1 {
		fmt.Print("Enter how long you would like the pomo to be << \n (1) 25/5 \n (2) 52/7 \n (3) 90/30 \n (4) Quit \n >> ")
		fmt.Scan(&option)

		switch option {
		case 1:
			duration = 25 * 60
			rest = 5 * 60
		case 2:
			duration = 52 * 60
			rest = 7 * 60
		case 3:
			duration = 90 * 60
			rest = 30 * 60
		case 4:
			return -1, -1
		case 5: // test
			duration = 1 * 60
			rest = 0.5 * 60
		}
		return duration, rest
	}

	return -1, -1
}

func pomodoro() {
	duration, rest := prompt(1)
	totalWorkTime := 0
	if duration == -1 || rest == -1 {
		return
	}
	now := time.Now()
	isBreak := false
	// defer tick.Stop()
	done := make(chan struct{})
	fmt.Println("Press 'q' to stop <<") 
	go func(){
		for {
			tick := time.NewTicker(1 * time.Second)

			greenBG := "\033[1;42m"
			whiteTxt := "\033[97m"
			reset := "\033[0m"

			myWorkTime := duration
			if isBreak {
				myWorkTime = rest
			}
			
			for range myWorkTime {
				t := <-tick.C
				elapsed := t.Sub(now)
				minutes := int(elapsed.Minutes())
				seconds := int(elapsed.Seconds()) % 60
				if myWorkTime > 60 {
					if float64(minutes*60)/float64(myWorkTime) >= 0.76 {
						greenBG = "\033[1;43m"
					}else if float64(minutes*60)/float64(myWorkTime) >= 0.95 {
						greenBG = "\033[1;91m"
					}
				}else if float64(seconds)/float64(myWorkTime) >= 0.8{
					greenBG = "\033[1;43m"
				}
				fmt.Printf("\r\033[2K%s %sPassed: %02d:%02d %s", greenBG, whiteTxt, minutes, seconds, reset)
				
				totalWorkTime += int(elapsed.Seconds())
			}
			fmt.Println("\nTime's up! 🍅")
			tick.Stop()
			playAlarm()
			isBreak = !isBreak // negating to alternate
		}
	}()
	fmt.Printf("\nPress 'q' to stop << ")
	option := ""
	// fmt.Printf("Press 'q' to stop << ")
	fmt.Scan(&option)
	if option == "q" {
		close(done)
	}
	if option == "q" {
		fmt.Printf("Nice session! Your total time was: %v minutes", totalWorkTime/60)
		return
	}
}

func playAlarm() {
	f, err := os.Open("notify.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Fatal(err)
	}
	c, err := oto.NewContext(d.SampleRate(), 1, 2, 8192)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	io.Copy(p, d)
}
