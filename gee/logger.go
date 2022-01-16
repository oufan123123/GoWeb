package gee

import (
	"log"
	"time"
)

func Logger() Handler {
	return func(c *Context) {
		// start conduct time
		t := time.Now()
		c.Next()
		// calculate time consume
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))

	}
}
