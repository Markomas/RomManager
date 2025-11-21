package api

import "time"

type writeCounter struct {
	total          int64
	written        int64
	progressFn     func(progress float64)
	lastReportTime time.Time
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.written += int64(n)
	if wc.total > 0 {
		now := time.Now()
		// Report progress every 100ms or when download is complete
		if now.Sub(wc.lastReportTime) > time.Second || wc.written == wc.total {
			wc.progressFn(float64(wc.written) / float64(wc.total))
			wc.lastReportTime = now
		}
	}
	return n, nil
}
