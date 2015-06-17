// Extract values and line number from line-delimited JSON.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const Version = "0.1.2"

// options carries the flags around
type options struct {
	key       string
	padlength int
}

// extracted carries around the extracted value together with its line number.
type extracted struct {
	lineno int64
	value  string
}

var (
	errKeyNotFound   = fmt.Errorf("key not found")
	errValueNotFound = fmt.Errorf("value not found")
	errInvalidType   = fmt.Errorf("invalid type")
)

// stringValue returns the value for a given key in dot notation. Work
// resursively for objects, but not lists of objects.
func stringValue(key string, doc map[string]interface{}) (string, error) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return "", errKeyNotFound
	}

	val, ok := doc[keys[0]]
	if !ok {
		return "", errKeyNotFound
	}

	switch val.(type) {
	case string:
		return val.(string), nil
	case map[string]interface{}:
		if len(keys) < 2 {
			return "", errValueNotFound
		}
		return stringValue(keys[1], val.(map[string]interface{}))
	case json.Number:
		return fmt.Sprintf("%s", val), nil
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', 6, 64), nil
	case int:
		return strconv.Itoa(val.(int)), nil
	default:
		return "", errInvalidType
	}

	return "", errValueNotFound
}

// extractor extracts a value for a key, which is given in options.
// TODO(miku): separate parallelism and processing via interface callbacks.
func extractor(queue chan []extracted, results chan extracted, opts options, wg *sync.WaitGroup) {
	defer wg.Done()

	for batch := range queue {
		for _, v := range batch {
			target := make(map[string]interface{})
			d := json.NewDecoder(strings.NewReader(v.value))
			d.UseNumber()

			if err := d.Decode(&target); err != nil {
				log.Fatal(err)
			}
			// drop on simple key access, if its a top level key
			value, err := stringValue(opts.key, target)
			if err != nil {
				log.Fatal(err)
			}
			results <- extracted{value: value, lineno: v.lineno}
		}
	}
}

// leftPad pads string s with padStr
func leftPad(s string, padStr string, pLen int) string {
	r := pLen - len(s)
	if r > 0 {
		return strings.Repeat(padStr, pLen-len(s)) + s
	}
	return s
}

// sink writes values as tab-separated values to stdout. pad line number, so
// it can be sorted alongside non-numeric columns.
func sink(c chan extracted, done chan bool, opts options) {
	w := bufio.NewWriter(os.Stdout)
	for v := range c {
		lineno := leftPad(fmt.Sprintf("%d", v.lineno), "0", opts.padlength)
		w.WriteString(fmt.Sprintf("%s\t%s\n", v.value, lineno))
	}
	w.Flush()
	done <- true
}

func main() {
	key := flag.String("key", "", "key to extract")
	numWorker := flag.Int("w", runtime.NumCPU(), "number of workers")
	batchSize := flag.Int("size", 20000, "size per batch")
	padlength := flag.Int("padlength", 10, "how many zeros as pad for line numbers")
	version := flag.Bool("v", false, "prints current program version")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		log.Fatal("input file required")
	}

	runtime.GOMAXPROCS(*numWorker)

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	var wg sync.WaitGroup

	queue := make(chan []extracted)
	results := make(chan extracted)
	done := make(chan bool)

	for i := 0; i < *numWorker; i++ {
		wg.Add(1)
		go extractor(queue, results, options{key: *key}, &wg)
	}

	go sink(results, done, options{padlength: int(*padlength)})

	var batch []extracted
	var lineno int64

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lineno++

		batch = append(batch, extracted{value: line, lineno: lineno})

		if len(batch)%*batchSize == 0 {
			cc := make([]extracted, len(batch))
			copy(cc, batch)
			queue <- cc
			batch = batch[:0]
		}
	}

	cc := make([]extracted, len(batch))
	copy(cc, batch)
	queue <- cc

	close(queue)
	wg.Wait()
	close(results)
	<-done
}
