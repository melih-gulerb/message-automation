package base

import (
	"fmt"
	"time"
)

func Log(message string) {
	fmt.Printf("\n[%s] %s", time.Now().Format("2006-01-02.15.04.05"), message)
}
