package main

import (
	"fmt"
	"os"
)

func calc(s, e int, writeChan chan string, done chan struct{}) {
	for i := s; i <= e; i++ {
		line := fmt.Sprintf("%06x\n", i)
		writeChan <- line
	}
	done <- struct{}{}
}

func main() {
	f, err := os.Create("f.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	grt := 10
	fn := 1677215
	done := make(chan struct{})
	writeChan := make(chan string)
	defer close(writeChan)

	go func() {
		for line := range writeChan {
			_, err := f.WriteString(line)
			if err != nil {
				panic(err)
			}
		}
	}()

	chunk := (fn / grt) + 1
	for i := 0; i <= fn; i += chunk {
		step := i + chunk - 1
		if step > min(fn) {
			step = fn
		}
		fmt.Printf("running %d %d\n", i, step)
		go calc(i, step, writeChan, done)
	}

	for i := range grt {
		<-done
		fmt.Printf("done #%d\n", i)
	}
	fmt.Println("exit")
}
