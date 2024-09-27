package handle

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var lim *Limiter

func init() {
	// Adjust this value to change the rate limiter's rate
	var rate uint16 = 20

	lim = NewLimiter(rate)

	go func() {
		for {
			select {
			case <-lim.ticker.C:
				lim.Lock()
				lim.reqs = lim.RPS
				lim.Unlock()
			case <-lim.ticker.C:
				lim.Lock()
				lim.reqs = lim.RPS
				lim.Unlock()
			}
		}
	}()
}

type Limiter struct {
	cnt    uint16
	RPS    uint16
	reqs   uint16
	ticker time.Ticker
	sync.Mutex
}

func NewLimiter(rps uint16) *Limiter {
	return &Limiter{RPS: rps, reqs: rps, ticker: *time.NewTicker(time.Second)}
}

type Req []byte

func Get(url string) {
	if lim.reqs > 0 {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			if lim.reqs == 0 {
				return
			}
			lim.Lock()
			lim.cnt++
			lim.reqs--
			lim.Unlock()

			client := &http.Client{}
			_, err := client.Get(url)
			if err != nil {
				panic(err.Error())
			}
			// fmt.Println("done", lim.cnt)
			wg.Done()
		}()

		wg.Wait()
	}
}

func Attack(method string, url string, t int, body ...[]byte) string {
	ch := make(chan struct{})
	T := time.Now()

	go func() {
		time.Sleep(time.Duration(t) * time.Second)
		close(ch)
	}()

	switch method {
	case "GET":
		for {
			select {
			case <-ch:
				return fmt.Sprintf("%v requests per %v seconds", lim.cnt, time.Since(T))
			default:
				go Get(url)
			}
		}
	case "POST":
		for {
			select {
			case <-ch:
				return fmt.Sprintf("%v requests per %v seconds", lim.cnt, time.Since(T))
			default:
				go Post(url, body[0])
			}
		}
	default:
		return "Invalid method"
	}
}

func Post(url string, body []byte) {
	if lim.reqs > 0 {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			if lim.reqs == 0 {
				return
			}
			lim.Lock()
			lim.cnt++
			lim.reqs--
			lim.Unlock()

			client := &http.Client{}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				panic(err.Error())
			}

			_, err = client.Do(req)
			if err != nil {
				panic(err.Error())
			}

			wg.Done()
		}()

		wg.Wait()
	}
}
