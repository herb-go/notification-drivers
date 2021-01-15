package cronqueue

import (
	"time"
)

type Config struct {
	TimeoutDuration  string
	IntervalDuration string
	ExecuteCount     int
}

func (c *Config) CreateQueue() (*Queue, error) {
	var err error
	q := New()
	if c.TimeoutDuration != "" {
		q.Timeout, err = time.ParseDuration(c.TimeoutDuration)
		if err != nil {
			return nil, err
		}
	}
	if c.IntervalDuration != "" {
		q.Interval, err = time.ParseDuration(c.IntervalDuration)
		if err != nil {
			return nil, err
		}
	}
	if c.ExecuteCount != 0 {
		q.ExecuteCount = c.ExecuteCount
	}
	return q, nil
}
