package main

import "fmt"

//The way we avoid a deadlock, is to make sure that philosophers prioritise forks on different sides based on their id.
//This way, we avoid the situation where all philosophers pick up the fork on their left side, and then wait for the fork on their right side.

//If a philosopher picks up the fork on their left side, then the philosopher who shares that fork will be unable to pick up that fork.
//Since that fork is their prioritised fork they will not hold any forks while waiting, thus avoiding the deadlock.

func philosopher(leftUp, leftDown, rightUp, rightDown chan bool, id int) {
	eatCount := 0
	//Don't terminate until all philosophers have eaten x times
	for eatCount < 4 {
		//Pick up the forks
		if id%2 == 0 {
			leftUp <- true
			rightUp <- true
		} else {
			rightUp <- true
			leftUp <- true
		}
		//Eat
		eatCount++
		fmt.Println("Philosopher", id, "ate", eatCount, "times")
		//Put down the forks
		leftDown <- false
		rightDown <- false
		fmt.Println("Philosopher", id, "put down forks and is now thinking")
	}
	fmt.Println("PHILOSOPHER", id, "IS DONE EATING")
}

func fork(ch1, ch2 chan bool) {
	//Don't terminate until all philosophers have eaten 100 times
	for {
		//Wait for a philosopher to pick up the fork
		<-ch1
		//Wait for a philosopher to put down the fork
		<-ch2
	}

}

func main() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	ch3 := make(chan bool)
	ch4 := make(chan bool)
	ch5 := make(chan bool)
	ch6 := make(chan bool)
	ch7 := make(chan bool)
	ch8 := make(chan bool)
	ch9 := make(chan bool)
	ch10 := make(chan bool)

	//Start the forks first, then philosophers
	go fork(ch1, ch2)
	go fork(ch3, ch4)
	go fork(ch5, ch6)
	go fork(ch7, ch8)
	go fork(ch9, ch10)

	go philosopher(ch9, ch10, ch1, ch2, 0)
	go philosopher(ch1, ch2, ch3, ch4, 1)
	go philosopher(ch3, ch4, ch5, ch6, 2)
	go philosopher(ch5, ch6, ch7, ch8, 3)
	go philosopher(ch7, ch8, ch9, ch10, 4)

	//Don't terminate until user input
	var exit string
	fmt.Scanln(&exit)
}
