package events

import (
	"fmt"
	"time"
)

type Fallback struct {
}

func NewEventsFallback() *Fallback {
	return &Fallback{}
}

func (lf *Fallback) Emit(mm Event) error {
	fmt.Printf("%s - [%s] : %s\n", time.Now().Format("2006-01-02 15:04:05"), mm.Key, mm.Data)
	return nil
}
