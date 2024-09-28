package handle

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
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
	RPS       uint16
	reqs      atomic.Int32
	toExecute chan struct{}
	done      chan struct{}
	ticker    time.Ticker
	cnt       atomic.Int32
	AtomicShards
}

type AtomicShards struct {
	shards []AtomicCounter
}

type AtomicCounter struct {
	cnt atomic.Int32
	_   [60]byte
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
	reqs := cfg.rps * cfg.duration
	var l limiter
	l.RPS = rps
	l.reqs.Add(int32(rps))
	l.done = make(chan struct{})
	l.toExecute = make(chan struct{}, rps)
	l.ticker = *time.NewTicker(time.Second)
	for i := 0; i < int(reqs); i++ {
		l.shards = append(l.shards, AtomicCounter{})
	}
	return &l
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
				lim.reqs.Add(int32(lim.RPS))
			case <-lim.done:
				lim.reqs.Store(0)
				return
			}
		}
	}
}

// Recieve desired METHOD, URL and optional body. Makes RPS * duration requests to URL.
// Attack will stop immediately after `duration` seconds or after all requests have sent.
func Attack(method, url string, body ...[]byte) string {
	ch := make(chan struct{})
	var cnt int32
	reqsNum := int(lim.RPS * cfg.duration)
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
				for i := 0; i < reqsNum; i++ {
					go getDetained(i, url)
				}
			}()

			<-ch
			for i := 0; i < reqsNum; i++ {
				cnt += lim.shards[i].cnt.Load()
			}
			return fmt.Sprintf("%v requests per %v", cnt, time.Since(T))
		case "POST":
			go func() {
				for i := 0; i < int(lim.RPS); i++ {
					lim.toExecute <- struct{}{}
				}
			}()

			T = time.Now()
			go ticker()

			go func() {
				for i := 0; i < reqsNum; i++ {
					go postDetained(i, url, body[0])
				}
			}()

			<-ch
			for i := 0; i < reqsNum; i++ {
				cnt += lim.shards[i].cnt.Load()
			}
			return fmt.Sprintf("%v requests per %v", cnt, time.Since(T))
		default:
			return "Invalid method"
		}
	} else {
		switch method {
		case "GET":
			go ticker()
			for {
				select {
				case <-ch:
					return fmt.Sprintf("%v requests per %v", lim.cnt.Load(), time.Since(T))
				case <-lim.done:
					select {
					case <-ch:
						return fmt.Sprintf("%v requests per %v", lim.cnt.Load(), time.Since(T))
					default:
						close(ch)
						return fmt.Sprintf("%v requests per %v", lim.cnt.Load(), time.Since(T))
					}
				default:
					if lim.reqs.Load() > 0 {
						go Get(url)
					}
				}
			}
		case "POST":
			go ticker()
			for {
				select {
				case <-ch:
					return fmt.Sprintf("%v requests per %v", lim.cnt.Load(), time.Since(T))
				case <-lim.done:
					select {
					case <-ch:
						return fmt.Sprintf("%v requests per %v", lim.cnt.Load(), time.Since(T))
					default:
						close(ch)
						return fmt.Sprintf("%v requests per %v", lim.cnt.Load(), time.Since(T))
					}
				default:
					if lim.reqs.Load() > 0 {
						go Post(url, body[0])
					}
				}
			}
		default:
			return "Invalid method"
		}
	}
}

// Makes a single GET request to the specified URL. Always respects the rate limiter.
func Get(url string) {
	select {
	case <-lim.done:
		return
	default:
		if lim.reqs.Load() == 0 {
			return
		} else {
			lim.reqs.Add(-1)
			lim.cnt.Add(1)
			fmt.Println(lim.reqs.Load(), lim.cnt.Load())

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
	select {
	case <-lim.done:
		return
	default:
		if lim.reqs.Load() == 0 {
			return
		} else {
			lim.reqs.Add(-1)
			lim.cnt.Add(1)

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
func getDetained(id int, url string) {
	select {
	case <-lim.done:
		return
	case <-lim.toExecute:
		lim.shards[id].cnt.Add(1)

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
func postDetained(id int, url string, body []byte) {
	select {
	case <-lim.done:
		return
	case <-lim.toExecute:
		lim.shards[id].cnt.Add(1)

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
