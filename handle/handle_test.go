package handle

import (
	"fmt"
	"sync"
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

func TestAtomicConter(t *testing.T) {
	tests := []struct {
		name     string
		reqs     uint16
		time     uint16
		wantFail bool
	}{
		{
			name:     "normal",
			reqs:     100,
			time:     10,
			wantFail: false,
		},
		{
			name:     "to much",
			reqs:     10000,
			time:     10000,
			wantFail: true,
		},
		{
			name:     "zero",
			reqs:     0,
			time:     0,
			wantFail: true,
		},
		{
			name:     "okay",
			reqs:     300,
			time:     5,
			wantFail: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exited := false
			fatal = func(format string, v ...any) { _ = format; _ = v; exited = true }

			ChangeConfig(tt.reqs, tt.time, true)
			cnt := 0
			wg := sync.WaitGroup{}
			wg.Add(int(tt.reqs * tt.time))
			for i := 0; i < int(tt.reqs*tt.time); i++ {
				lim.shards = append(lim.shards, AtomicCounter{})
				go func(i int) {
					lim.shards[i].cnt.Add(1)
					wg.Done()
				}(i)
			}
			wg.Wait()
			for i := 0; i < int(tt.reqs*tt.time); i++ {
				cnt += int(lim.shards[i].cnt.Load())
			}
			if tt.wantFail == exited {
				t.Skip()
			}
			if cnt != int(tt.reqs*tt.time) {
				fmt.Println(cnt, exited, "-", int(tt.reqs*tt.time), tt.wantFail)
				t.Fail()
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		url      string
		wantStop bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				url:      "http://easydev.club/api/v1/todos",
				wantStop: false,
			},
		},
		{
			name: "stopped",
			args: args{
				url:      "http://easydev.club/api/v1/todos",
				wantStop: true,
			},
		},
		{
			name: "wrong",
			args: args{
				url:      "http://",
				wantStop: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ChangeConfig(1, 1, false)
			if tt.args.wantStop {
				close(lim.done)
			}
			go Get(tt.args.url)

			time.Sleep(time.Second) // ПРОСТО ЧТОБЫ НЕ ДОБАВЛЯТЬ ВЕЙТГРУПП В ФУНКЦИЮ

			if tt.args.wantStop {
				if lim.reqs.Load() == 1 && lim.cnt.Load() == 0 {
					t.Skip()
				}
			}
			if lim.reqs.Load() != 0 && lim.cnt.Load() != 1 {
				t.Fail()
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		url      string
		body     []byte
		wantStop bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				url:      "http://easydev.club/api/v1/todos",
				wantStop: false,
			},
		},
		{
			name: "stopped",
			args: args{
				url:      "http://easydev.club/api/v1/todos",
				wantStop: true,
			},
		},
		{
			name: "wrong",
			args: args{
				url:      "http://",
				wantStop: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ChangeConfig(1, 1, false)
			if tt.args.wantStop {
				close(lim.done)
			}
			go Post(tt.args.url, tt.args.body)

			time.Sleep(time.Second) // ПРОСТО ЧТОБЫ НЕ ДОБАВЛЯТЬ ВЕЙТГРУПП В ФУНКЦИЮ

			if tt.args.wantStop {
				if lim.reqs.Load() == 1 && lim.cnt.Load() == 0 {
					t.Skip()
				}
			}
			if lim.reqs.Load() != 0 && lim.cnt.Load() != 1 {
				t.Fail()
			}
		})
	}
}

func TestAttack(t *testing.T) {
	type args struct {
		method   string
		detained bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normalGet",
			args: args{
				method:   "GET",
				detained: false,
			},
		},
		{
			name: "normalGetDetained",
			args: args{
				method:   "GET",
				detained: true,
			},
		},
		{
			name: "normalPost",
			args: args{
				method:   "POST",
				detained: false,
			},
		},
		{
			name: "normalPostDetained",
			args: args{
				method:   "POST",
				detained: true,
			},
		},
		{
			name: "Invalid method",
			args: args{
				method:   "OPTIONS",
				detained: false,
			},
			want: "Invalid method",
		},
	}
	for _, tt := range tests {
		ChangeConfig(1, 1, tt.args.detained)
		t.Run(tt.name, func(t *testing.T) {
			got := ""
			if tt.args.method == "POST" {
				got = Attack(tt.args.method, "https://easydev.club/api/v1/todos", []byte(`{"title": "Test"}`))
			} else {
				got = Attack(tt.args.method, "https://easydev.club/api/v1/todos")
			}
			if tt.want != "" && got != tt.want {
				t.Errorf("Attack() = %v, want %v", got, tt.want)
			}
		})
	}
}
