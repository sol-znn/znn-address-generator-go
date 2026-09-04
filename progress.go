package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// progressBar is a minimal terminal progress indicator, standing in for the
// Dart `console_bars` FillingBar used by the original tool.
type progressBar struct {
	mu        sync.Mutex
	desc      string
	total     int64
	count     int64
	start     time.Time
	width     int
	showTotal bool
}

func newProgressBar(desc string, total int64, showTotal bool) *progressBar {
	return &progressBar{
		desc:      desc,
		total:     total,
		start:     time.Now(),
		width:     50,
		showTotal: showTotal,
	}
}

func (p *progressBar) increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.count++
	p.render()
}

func (p *progressBar) render() {
	elapsed := time.Since(p.start).Round(time.Second)

	var filled int
	var percent float64
	if p.total > 0 {
		percent = float64(p.count) / float64(p.total)
		if percent > 1 {
			percent = 1
		}
		filled = int(percent * float64(p.width))
	}

	bar := strings.Repeat("=", filled) + strings.Repeat(" ", p.width-filled)

	if p.showTotal {
		fmt.Printf("\r%s [%s] %d/%d (%.1f%%) %s", p.desc, bar, p.count, p.total, percent*100, elapsed)
	} else {
		fmt.Printf("\r%s [%s] %d %s", p.desc, bar, p.count, elapsed)
	}
}
