package main

import (
	"bytes"
	"fmt"
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
	fmt.Println()
	fmt.Printf("testing is completed (%d out of %d requests)\n", counter, numRequests*numClients)
	fmt.Printf("total elapsed time: %d ms\n", time.Since(startTime).Milliseconds())
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
