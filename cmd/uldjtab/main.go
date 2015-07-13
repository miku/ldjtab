package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miku/ldjtab"
	"github.com/mreiferson/go-ujson"
)

func main() {
	key := flag.String("key", "", "key to extract")
	version := flag.Bool("v", false, "prints current program version")

	flag.Parse()

	if *version {
		fmt.Println(ldjtab.Version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		log.Fatal("input file required")
	}

	if *key == "" {
		log.Fatal("key required")
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	var lineno int64
	w := bufio.NewWriter(os.Stdout)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lineno++
		obj, err := ujson.NewFromBytes([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
		// pad of 10 allows files with up to 10B lines
		w.WriteString(fmt.Sprintf("%s\t%010d\n", obj.Get(*key).String(), lineno))
	}
	w.Flush()
}
