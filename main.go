package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) <= 1 {
		log.Fatal("you should provide arguments")
	}

	fieldNoStr, ok := strings.CutPrefix(args[0], "-f")
	if !ok {
		log.Fatalf("no -f argument was found in %s", args[0])
	}
	fieldNo, err := strconv.Atoi(fieldNoStr)
	if err != nil {
		log.Fatalf("string conversion of %s to integer", fieldNoStr)
	}

	fname := args[len(args)-1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("open file %s: %s", fname, err.Error())
	}

	s := bufio.NewScanner(file)

	for s.Scan() {
		line := s.Text()
		splits := strings.Split(line, "	")
		fmt.Printf("%v\n", splits[fieldNo-1])
		// fmt.Printf("%s\n", line)
	}
}
