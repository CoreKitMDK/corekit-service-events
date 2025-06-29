package events

import (
	"fmt"
	"time"
)

type Console struct{}

func NewEventsConsole() *Console {
	return &Console{}
}

func (c Console) Emit(m Event) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(fmt.Sprintf("%s - [CONSOLE] [%s] : %s = %s\n", timestamp, m.formatTags(), m.Key, m.Data))
	return nil
}
