package main

import "fmt"

func philosopher(ch1, ch2 chan bool) {
	eatCount := 0
	shouldThink := true
	for !(eatCount > 3) {
		select {
		case <-ch1:
			ch1 <- false
			if <-ch2 {
				ch2 <- false
				eatCount++
				fmt.Println("This philosopher has eaten ", eatCount, " times!")
				shouldThink = true
				ch2 <- true
			} else {
				ch2 <- false
			}
			ch1 <- true
			if !shouldThink {
				fmt.Println("This philosopher is currently philosophising...")
				shouldThink = false
			}
		case <-ch2:
			ch2 <- false
			if <-ch1 {
				ch1 <- false
				eatCount++
				fmt.Println("This philosopher has eaten ", eatCount, " times!")
				ch1 <- true
			} else {
				ch1 <- false
			}
			ch2 <- true
		default:
			ch1 <- false
			ch2 <- false

		}
	}
}

func fork() {

}

func main() {
	ch1 := make(chan bool, 3)
	ch2 := make(chan bool, 3)
	ch3 := make(chan bool, 3)
	ch4 := make(chan bool, 3)
	ch5 := make(chan bool, 3)
	ch1 <- true
	ch2 <- true
	ch3 <- true
	ch4 <- true
	ch5 <- true

	go philosopher(ch1, ch2)
	go fork()
	go philosopher(ch2, ch3)
	go fork()
	go philosopher(ch3, ch4)
	go fork()
	go philosopher(ch4, ch5)
	go fork()
	go philosopher(ch5, ch1)
	go fork()
	var ost string
	fmt.Scan(&ost)
	fmt.Println(ost)
}
