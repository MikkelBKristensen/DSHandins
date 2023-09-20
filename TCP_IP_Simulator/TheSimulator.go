package TCP_IP_Simulator

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

func sendSyn(ch chan int) {
	synNum := random(0, 1000)
	ch <- synNum
}

func sender(ch chan int) {

}

func receiver() {

}

func main() {

	fmt.Println(random(0, 1000))

}
