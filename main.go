package main

import (
	"bufio"
	"fmt"
	"io"
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
				log.Printf("close file descriptor: %s", err.Error())
			}
		}()
	}

	w := bufio.NewWriter(os.Stdout)

	err = cut(w, file, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("cut data: %w", err))
		os.Exit(1)
	}

	err = w.Flush()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("flush stdout wrapper: %w", err))
		os.Exit(1)
	}
}

func readArgs() (*Config, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, fmt.Errorf("you should provide some arguments")
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

				if fieldNo < 1 {
					return nil, fmt.Errorf("fields are numbered from 1")
				}

				fieldIdx := fieldNo - 1

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

func cut(w io.Writer, r io.Reader, cfg *Config) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		splits := strings.Split(line, cfg.delimiter)
		if len(splits) == 1 {
			_, err := fmt.Fprintf(w, "%v\n", splits[0])
			if err != nil {
				return fmt.Errorf("write to stdout: %w", err)
			}

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

		_, err := fmt.Fprintf(w, "%s\n", l.String())
		if err != nil {
			return fmt.Errorf("write to stdout: %w", err)
		}
	}

	if s.Err() != nil {
		return fmt.Errorf("scanning data: %w", s.Err())
	}

	return nil
}
