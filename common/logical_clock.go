package common

import (
	"fmt"
	"sync/atomic"
)

type LogicalTimestamp struct {
	value uint64
}

var logicalTimestamp = LogicalTimestamp{value: 0}

func PeekLogicalTimestamp() LogicalTimestamp {
	return logicalTimestamp
}

func AdvanceLogicalTimestamp() LogicalTimestamp {
	return LogicalTimestamp{atomic.AddUint64(&(logicalTimestamp.value), 1)}
}

func (l *LogicalTimestamp) String() string {
	return fmt.Sprintf("[%10d]", l.value)
}
