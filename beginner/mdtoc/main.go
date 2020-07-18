package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	// open STDIN or open file from path argument
	var input io.Reader
	if len(os.Args) > 1 {
		fh, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fh.Close()
		input = fh
	} else {
		input = os.Stdin
	}

	// iterate lines in input, for each `# headline` output
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		level, title := extractHeadline(line)
		if level == 0 {
			continue
		}
		link := titleAsLink(title)
		fmt.Printf("%s- [%s](#%s)\n", strings.Repeat(" ", level-1), title, link)
	}
}

// extractHeadline returns the level (1..n) and title (text) extracted from the
// given line. If the line is not a headline then level=0 will be returned
func extractHeadline(line string) (level int, title string) {
	for _, c := range line {
		if c != '#' {
			break
		}
		level++
	}
	if level == 0 {
		return
	}
	title = strings.TrimSpace(line[level:])
	return

}

// titleAsLink returns
func titleAsLink(title string) string {
	title = strings.TrimSpace(strings.ToLower(title))
	link := ""
	dashes := 0
	for _, c := range title {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			link += string(c)
			dashes = 0
		} else if dashes == 0 {
			link += "-"
			dashes++
		}
	}
	return strings.Trim(link, "-")
}
