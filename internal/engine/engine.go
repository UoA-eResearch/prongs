package engine

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thomaslaurenson/prongs/internal/scanner"
)

// Run executes all active scanners against all hosts concurrently.
// This replaces the copy-pasted ThreadPoolExecutor block in every Python scanner.
func Run(scanners []scanner.Scanner, hosts []net.IP, concurrency int, prettyPrint bool) {
	type job struct {
		ip      net.IP
		scanner scanner.Scanner
	}

	total := int64(len(hosts) * len(scanners))
	jobs := make(chan job, total)
	results := make(chan scanner.Result, 100)
	var processed atomic.Int64

	// Fill the work queue upfront
	for _, ip := range hosts {
		for _, s := range scanners {
			jobs <- job{ip, s}
		}
	}
	close(jobs)

	// Worker pool
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if r, found := j.scanner.Run(j.ip); found {
					results <- r
				}
				processed.Add(1)
			}
		}()
	}

	// Close results channel once all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Progress ticker (pretty mode only).
	// stopTicker signals the goroutine to exit; tickerStopped confirms it has.
	var stopTicker chan struct{}
	var tickerStopped chan struct{}
	if prettyPrint {
		ticker := time.NewTicker(500 * time.Millisecond)
		stopTicker = make(chan struct{})
		tickerStopped = make(chan struct{})
		go func() {
			defer close(tickerStopped)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					fmt.Printf("\rProgress: %d/%d", processed.Load(), total)
				case <-stopTicker:
					return
				}
			}
		}()
	}

	// Print results as they arrive - no buffering, no post-scan loop
	for r := range results {
		printResult(r, prettyPrint)
	}

	if prettyPrint {
		close(stopTicker)
		<-tickerStopped
		fmt.Printf("\rProgress: %d/%d\n", processed.Load(), total)
		fmt.Printf("Total hosts/checks: %d/%d\n", len(hosts), processed.Load())
	}
}

func printResult(r scanner.Result, prettyPrint bool) {
	if prettyPrint {
		fmt.Printf("🚨 %s:%d - %s\n", r.IP, r.Port, r.ScanType)
	} else {
		// Preserve exact Python output format: timestamp\tip\tscanner\tport
		fmt.Printf("%s\t%s\t%s\t%d\n",
			r.Timestamp.UTC().Format(time.RFC3339),
			r.IP,
			r.ScanType,
			r.Port,
		)
	}
}
