package snowflake

import (
	"sync"
	"time"
)

const (
	epoch          int64 = 1704067200000
	sequenceBits         = 12
	nodeBits             = 10
	maxSequence    int64 = -1 ^ (-1 << sequenceBits)
	maxNode        int64 = -1 ^ (-1 << nodeBits)
	nodeShift            = sequenceBits
	timeShift            = sequenceBits + nodeBits
)

type Generator struct {
	mu        sync.Mutex
	node      int64
	lastStamp int64
	sequence  int64
}

func New(node int64) *Generator {
	if node < 0 || node > maxNode {
		node = 1
	}
	return &Generator{node: node}
}

func (g *Generator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := currentMillis()
	if now == g.lastStamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			for now <= g.lastStamp {
				now = currentMillis()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastStamp = now
	return ((now - epoch) << timeShift) | (g.node << nodeShift) | g.sequence
}

func currentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
