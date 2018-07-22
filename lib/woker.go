package lib

import (
	"sync"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
)

// Basic worker struct that does all the job
// Uses boolean channel as a limiter for a K number of concurrent jobs
type worker struct {
	sum int
	wg sync.WaitGroup
	sync.Mutex
	limiter chan bool
}

func NewWorker(k int) *worker {
	return &worker{limiter: make(chan bool, k)}
}

// Runs concurrent jobs for specified url
func (w *worker) Run(url string) {
	w.limiter <- true
	w.wg.Add(1)
	go func() {
		defer func() {
			<- w.limiter
			w.wg.Done()
		}()

		data, err := w.getData(url)

		if err != nil {
			fmt.Println(err)
			return
		}
		sum := w.getCount(data, "Go")
		fmt.Printf("Count for %s: %d\n", url, sum)
		w.storeSum(sum)
	}()
}

// getData tries to reach the specified URL
// Retrieves responce's body on a success
func (w *worker) getData(s string) ([]byte, error) {
	resp, err := http.Get(s)
	if err != nil {
		return []byte{}, fmt.Errorf("get request failed with an error: %s", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("could not read the response body, error: %s", err)
	}

    return data, nil
}

func (w *worker) getCount(data []byte, s string) int {
	return bytes.Count(data, []byte(s))
}

func (w *worker) storeSum(i int) {
	if i < 0 {
		return
	}
	w.Lock()
	w.sum += i
	w.Unlock()
}

func (w worker) GetSum() int {
	return w.sum
}

func (w *worker) Wait()  {
	w.wg.Wait()
}