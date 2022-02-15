package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 120 * time.Second}
)

func main() {
	args := os.Args
	if len(args) != 4 {
		log.Fatalln(`command should be in the form of ./<tester> <upper threshold of concurrent clients> <number of requests> <replica list>`)
	}

	numClients, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalln(err)
	}

	numRequests, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalln(err)
	}

	replicas := hosts(args[3])
	wg := &sync.WaitGroup{}
	var counter uint64
	startTime := time.Now().UTC()
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go start(i, &counter, numRequests, replicas, wg)
	}

	wg.Wait()
	latency := time.Since(startTime).Milliseconds()

	persist(numClients, numRequests, int(counter), latency)

	fmt.Println()
	fmt.Printf("testing is completed (%d out of %d requests)\n", counter, numRequests*numClients)
	fmt.Printf("total elapsed time: %d ms\n", latency)
}

func hosts(arg string) []string {
	list := strings.Split(arg, `,`)
	if len(list) == 0 {
		log.Fatalln(`replica list is empty`)
	}

	return list
}

func start(id int, countAddr *uint64, numRequests int, replicas []string, wg *sync.WaitGroup) {
	for i := 0; i < numRequests; i++ {
		replica := replicas[id%len(replicas)]
		val := strconv.Itoa(rand.Intn(1000))

		fmt.Printf(`client: %d, replica: %d, value: %s`, id, id%len(replicas), val)
		fmt.Println()

		res, err := httpClient.Post(`http://`+replica+`/replica/request`, `text/plain`, bytes.NewBuffer([]byte(val)))
		if err != nil {
			log.Println(`ERROR: `, err, val)
			break
		}

		if res.StatusCode != http.StatusOK {
			log.Println(`Failed Response with code:`, res.StatusCode, `of client:`, id, `for val:`, val)
			res.Body.Close()
			continue
		}
		res.Body.Close()
		atomic.AddUint64(countAddr, 1)
	}
	wg.Done()
}

func persist(clients, reqs, success int, latency int64) {
	fileName := `results.csv`
	var data [][]string
	var file *os.File
	defer file.Close()

	if fileExists(fileName) {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalln(err)
		}

		r := csv.NewReader(file)
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			data = append(data, record)
		}
	} else {
		data = append(data, []string{`clients`, `requests per client`, `total requests`, `success requests`, `latency in ms`})
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	w := csv.NewWriter(file)

	// appending new data item [#clients, #reqs_per_client, #total_reqs #success_reqs, latency_in_ms]
	data = append(data, []string{strconv.Itoa(clients), strconv.Itoa(reqs), strconv.Itoa(clients * reqs), strconv.Itoa(success), strconv.Itoa(int(latency))})

	err = w.WriteAll(data)
	if err != nil {
		log.Fatalln(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
