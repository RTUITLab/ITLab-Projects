package main_test

import (
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestFubc_Bench(t *testing.T) {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/api/v1/projects/issues", nil)
	if err != nil {
		t.Log(err)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(t *testing.T, req *http.Request, wg *sync.WaitGroup) {
			time.Sleep(10 * time.Millisecond)
			defer wg.Done()
			log.Println("Gor start")
			if resp, err := http.DefaultClient.Do(req); err != nil {
				log.Println(err)
			} else {
				log.Println(resp.StatusCode)
			}
		}(t, req, &wg)
	}
	wg.Wait()

	t.Log("DOne")

}
