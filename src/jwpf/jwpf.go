package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"../fstring"
)

func packet(url string, client *http.Client, req *http.Request, nWorker int) error {

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

func createCookie(cookiesString []string) []http.Cookie {
	if cookiesString == nil {
		return nil
	}
	var cookies []http.Cookie
	for _, data := range cookiesString {
		d := strings.Split(data, "=")
		name, value := d[0], d[1]
		cookies = append(
			cookies,
			http.Cookie{Name: name, Value: value})
	}
	return cookies
}
func addCookie(req *http.Request, cookies []http.Cookie) {
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}
}

func worker(target string, wordlist []string, cookies []string, wg *sync.WaitGroup, nWorker int) {
	defer wg.Done()
	replacer := strings.NewReplacer(" ", "%20")
	cookie := createCookie(cookies)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout: 5 * time.Second},
		Timeout: time.Second * 10}
	for _, word := range wordlist {
		if strings.Contains(word, " ") {
			word = replacer.Replace(word)
		}
		target = appendslash(target)
		req, reqErr := http.NewRequest("GET", target+word, nil)
		if reqErr != nil {
			//NOTE: When url cannot be processed, it will be skipped.
			continue
		}
		if cookie != nil {
			addCookie(req, cookie)
		}
		var packErr error
		for i := 0; i < 30; i++ {
			packErr = packet(target+word, client, req, nWorker)
			if packErr == nil {
				break
			}
			sleep(1)
		}
		if packErr != nil {
			fmt.Printf(" [%d] WORKER:  %s\n", nWorker, packErr.Error())
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(" Missing argument.\n Usage:\t%s <url> <dictionary> <threads> <cookie1> ... <cookie2>\n", os.Args[0])
		os.Exit(0)
	}

	nThreads, atoiErr := strconv.Atoi(os.Args[3])
	if atoiErr != nil {
		fmt.Printf(" Error. Threads argument must be an integer.\n")
		os.Exit(0)
	}
	if nThreads < 1 {
		fmt.Printf(" No joke here. Threads argument must be greater than zero\n")
		os.Exit(0)
	}
	target := os.Args[1]
	filepath := os.Args[2]
	var cookies []string
	if len(os.Args) > 3 {
		cookies = os.Args[4:]
	} else {
		cookies = nil
	}

	//Create big dictionary containing the wordlist
	var dict []string
	// Fill the dictionary
	readFile(filepath, &dict)

	// Split the loaded dictionary in nThreads lists
	bi := fstring.ListDivider(dict, nThreads)

	// Create a WaitGroup to wait until all go routines end.
	var wg sync.WaitGroup
	wg.Add(nThreads)

	// Create nThreads go routines
	for i := 0; i < nThreads; i++ {
		go worker(target, bi[i], cookies, &wg, i)
	}

	wg.Wait()
	fmt.Println("\n Process terminated.")
}
