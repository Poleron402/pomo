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
			duration = 10
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
	isBreak := false
	
	done := make(chan struct{})
	fmt.Printf("Working on %v/%v schedule... Press 'q' to stop <<", duration, rest) 
	go func(){
		for {
			greenBG := "\033[1;42m"
			whiteTxt := "\033[97m"
			reset := "\033[0m"
			
			now := time.Now()
			tick := time.NewTicker(1 * time.Second)
			
			myWorkTime := duration
			prompt := "Working timer: "
			if isBreak {
				myWorkTime = rest
				prompt = "Resting timer: "
			}
			OuterLoop:
			for {
				select {
					case <-done:
						fmt.Println("\nStopping timer...")
						return
					case t := <-tick.C:

						elapsed := t.Sub(now)
						minutes := int(elapsed.Minutes())
						seconds := int(elapsed.Seconds()) % 60
						if myWorkTime > 60 {
							if float64(minutes*60)/float64(myWorkTime) >= 0.95 {
								greenBG = "\033[1;41m"
							}else if float64(minutes*60)/float64(myWorkTime) >= 0.75 {
								greenBG = "\033[1;43m"
							}
						}else {
							if float64(seconds)/float64(myWorkTime) >= 0.90 {
								greenBG = "\033[1;41m"
							}else if float64(seconds)/float64(myWorkTime) >= 0.76{
								greenBG = "\033[1;43m"
							}
						}

						fmt.Printf("\r\033[2K%s %s%v%02d:%02d %s", greenBG, whiteTxt, prompt, minutes, seconds, reset)
						
						totalWorkTime = int(elapsed.Seconds())

						if int(elapsed.Seconds()) >= myWorkTime {
							fmt.Println("\nTime's up! 🍅")
							tick.Stop()
							playAlarm()
							isBreak = !isBreak
							break OuterLoop
						}
				}
			}
		}
	}()
	option := ""

	fmt.Scan(&option)
	if option == "q" {
		close(done)
		return
	}
	fmt.Printf("Nice session! Your total work time was: %v minutes\n", totalWorkTime/60)
	
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
