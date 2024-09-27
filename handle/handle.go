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
	detain   bool
}

type limiter struct {
	cnt       uint16
	RPS       uint16
	reqs      uint16
	toExecute chan struct{}
	done      chan struct{}
	ticker    time.Ticker
	sync.Mutex
}

// initializes default configuration for programm.
func init() {
	cfg = &config{
		rps:      100,
		duration: 5,
		detain:   false,
	}

	cfg.staticCheck()

	lim = newLimiter(cfg.rps)
}

// ChangeConfig changes the configuration.
// rps: the desired number of requests per second.
// duration: the desired duration of requests repetitions.
// method: the desired Method of generation of detained requests where.
// Note: It only updates the desired configuration.
func ChangeConfig(rps, duration uint16, method bool) {
	cfg = &config{
		rps:      rps,      // desired requests per second
		duration: duration, // desired Duration of requests repetitions
		detain:   method,   // desired Method of generation of requests where 1 - "throw out" and 0 - "detain"
	}

	cfg.staticCheck()

	lim = newLimiter(cfg.rps)
}

var fatal = log.Fatalf

// checks configuration for valid fields
func (cfg *config) staticCheck() {
	if cfg.rps < 1 || cfg.duration < 1 {
		fatal("Invalid configuration")
	} else if cfg.rps > 1000 || cfg.duration > 60 {
		fatal("Are you sure? The RPS and Duration is unlikely to exceed 1000 and 60, as it would take too much.")
	}
}

// makes new limiter struct with rps from config
func newLimiter(rps uint16) *limiter {
	return &limiter{
		RPS:       rps,
		reqs:      rps,
		done:      make(chan struct{}),
		toExecute: make(chan struct{}, rps),
		ticker:    *time.NewTicker(time.Second)}
}

// starts a global ticker handler
func ticker() {
	for {
		switch cfg.detain {
		case true:
			select {
			case <-lim.ticker.C:
				go func() {
					for i := 0; i < int(lim.RPS); i++ {
						lim.toExecute <- struct{}{}
					}
				}()
			case <-lim.done:
				return
			}
		case false:
			select {
			case <-lim.ticker.C:
				lim.Lock()
				lim.reqs = lim.RPS
				lim.Unlock()
			case <-lim.done:
				lim.Lock()
				lim.reqs = 0
				lim.Unlock()
				return
			}
		}
	}
}

// Recieve desired METHOD, URL and optional body. Makes RPS * duration requests to URL.
// Attack will stop immediately after `duration` seconds or after all requests have sent.
func Attack(method string, url string, body ...[]byte) string {
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(cfg.duration) * time.Second)
		close(ch)
	}()

	T := time.Now()
	if cfg.detain {
		switch method {
		case "GET":
			go func() {
				for i := 0; i < int(lim.RPS); i++ {
					lim.toExecute <- struct{}{}
				}
			}()

			T = time.Now()
			go ticker()

			go func() {
				for i := 0; i < int(lim.RPS*cfg.duration); i++ {
					go getDetained(url)
				}
			}()

			<-ch
			return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
		case "POST":
			go func() {
				for i := 0; i < int(lim.RPS); i++ {
					lim.toExecute <- struct{}{}
				}
			}()

			T = time.Now()
			go ticker()

			go func() {
				for i := 0; i < int(lim.RPS*cfg.duration); i++ {
					go postDetained(url, body[0])
				}
			}()

			<-ch
			return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
		}
	} else {
		switch method {
		case "GET":
			go ticker()
			for {
				select {
				case <-ch:
					return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
				case <-lim.done:
					select {
					case <-ch:
						return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
					default:
						close(ch)
						return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
					}
				default:
					if lim.reqs > 0 {
						go Get(url)
					}
				}
			}
		case "POST":
			go ticker()
			for {
				select {
				case <-ch:
					return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
				case <-lim.done:
					select {
					case <-ch:
						return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
					default:
						close(ch)
						return fmt.Sprintf("%v requests per %v", lim.cnt, time.Since(T))
					}
				default:
					if lim.reqs > 0 {
						go Post(url, body[0])
					}
				}
			}
		default:
			return "Invalid method"
		}
	}
	return "Something went wrong"
}

// Makes a single GET request to the specified URL. Always respects the rate limiter.
func Get(url string) {
	lim.Lock()
	select {
	case <-lim.done:
		return
	default:
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
	lim.Lock()
	select {
	case <-lim.done:
		return
	default:
		if lim.reqs == 0 {
			lim.Unlock()
			return
		} else {
			lim.cnt++
			lim.reqs--
			lim.Unlock()

			client := &http.Client{}

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			if err != nil {
				select {
				case <-lim.done:
				default:
					close(lim.done)
				}
			}

			_, err = client.Do(req)
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

// Makes a single GET request to the specified URL. Always respects the rate limiter.
func getDetained(url string) {
	select {
	case <-lim.done:
		return
	case <-lim.toExecute:
		lim.Lock()
		lim.cnt++
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
func postDetained(url string, body []byte) {
	select {
	case <-lim.done:
		return
	case <-lim.toExecute:
		lim.Lock()
		lim.cnt++
		lim.Unlock()

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			select {
			case <-lim.done:
			default:
				close(lim.done)
			}
		}

		_, err = client.Do(req)
		if err != nil {
			select {
			case <-lim.done:
			default:
				close(lim.done)
			}
		}
	}
}
