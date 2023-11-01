package Fan_out

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func readData(file string) <-chan string {
	f, err := os.Open(file) //opens the file for reading
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan string) //channel declared

	//returns a scanner to read from f
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines) //scanning it line-by-line token

	//loop through the fileScanner based on our token split
	go func() {
		for fileScanner.Scan() {
			val := fileScanner.Text() //returns the recent token
			out <- val                //passed the token value to our channel
		}

		close(out) //closed the channel when all content of file is read

		//closed the file
		err := f.Close()
		if err != nil {
			fmt.Printf("Unable to close an opened file: %v\n", err.Error())
			return
		}
	}()

	return out
}

func fanInMergeData(ch1, ch2 <-chan string) chan string {
	chRes := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)

	//reads from 1st channel
	go func() {
		for val := range ch1 {
			chRes <- val
		}
		wg.Done()
	}()

	//reads from 2nd channel
	go func() {
		for val := range ch2 {
			chRes <- val
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()    //waits till the goroutines are completed and wg marked Done
		close(chRes) //close the result channel
	}()

	return chRes
}

func main() {
	ch1 := readData("text1.txt")
	ch2 := readData("text2.txt")

	//receive data from multiple channels and place it on result channel - FanIn
	chRes := fanInMergeData(ch1, ch2)

	//some logic with the result channel
	for val := range chRes {
		fmt.Println(val)
	}
}
