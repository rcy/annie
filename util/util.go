package util

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"strings"
	"time"
)

func ParseTime(str string) (time.Time, error) {
	result, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		result, err = time.Parse("2006-01-02T15:04:05Z", str)
	}
	return result, err
}

func Since(tstr string) string {
	t, err := ParseTime(tstr)
	if err != nil {
		log.Fatal(err)
	}
	return Ago(time.Now().Sub(t).Round(time.Second))
}

func Ago(d time.Duration) string {
	if d.Hours() >= 48.0 {
		return fmt.Sprintf("%dd", int(math.Round(d.Hours()/24)))
	} else {
		return d.String()
	}
}

// from a uri like https://www.google.com/abc?def=123 return google.com
func BareDomain(uri string) string {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		// just punt and return the original uri
		return uri
	}
	return strings.Replace(parsedUrl.Host, "www.", "", 1)
}

func Getenv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s not set!", key)
	} else {
		log.Printf("%s=%s\n", key, val)
	}

	return val
}
