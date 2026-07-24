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
	fieldIdx  []int
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
			idxList := []string{}
			splits := strings.Split(fieldNoStr, `"`)
			if len(splits) == 1 {
				if len(splits[0]) == 1 {
					idxList = append(idxList, splits[0])
				} else {
					splits = strings.Split(splits[0], " ")
					if len(splits) == 1 {
						splits = strings.Split(splits[0], ",")
					}
					idxList = append(idxList, splits...)
				}
			}

			for _, idx := range idxList {
				fieldNo, err := strconv.Atoi(idx)
				if err != nil {
					log.Fatalf("string conversion of %s to integer", fieldNoStr)
				}
				fieldIdx := fieldNo - 1
				cfg.fieldIdx = append(cfg.fieldIdx, fieldIdx)
			}

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

		l := strings.Builder{}

		for i, idx := range cfg.fieldIdx {
			if len(splits) < idx+1 {
				continue
			}
			if i == len(cfg.fieldIdx)-1 {
				l.WriteString(splits[idx])
			} else {
				l.WriteString(splits[idx] + cfg.delimiter)
			}
		}

		fmt.Printf("%s\n", l.String())
	}
}
