package handle

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var lim *limiter
var cfg *config

type config struct {
	rps      uint16
	duration uint16
	method   byte
}

type limiter struct {
	cnt    uint16
	RPS    uint16
	reqs   uint16
	done   chan struct{}
	ticker time.Ticker
	sync.Mutex
}

// initializes default configuration for programm.
func init() {
	cfg = &config{
		rps:      100,
		duration: 5,
		method:   1,
	}

	cfg.staticCheck()

	lim = newLimiter(cfg.rps)
}

// ChangeConfig changes the configuration
// rps: the desired number of requests per second
// duration: the desired duration of requests repetitions
// method: the desired Method of generation of requests (1 - "throw out", 0 - "detain")
// Note: It only updates the desired configuration.
func ChangeConfig(rps, duration uint16, method byte) {
	cfg = &config{
		rps:      rps,      // desired requests per second
		duration: duration, // desired Duration of requests repetitions
		method:   method,   // desired Method of generation of requests where 1 - "throw out" and 0 - "detain"
	}

	cfg.staticCheck()

	lim = newLimiter(cfg.rps)
}

var fatal = log.Fatalf

// checks configuration for valid fields
func (cfg *config) staticCheck() {
	if cfg.rps < 1 || cfg.duration < 1 || cfg.method > 1 {
		fatal("Invalid configuration")
	} else if cfg.rps > 1000 || cfg.duration > 60 {
		fatal("Are you sure? The RPS and Duration is unlikely to exceed 1000 and 60, as it would take too much.")
	}
}

// starts a tiker for whole programm
func tiker() {
	for {
		select {
		case <-lim.ticker.C:
			lim.Lock()
			lim.reqs = lim.RPS
			lim.Unlock()
		case <-lim.done:
			lim.Lock()
			lim.reqs = 0
			lim.ticker.Stop()
			lim.Unlock()
			return
		}
	}
}

// makes new limiter struct with rps from config
func newLimiter(rps uint16) *limiter {
	return &limiter{RPS: rps, reqs: rps, done: make(chan struct{}), ticker: *time.NewTicker(time.Second)}
}

// Makes a single GET request to the specified URL. Always respects the rate limiter.
func Get(url string) {
	select {
	case <-lim.done:
		return
	default:
		lim.Lock()
		if lim.reqs == 0 {
			lim.Unlock()
			return
		} else {
			lim.cnt++
			lim.reqs--
			lim.Unlock()

			client := &http.Client{}
			_, err := client.Get(url)
			if err != nil {
				select {
				case <-lim.done:
				default:
					close(lim.done)
				}
			}
		}
	}
}

// Makes a single POST request to the specified URL with the provided body. Always respects the rate limiter.
func Post(url string, body []byte) {
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

// Recieve desired METHOD, URL and optional body. Makes RPS * duration requests to URL.
// Attack will stop immediately after `duration` seconds
func Attack(method string, url string, body ...[]byte) string {
	T := time.Now()
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(cfg.duration) * time.Second)
		lim.ticker.Stop()
		close(ch)
	}()

	go tiker()
	switch cfg.method {
	case 1:
		switch method {
		case "GET":
			for {
				select {
				case <-ch:
					return fmt.Sprintf("%v requests per %v seconds", lim.cnt, time.Since(T))
				case <-lim.done:
					select {
					case <-ch:
						return fmt.Sprintf("%v requests per %v seconds", lim.cnt, time.Since(T))
					default:
						close(ch)
						return fmt.Sprintf("%v requests per %v seconds", lim.cnt, time.Since(T))
					}
				default:
					if lim.reqs > 0 {
						go Get(url)
					}
				}
			}
		case "POST":
		default:
			return "Invalid method"
		}
	case 0:

	}
	return "Something went wrong"
}

/*
func Get(url string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		select {
		case <-lim.done:
			lim.reqs = 0
			return
		default:
			if lim.reqs == 0 {
				return
			} else {
				lim.Lock()
				lim.cnt++
				lim.reqs--
				lim.Unlock()

				client := &http.Client{}
				_, err := client.Get(url)
				if err != nil {
					select {
					case <-lim.done:
					default:
						close(lim.done)
					}
				}
				wg.Done()
			}
		}
	}()
	wg.Wait()
}*/
