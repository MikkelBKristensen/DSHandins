package main

import (
	"fmt"
	"math/rand"
)

/*
 1) Implement TCP/IP using threads. The implementation needs to show that we have a good understanding of the protocol.
*/
//random number generator
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func sendSyn(ch chan int) int {
	synNum := random(0, 1000)
	ch <- synNum
	println("Sending SYN: ", synNum)
	return synNum
}

func confirmAck(mySyn, ack int) bool {
	if (mySyn + 1) != ack {
		return false
	}
	print("Received ACK: ", ack, "\n")
	return true
}

func senderInitialise(ch chan int) bool {
	mySyn := sendSyn(ch) //send syn
	synReceived := <-ch
	if confirmAck(mySyn, <-ch) {
		mySyn = mySyn + 1
		ch <- mySyn           //send seq
		ch <- synReceived + 1 //send ack
		ch <- mySyn + 1       //send seq
		return true
	} else {
		fmt.Println("Ack not received - RST")
		return false
	}
}

func receiverInitialise(ch chan int) bool {
	synReceived := <-ch
	mySyn := sendSyn(ch)  //send own syn
	ch <- synReceived + 1 //send ack
	if <-ch != synReceived+1 {
		return false
	}
	if confirmAck(mySyn, <-ch) {
		fmt.Println("Connection established")
		return true
	} else {
		fmt.Println("Ack not received - RST")
		return false
	}
}

func sender(ch chan int) {
	if senderInitialise(ch) {
		data := 754
		ch <- data //send data
	}
}

func receiver(ch chan int) {
	if receiverInitialise(ch) {
		receiveSeq := <-ch //receive seq
		println("Received seq: ", receiveSeq)
		receivedData := <-ch //receive data
		fmt.Println("Received data: ", receivedData)
	}
}

func main() {

	ch := make(chan int)
	go sender(ch)
	go receiver(ch)

	var exit string
	fmt.Scan(&exit)
	fmt.Println(exit)

}
