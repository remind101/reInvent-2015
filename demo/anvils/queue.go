package main

import (
	"encoding/json"
	"net/url"
	"time"

	"gopkg.in/redis.v3"
)

type Queue struct {
	*redis.Client
}

func newQueue(uri string) *Queue {
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	var pass string
	if u.User != nil {
		pass, _ = u.User.Password()
	}
	c := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: pass,
	})
	return &Queue{c}
}

func (q *Queue) Push(r DropRequest) {
	raw, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	cmd := q.RPush(q.queueName(), string(raw))
	if err := cmd.Err(); err != nil {
		panic(err)
	}
}

func (q *Queue) Pop() DropRequest {
	var r DropRequest

	cmd := q.BLPop(100*time.Hour, q.queueName())
	if err := cmd.Err(); err != nil {
		panic(err)
	}

	result, err := cmd.Result()
	if err != nil {
		panic(err)
	}

	raw := []byte(result[1])
	if err := json.Unmarshal(raw, &r); err != nil {
		panic(err)
	}

	return r
}

func (q *Queue) queueName() string {
	return "queue"
}
