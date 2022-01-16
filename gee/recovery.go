package gee

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

func Trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}

	return str.String()

}

// trace error and continue
func Recovery() Handler {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", Trace(message))
				c.Fail(500, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
