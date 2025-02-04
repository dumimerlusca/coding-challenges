package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Result struct {
	bytesCount *int
	linesCount *int
	wordCount  *int
	charCount  *int
}

func main() {
	var bytesCount, linesCount, wordCount, charCount bool

	flag.BoolVar(&bytesCount, "c", false, "number of bytes in the file")
	flag.BoolVar(&linesCount, "l", false, "number of lines in the file")
	flag.BoolVar(&wordCount, "w", false, "number of words in the file")
	flag.BoolVar(&charCount, "m", false, "number of characters in the file")
	flag.Parse()
	args := flag.Args()

	if !bytesCount && !linesCount && !wordCount && !charCount {
		bytesCount = true
		linesCount = true
		wordCount = true
		charCount = true
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(fmt.Errorf("failed to read os.Stdin.Stat %w", err))
	}

	useStdin := (info.Mode() & os.ModeCharDevice) == 0
	var input []byte
	var filename string

	if useStdin {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(fmt.Errorf("failed to read from os.Stdin %w", err))
		}
	} else {
		if len(args) != 1 {
			panic(fmt.Sprintf("expected one filename as argument, received: %s", strings.Join(args, ",")))
		}

		filename = args[0]

		input, err = os.ReadFile(filename)

		if err != nil {
			panic(fmt.Errorf("failed to read from file"))
		}
	}

	result := wc(wcArgs{bytesCount: bytesCount, linesCount: linesCount, wordCount: wordCount, charCount: charCount}, input)

	stringRes := []string{}

	if result.linesCount != nil {
		stringRes = append(stringRes, strconv.FormatInt(int64(*result.linesCount), 10))
	}
	if result.wordCount != nil {
		stringRes = append(stringRes, strconv.FormatInt(int64(*result.wordCount), 10))
	}
	if result.charCount != nil {
		stringRes = append(stringRes, strconv.FormatInt(int64(*result.charCount), 10))
	}
	if result.bytesCount != nil {
		stringRes = append(stringRes, strconv.FormatInt(int64(*result.bytesCount), 10))
	}

	if !useStdin {
		stringRes = append(stringRes, filename)
	}

	fmt.Println(strings.Join(stringRes, "    "))
}

type wcArgs struct {
	bytesCount bool
	linesCount bool
	wordCount  bool
	charCount  bool
}

func wc(args wcArgs, input []byte) *Result {
	result := new(Result)

	if args.bytesCount {
		count := len(input)
		result.bytesCount = &count
	}

	if args.linesCount {
		count := bytes.Count(input, []byte("\n"))
		result.linesCount = &count
	}

	if args.wordCount {
		count := countWords(input)
		result.wordCount = &count
	}

	if args.charCount {
		count := countChars(input)
		result.charCount = &count
	}

	return result
}

func countWords(data []byte) int {
	var count int
	inWord := false

	for _, b := range data {
		if unicode.IsSpace(rune(b)) {
			inWord = false
		} else if !inWord {
			inWord = true
			count++
		}
	}

	return count
}

func countChars(data []byte) int {
	count := 0
	for len(data) > 0 {
		_, size := utf8.DecodeRune(data)
		data = data[size:]
		count++
	}
	return count
}
