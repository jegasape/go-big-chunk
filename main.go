package main

import (
	"log"
	"math/rand"
	"os"
)

func randtext(length int) string {
	const charset string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func calc(s, e int, writerChan chan string, done chan struct{}) {
	for i := s; i <= e; i++ {
		line := randtext(6) + "\n"
		writerChan <- line
	}
	done <- struct{}{}
}

func main() {
	f, err := os.Create("f.out")
	if err != nil {
		log.Printf("err creating file %v", f)
	}
	defer f.Close()

	grt := 10
	final := 10_000_000
	done := make(chan struct{})
	writeChan := make(chan string)

	go func() {
		for line := range writeChan {
			_, err := f.WriteString(line)
			if err != nil {
				log.Printf("err writting in file %v", err.Error())
			}
		}
	}()

	chunk := (final / grt) + 1
	for i := 0; i <= final; i += chunk {
		step := i + chunk - 1
		if step > min(final) {
			step = final
		}
		log.Printf("running %d %d\n", i, step)
		go calc(i, step, writeChan, done)
	}

	for i := range grt {
		<-done
		log.Printf("done #%d\n", i)
	}
	print("exit\n")
}
