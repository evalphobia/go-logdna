package logdna

import (
	"sync"
	"time"
)

// Daemon is a background daemon for sending logs.
// This struct stores logs and sends them to LogDNA once per checkpoint interval.
type Daemon struct {
	flushLogs func([]*logPayload) error

	spoolMu            sync.Mutex
	spool              []*logPayload
	checkpointSize     int
	checkpointInterval time.Duration
	stopSignal         chan struct{}
}

// NewDaemon creates new Daemon.
// size is the number of logs to send LogDNA API in single checkpoint.
// interval is the length of the checkpoint interval.
// fn is function called at each checkpoint to sends logs to the LogDNA API.
func NewDaemon(size int, interval time.Duration, fn func([]*logPayload) error) *Daemon {
	if size < 1 {
		size = 10
	}
	if interval == 0 {
		interval = 1 * time.Second
	}

	return &Daemon{
		spool:              make([]*logPayload, 0, 4096),
		checkpointSize:     size,
		checkpointInterval: interval,
		stopSignal:         make(chan struct{}),
		flushLogs:          fn,
	}
}

// Add adds log data to the daemon.
func (d *Daemon) Add(logs ...*logPayload) {
	d.spoolMu.Lock()
	d.spool = append(d.spool, logs...)
	d.spoolMu.Unlock()
}

// Flush gets logs from the internal spool and execute the flushLogs function.
func (d *Daemon) Flush() {
	d.spoolMu.Lock()
	var logs []*logPayload
	logs, d.spool = shiftLog(d.spool, d.checkpointSize)
	d.spoolMu.Unlock()
	_ = d.flushLogs(logs)
}

// shiftLog retrieves logs.
func shiftLog(slice []*logPayload, size int) (part []*logPayload, all []*logPayload) {
	if len(slice) <= size {
		return slice, slice[:0]
	}
	return slice[:size], slice[size:]
}

// Run sets the timer to flush data once per checkpoint interval as a background daemon in a goroutine.
func (d *Daemon) Run() {
	ticker := time.NewTicker(d.checkpointInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				d.Flush()
			case <-d.stopSignal:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the daemon.
func (d *Daemon) Stop() {
	d.stopSignal <- struct{}{}
}
