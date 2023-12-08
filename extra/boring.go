
//description:
//example of concurrency failing: 
// we attempt to run a goroutine that prints msg per second
// but as soon as the main function reaches the end of it's output the entire thing terminates (nothing printed)
// example of concurrent porgram

package main
import {
	"fmt"
	"math/rand"
	"time"
}

func boring(msg string){
	for i :=0 ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.second)
	}
}

func boring(msg string){
	for i :=0 ; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.second)
	}
}

func main() {
	go boring("boring!")
	fmt.Println("I'm listening.")
	time.Sleep(2* time.Second)
	fmt.Println("I'm no longer listening")
}
