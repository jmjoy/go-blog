package session

import (
	"time"
)

type Session struct {
}

// 针对每一个用户的Sesseion个体
type Single struct {
	data  map[string]string // 数据数组
	mtime time.Time         // 修改时间
}

const (
	OUT_TIME = 30 * 60
)

var (
	// 全局Session数组，键为token
	sessions = make(map[string]*Single)
)

func init() {
	go func() {
		// 每1分钟运行一次gc
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			gc()
		}
	}()
}

// upsert
func (this *Session) Set(token, key, value string) {
	sing, ok := sessions[token]
	if !ok {
		sing = new(Single)
		sing.data = make(map[string]string)
		sessions[token] = sing
	}
	sing.data[key] = value
	sing.mtime = time.Now()
}

func (this *Session) Get(token, key string) (value string) {
	sing, ok := sessions[token]
	if !ok {
		return
	}
	value = sing.data[key]
	sing.mtime = time.Now()
	return
}

func (this *Session) Del(token, key string) {
	sing, ok := sessions[token]
	if !ok {
		return
	}
	delete(sessions, key)
	sing.mtime = time.Now()
}

func (this *Session) Drop(token string) {
	delete(sessions, token)
}

func gc() {
	for k := range sessions {
		sing := sessions[k]
		// 看看是不是过期了
		if sing.mtime.Unix()+OUT_TIME < time.Now().Unix() {
			delete(sessions, k)
		}
	}
}
