package handle

import (
	"fmt"
	"testing"
	"time"
)

func TestReqTime(t *testing.T) {
	ch := make(chan struct{})
	T := time.Now()

	for i := 0; i < int(lim.RPS); i++ {
		fmt.Println(i)
		select {
		case <-ch:
		case <-lim.done:
			select {
			case <-ch:
			default:
				close(ch)
			}
		default:
			go Get("http://easydev.club/api/v1/todos")
		}
	}
	fmt.Println(time.Since(T))
}

func TestConfig_StaitCheck(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config
		exited bool
	}{
		{
			name:   "default",
			cfg:    &config{rps: 100, duration: 5, detain: false},
			exited: false,
		},
		{
			name:   "zero",
			cfg:    &config{rps: 0, duration: 0, detain: false},
			exited: true,
		},
		{
			name:   "normal",
			cfg:    &config{rps: 233, duration: 10, detain: false},
			exited: false,
		},
		{
			name:   "err",
			cfg:    &config{rps: 10000, duration: 10000, detain: false},
			exited: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exited := false
			fatal = func(format string, v ...any) { _ = format; _ = v; exited = true }
			tt.cfg.staticCheck()
			if exited != tt.exited {
				t.Fail()
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Get(tt.args.url)
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		url  string
		body []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Post(tt.args.url, tt.args.body)
		})
	}
}

func TestAttack(t *testing.T) {
	type args struct {
		detain string
		url    string
		body   [][]byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Attack(tt.args.detain, tt.args.url, tt.args.body...); got != tt.want {
				t.Errorf("Attack() = %v, want %v", got, tt.want)
			}
		})
	}
}
