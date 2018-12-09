package fstring

import (
	"fmt"
	"strings"
	"time"
)

const (
	DATEFORMAT string = "2006-01-02 15:04:05"
)

func RED(text string) string {
	return ("\033[1;40;31m" + text + "\033[0;0m")
}
func GREEN(text string) string {
	return ("\033[1;40;32m" + text + "\033[0;0m")
}
func BLUE(text string) string {
	return ("\033[1;40;94m" + text + "\033[0;0m")
}

func trim(text string) string {
	return strings.TrimSuffix(text, "\n")
}

func ListDivider(list []string, n_lists int) [][]string {
	var bi [][]string

	chunkSize := (len(list) + n_lists - 1) / n_lists

	for i := 0; i < len(list); i += chunkSize {
		end := i + chunkSize

		if end > len(list) {
			end = len(list)
		}

		bi = append(bi, list[i:end])
	}
	return bi
}

func FormatLog(url string, code int) string {
	now := time.Now().Format(DATEFORMAT)
	result := fmt.Sprintf("[%s] %d %s\n", now, code, url)
	return result
}
