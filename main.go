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
	cfg, err := readArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("read arguments: %w", err))
		os.Exit(1)
	}

	var file *os.File

	if cfg.fname == "" || cfg.fname == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(cfg.fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("open file %s: %w", cfg.fname, err))
			os.Exit(1)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatalf("close file descriptor: %s", err.Error())
			}
		}()
	}

	buffer, err := cut(cfg, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("cut data: %w", err))
		os.Exit(1)
	}

	_, err = fmt.Fprint(os.Stdout, buffer.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("print result to stdout: %w", err))
		os.Exit(1)
	}
}

func readArgs() (*Config, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("you should provide arguments")
	}

	cfg := Config{
		fname:     "",
		fieldIdx:  nil,
		delimiter: "	", // tab by default
	}

	for _, arg := range args {
		if fieldNoStr, ok := strings.CutPrefix(arg, "-f"); ok {
			idxList := []string{}
			if len(fieldNoStr) == 1 {
				idxList = append(idxList, fieldNoStr)
			} else {
				splits := strings.Split(fieldNoStr, " ")
				if len(splits) == 1 {
					splits = strings.Split(splits[0], ",")
				}
				idxList = append(idxList, splits...)
			}

			for _, idx := range idxList {
				fieldNo, err := strconv.Atoi(idx)
				if err != nil {
					return nil, fmt.Errorf("string conversion of %s to integer", fieldNoStr)
				}

				var fieldIdx int
				if fieldNo == 0 {
					fieldIdx = 0
				} else {
					fieldIdx = fieldNo - 1
				}

				cfg.fieldIdx = append(cfg.fieldIdx, fieldIdx)
			}

			continue
		}

		if delimiter, ok := strings.CutPrefix(arg, "-d"); ok {
			if len(delimiter) == 0 {
				return nil, fmt.Errorf("specify delimiter or don't use -d at all")
			}
			cfg.delimiter = delimiter

			continue
		}

		cfg.fname = arg
	}

	if cfg.fieldIdx == nil {
		return nil, fmt.Errorf("no -f argument was found in args: %v", args)
	}

	return &cfg, nil
}

func cut(cfg *Config, file *os.File) (*strings.Builder, error) {
	s := bufio.NewScanner(file)
	b := strings.Builder{}
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

		_, err := fmt.Fprintf(&b, "%s\n", l.String())
		if err != nil {
			return nil, fmt.Errorf("write to string builder: %w", err)
		}
	}

	if s.Err() != nil {
		return nil, fmt.Errorf("scanning data: %w", s.Err())
	}

	return &b, nil
}
