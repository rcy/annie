package day

import (
	"goirc/shell"
	"math/rand"
	"strings"
	"time"
)

type cache struct {
	cmd      string
	data     []string
	ts       time.Time
	location *time.Location
	day      string
}

func NewCache(cmd string) cache {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	return cache{cmd: cmd, location: location}
}

func (c *cache) Load() error {
	r, err := shell.Command(c.cmd)
	if err != nil {
		return err
	}
	r = strings.TrimSpace(r)
	c.data = strings.Split(r, "\n")
	c.ts = time.Now()
	c.day = c.currentDay()

	return nil
}

func (c *cache) currentDay() string {
	return time.Now().In(c.location).Format("2006-01-02")
}

func (c *cache) Pop() (string, error) {
	if c.day != c.currentDay() {
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
