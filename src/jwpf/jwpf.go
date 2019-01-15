package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"../fstring"
)

func packet(url string, client *http.Client, req *http.Request, n_worker int) error {

	reqRes, reqErr := client.Do(req)
	if reqErr != nil {
		return reqErr
	}
	defer reqRes.Body.Close()

	statusCode := reqRes.StatusCode
	logString := fstring.FormatLog(url, statusCode)
	if statusCode >= 200 && statusCode <= 226 {
		fmt.Printf("%s", fstring.GREEN(logString))
	} else if statusCode == 404 {
		//pass
	} else if statusCode >= 400 && statusCode <= 451 {
		fmt.Printf("%s", fstring.RED(logString))
	} else if statusCode >= 300 && statusCode <= 308 {
		fmt.Printf("%s", fstring.BLUE(logString))
	} else {
		fmt.Printf("%s", logString)
	}
	return nil
}

func sleep(sec int) {
	time.Sleep(time.Second * time.Duration(sec))
}

func readFile(filename string, dict *[]string) {
	file, errOpen := os.Open(filename)
	if errOpen != nil {
		fmt.Printf(" Cannot read file: %s", errOpen.Error())
		os.Exit(1)
	}
	defer file.Close()
	fmt.Printf("File opened. Filling channel...\n")
	s := bufio.NewScanner(file)
	for s.Scan() {
		*dict = append(*dict, s.Text())
	}
}
func appendslash(text string) string {
	// This adds a "/" at the end of the string
	// if there isn't one yet
	l := len(text) - 1
	if text[l] != '/' {
		return text + "/"
	}
	return text
}

func worker(target string, wordlist []string, flag *bool, n_worker int) {
	for _, word := range wordlist {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				TLSHandshakeTimeout: 5 * time.Second},
			Timeout: time.Second * 10}
		target = appendslash(target)
		req, _ := http.NewRequest("GET", target+word, nil)
		var reqError error
		for i := 0; i < 5; i++ {
			reqError = packet(target+word, client, req, n_worker)
			if reqError == nil {
				break
			}
		}
		if reqError != nil {
			fmt.Printf(" [%d] WORKER:  %s\n", n_worker, reqError.Error())
		}
	}
	*flag = false
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(" Missing argument.\n Usage:\t%s <url> <dictionary> <threads>\n", os.Args[0])
		os.Exit(0)
	}

	n_threads, atoiErr := strconv.Atoi(os.Args[3])
	if atoiErr != nil {
		fmt.Printf(" Error. Threads argument must be an integer.\n")
		os.Exit(0)
	}
	if n_threads < 1 {
		fmt.Printf(" No joke here. Threads argument must be greater than zero\n")
		os.Exit(0)
	}
	target := os.Args[1]
	filepath := os.Args[2]

	//Create big dictionary containing the wordlist
	var dict []string
	// Fill the dictionary
	readFile(filepath, &dict)

	// Split the loaded dictionary in n_threads lists
	bi := fstring.ListDivider(dict, n_threads)

	// Create the flags that will check the status of threads
	flags := make([]bool, n_threads)
	// Create n_threads go routines
	for c, _ := range flags {
		flags[c] = true
		go worker(target, bi[c], &flags[c], c)
	}

	for {
		// Check every second if all thread has finished.
		for _, flag := range flags {
			if flag {
				break
			}
			fmt.Printf("\n\tProcess terminated.\n")
			os.Exit(0)

		}
		sleep(1)

	}

}
