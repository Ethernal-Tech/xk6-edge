package ethereum

import (
	"sync"
)

var (
	selected = -1
	once     sync.Once
)

func selection(VUID int) {
	once.Do(func() {
		selected = VUID
	})

	if selected != VUID {
		return
	}

	go polling(VUID)
}

func polling(VUID int) {
	for {
		// polling logic
	}
}
