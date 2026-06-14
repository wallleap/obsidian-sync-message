package util

import (
	"strconv"
	"sync"
	"time"
)

const (
	workerIDBits  = 10
	sequenceBits  = 12
	workerIDShift = sequenceBits
	timestampShift = sequenceBits + workerIDBits
	maxWorkerID   = -1 ^ (-1 << workerIDBits)
	maxSequence   = -1 ^ (-1 << sequenceBits)
	epoch         = 1609459200000 // 2021-01-01 00:00:00 UTC
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorkerID
	}
	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now < s.timestamp {
		panic("Clock moved backwards")
	}

	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	return (now-epoch)<<timestampShift |
		s.workerID<<workerIDShift |
		s.sequence
}

func (s *Snowflake) GenerateString() string {
	return strconv.FormatInt(s.Generate(), 10)
}

var (
	ErrInvalidWorkerID = &snowflakeError{"invalid worker ID"}
)

type snowflakeError struct {
	msg string
}

func (e *snowflakeError) Error() string {
	return e.msg
}
