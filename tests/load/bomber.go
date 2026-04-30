package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	baseURL    = "http://localhost:8080"
	warmupCobs = 300
)

type sample struct {
	endpoint string
	latency  time.Duration
	status   int
	ok       bool
}

var (
	cobTxIDs   []string
	cobTxIDsMu sync.RWMutex
)

func main() {
	target := baseURL
	if len(os.Args) > 1 {
		target = os.Args[1]
	}

	// connection pool shared across all runs
	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
	}
	client := &http.Client{Timeout: 10 * time.Second, Transport: transport}

	fmt.Printf("╔══════════════════════════════════════════╗\n")
	fmt.Printf("║       PIX PSP Simulator — Load Test      ║\n")
	fmt.Printf("╚══════════════════════════════════════════╝\n\n")
	fmt.Printf("Target: %s\n\n", target)

	// ── Warm-up ───────────────────────────────────────────────────────────────
	fmt.Printf("Warm-up: creating %d cobranças...\n", warmupCobs)
	for i := 0; i < warmupCobs; i++ {
		txid := randomTxID()
		body := map[string]any{
			"chave": "+5511999990000",
			"valor": map[string]any{"original": "1.00"},
		}
		if s := doRequest(client, "PUT", target+"/cob/"+txid, body, ""); s.ok {
			cobTxIDsMu.Lock()
			cobTxIDs = append(cobTxIDs, txid)
			cobTxIDsMu.Unlock()
		}
	}
	fmt.Printf("Warm-up done. %d cobranças available.\n\n", len(cobTxIDs))

	// ── Sweep concurrency levels ──────────────────────────────────────────────
	levels := []int{1, 5, 10, 25, 50, 100, 200}
	type runResult struct {
		concurrency int
		rps         float64
		avgMs       float64
		p50Ms       float64
		p95Ms       float64
		p99Ms       float64
		errPct      float64
	}
	var runs []runResult

	for _, c := range levels {
		rps, avg, p50, p95, p99, errPct := runLevel(client, target, c, 8*time.Second)
		runs = append(runs, runResult{c, rps, avg, p50, p95, p99, errPct})
		fmt.Printf("  concurrency=%-4d  rps=%-8.0f  avg=%-8.2fms  p50=%-8.2fms  p95=%-8.2fms  p99=%-8.2fms  err=%.1f%%\n",
			c, rps, avg, p50, p95, p99, errPct)
	}

	// ── Find best concurrency and run 30s final test ──────────────────────────
	best := runs[0]
	for _, r := range runs {
		if r.rps > best.rps && r.errPct < 5.0 {
			best = r
		}
	}

	fmt.Printf("\nBest concurrency: %d (%.0f req/s, %.1f%% errors)\n", best.concurrency, best.rps, best.errPct)
	fmt.Printf("\nRunning 30s final test at concurrency=%d...\n", best.concurrency)

	var (
		totalReqs   int64
		totalErrors int64
		mu          sync.Mutex
		allSamples  []sample
	)

	results := make(chan sample, best.concurrency*500)
	deadline := time.Now().Add(30 * time.Second)

	var wg sync.WaitGroup
	for i := 0; i < best.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(deadline) {
				s := dispatch(client, target)
				results <- s
				atomic.AddInt64(&totalReqs, 1)
				if !s.ok {
					atomic.AddInt64(&totalErrors, 1)
				}
			}
		}()
	}

	collected := make(chan struct{})
	go func() {
		for s := range results {
			mu.Lock()
			allSamples = append(allSamples, s)
			mu.Unlock()
		}
		close(collected)
	}()

	wg.Wait()
	close(results)
	<-collected

	printFinalReport(30*time.Second, allSamples, int(totalReqs), int(totalErrors))

	// ── Write-focused benchmark ───────────────────────────────────────────────
	const writeConcurrency = 250
	fmt.Printf("\n\nWrite-focused benchmark (%d goroutines, 30s — POST /cob + PUT /cob/{txid} only)...\n", writeConcurrency)

	var (
		wTotalReqs   int64
		wTotalErrors int64
		wMu          sync.Mutex
		wSamples     []sample
	)

	wResults := make(chan sample, writeConcurrency*500)
	wDeadline := time.Now().Add(30 * time.Second)

	var wwg sync.WaitGroup
	for i := 0; i < writeConcurrency; i++ {
		wwg.Add(1)
		go func() {
			defer wwg.Done()
			for time.Now().Before(wDeadline) {
				s := dispatchWrite(client, target)
				wResults <- s
				atomic.AddInt64(&wTotalReqs, 1)
				if !s.ok {
					atomic.AddInt64(&wTotalErrors, 1)
				}
			}
		}()
	}

	wCollected := make(chan struct{})
	go func() {
		for s := range wResults {
			wMu.Lock()
			wSamples = append(wSamples, s)
			wMu.Unlock()
		}
		close(wCollected)
	}()

	wwg.Wait()
	close(wResults)
	<-wCollected

	printFinalReport(30*time.Second, wSamples, int(wTotalReqs), int(wTotalErrors))
}

func runLevel(client *http.Client, base string, concurrency int, dur time.Duration) (rps, avgMs, p50Ms, p95Ms, p99Ms, errPct float64) {
	var (
		totalErrors int64
		mu          sync.Mutex
		latencies   []time.Duration
	)

	results := make(chan sample, concurrency*200)
	deadline := time.Now().Add(dur)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(deadline) {
				s := dispatch(client, base)
				results <- s
				if !s.ok {
					atomic.AddInt64(&totalErrors, 1)
				}
			}
		}()
	}

	collected := make(chan struct{})
	go func() {
		for s := range results {
			mu.Lock()
			latencies = append(latencies, s.latency)
			mu.Unlock()
		}
		close(collected)
	}()

	wg.Wait()
	close(results)
	<-collected

	n := len(latencies)
	if n == 0 {
		return
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	rps = float64(n) / dur.Seconds()
	avgMs = msOf(avgDuration(latencies))
	p50Ms = msOf(percentile(latencies, 50))
	p95Ms = msOf(percentile(latencies, 95))
	p99Ms = msOf(percentile(latencies, 99))
	errPct = 100.0 * float64(totalErrors) / float64(n)
	return
}

func dispatch(client *http.Client, base string) sample {
	r := rand.Intn(100)
	switch {
	case r < 50:
		txid := randomExistingTxID()
		if txid == "" {
			txid = randomTxID()
		}
		return doRequest(client, "GET", base+"/cob/"+txid, nil, "GET /cob/{txid}")
	case r < 65:
		body := map[string]any{
			"chave": "+5511999990000",
			"valor": map[string]any{"original": "10.00"},
		}
		return doRequest(client, "POST", base+"/cob", body, "POST /cob")
	case r < 80:
		txid := randomTxID()
		body := map[string]any{
			"chave": "+5511999990000",
			"valor": map[string]any{"original": "25.00"},
		}
		s := doRequest(client, "PUT", base+"/cob/"+txid, body, "PUT /cob/{txid}")
		if s.ok {
			cobTxIDsMu.Lock()
			cobTxIDs = append(cobTxIDs, txid)
			cobTxIDsMu.Unlock()
		}
		return s
	case r < 90:
		return doRequest(client, "GET", base+"/cob?status=ATIVA", nil, "GET /cob")
	default:
		txid := randomExistingTxID()
		if txid == "" {
			return sample{endpoint: "POST /cob/simulate", ok: false}
		}
		body := map[string]any{"txid": txid, "valor": "1.00", "infopagador": "load-test"}
		return doRequest(client, "POST", base+"/cob/simulate", body, "POST /cob/simulate")
	}
}

func doRequest(client *http.Client, method, url string, body any, endpoint string) sample {
	start := time.Now()
	var req *http.Request
	var err error
	if body != nil {
		data, _ := json.Marshal(body)
		req, err = http.NewRequest(method, url, bytes.NewReader(data))
		if err != nil {
			return sample{endpoint: endpoint, latency: time.Since(start), ok: false}
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return sample{endpoint: endpoint, latency: time.Since(start), ok: false}
		}
	}
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return sample{endpoint: endpoint, latency: latency, ok: false}
	}
	defer resp.Body.Close()
	return sample{endpoint: endpoint, latency: latency, status: resp.StatusCode, ok: resp.StatusCode < 500}
}

func printFinalReport(elapsed time.Duration, samples []sample, total, errors int) {
	rps := float64(total) / elapsed.Seconds()

	byEndpoint := map[string][]time.Duration{}
	for _, s := range samples {
		byEndpoint[s.endpoint] = append(byEndpoint[s.endpoint], s.latency)
	}

	all := make([]time.Duration, 0, len(samples))
	for _, s := range samples {
		all = append(all, s.latency)
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })

	fmt.Printf("\n╔══════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                   FINAL LOAD TEST REPORT (30s)                  ║\n")
	fmt.Printf("╚══════════════════════════════════════════════════════════════════╝\n\n")
	fmt.Printf("  Total requests  : %d\n", total)
	fmt.Printf("  Successful      : %d (%.1f%%)\n", total-errors, 100*float64(total-errors)/float64(total))
	fmt.Printf("  Errors          : %d (%.1f%%)\n", errors, 100*float64(errors)/float64(total))
	fmt.Printf("  Throughput      : %.0f req/s\n\n", rps)

	printTable("OVERALL", all)

	endpoints := make([]string, 0, len(byEndpoint))
	for ep := range byEndpoint {
		endpoints = append(endpoints, ep)
	}
	sort.Strings(endpoints)

	for _, ep := range endpoints {
		lats := byEndpoint[ep]
		sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })
		pct := 100.0 * float64(len(lats)) / float64(total)
		fmt.Printf("\n  ── %-30s  %d reqs (%.0f%%)\n", ep, len(lats), pct)
		printTable("", lats)
	}
}

func printTable(label string, latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}
	avg := avgDuration(latencies)
	min := latencies[0]
	max := latencies[len(latencies)-1]
	p50 := percentile(latencies, 50)
	p90 := percentile(latencies, 90)
	p95 := percentile(latencies, 95)
	p99 := percentile(latencies, 99)
	if label != "" {
		fmt.Printf("  ── %s\n", label)
	}
	fmt.Printf("  %-10s %-10s %-10s %-10s %-10s %-10s %-10s\n", "avg", "min", "max", "p50", "p90", "p95", "p99")
	fmt.Printf("  %-10s %-10s %-10s %-10s %-10s %-10s %-10s\n",
		fmtD(avg), fmtD(min), fmtD(max), fmtD(p50), fmtD(p90), fmtD(p95), fmtD(p99))
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	return sorted[int(float64(len(sorted)-1)*float64(p)/100.0)]
}

func avgDuration(d []time.Duration) time.Duration {
	var total time.Duration
	for _, v := range d {
		total += v
	}
	return total / time.Duration(len(d))
}

func msOf(d time.Duration) float64   { return float64(d.Microseconds()) / 1000.0 }
func fmtD(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
}

func randomTxID() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func randomExistingTxID() string {
	cobTxIDsMu.RLock()
	defer cobTxIDsMu.RUnlock()
	if len(cobTxIDs) == 0 {
		return ""
	}
	return cobTxIDs[rand.Intn(len(cobTxIDs))]
}

func dispatchWrite(client *http.Client, base string) sample {
	if rand.Intn(2) == 0 {
		body := map[string]any{
			"chave": "+5511999990000",
			"valor": map[string]any{"original": "10.00"},
		}
		return doRequest(client, "POST", base+"/cob", body, "POST /cob")
	}
	txid := randomTxID()
	body := map[string]any{
		"chave": "+5511999990000",
		"valor": map[string]any{"original": "25.00"},
	}
	s := doRequest(client, "PUT", base+"/cob/"+txid, body, "PUT /cob/{txid}")
	if s.ok {
		cobTxIDsMu.Lock()
		cobTxIDs = append(cobTxIDs, txid)
		cobTxIDsMu.Unlock()
	}
	return s
}
