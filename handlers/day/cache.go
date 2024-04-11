package day

import (
	"goirc/shell"
	"math/rand"
	"strings"
	"time"
)

type cache struct {
	cmd    string
	data   []string
	ts     time.Time
	maxAge time.Duration
}

func NewCache(cmd string, maxAge time.Duration) cache {
	return cache{cmd: cmd, maxAge: maxAge}
}

func (c *cache) Load() error {
	r, err := shell.Command(c.cmd)
	if err != nil {
		return err
	}
	r = strings.TrimSpace(r)
	c.data = strings.Split(r, "\n")
	c.ts = time.Now()
	return nil
}

func (c *cache) Pop() (string, error) {
	if time.Since(c.ts) >= c.maxAge {
		err := c.Load()
		if err != nil {
			return "", err
		}
	}

	if len(c.data) == 0 {
		err := c.Load()
		if err != nil {
			return "", err
		}
		return "EOF", nil
	}

	i := rand.Intn(len(c.data))
	result := c.data[i]

	c.data = append(c.data[0:i], c.data[i+1:]...)

	return result, nil
}
