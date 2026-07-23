package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	fname     string
	fieldIdx  *int
	delimiter string
}

func main() {
	args := os.Args[1:]
	if len(args) <= 1 {
		log.Fatal("you should provide arguments")
	}

	cfg := Config{
		fname:     "",
		fieldIdx:  nil,
		delimiter: "	", // tab by default
	}

	for i, arg := range args {
		if i == len(args)-1 {
			cfg.fname = arg

			continue
		}

		if fieldNoStr, ok := strings.CutPrefix(arg, "-f"); ok {
			fieldNo, err := strconv.Atoi(fieldNoStr)
			if err != nil {
				log.Fatalf("string conversion of %s to integer", fieldNoStr)
			}
			fieldIdx := fieldNo - 1
			cfg.fieldIdx = &fieldIdx

			continue
		}

		if delimiter, ok := strings.CutPrefix(arg, "-d"); ok {
			cfg.delimiter = delimiter

			continue
		}
	}

	if cfg.fieldIdx == nil {
		log.Fatalf("no -f argument was found in args: %v", args)
	}

	fname := args[len(args)-1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("open file %s: %s", fname, err.Error())
	}

	s := bufio.NewScanner(file)

	for s.Scan() {
		line := s.Text()
		splits := strings.Split(line, cfg.delimiter)
		if len(splits) == 1 {
			fmt.Printf("%v\n", splits[0])

			continue
		}

		if len(splits) < *cfg.fieldIdx+1 {
			fmt.Println("")

			continue
		}

		fmt.Printf("%v\n", splits[*cfg.fieldIdx])
	}
}
