package gotcha

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/types"
)

type verboseReporter struct {
	levels     map[string][]string
	prefix     bool
	skipped    int
	summarizer summarizer
	mu         *sync.Mutex
}

func newVerboseReporter(prefix bool) *verboseReporter {
	return &verboseReporter{
		levels:     make(map[string][]string),
		prefix:     prefix,
		summarizer: summarizer{},
		mu:         &sync.Mutex{},
	}
}

func (r *verboseReporter) AnnounceSuite(description string) {
	fmt.Println()
	color.New().Add(color.Bold).Println(description)
	fmt.Println()
}

func (r *verboseReporter) PrintSingleSpec(spec *types.SpecSummary, prefix string, fn ColorFunc) {
	size := len(spec.ComponentTexts[1:]) - 1
	fullComponent := ""
	for i, component := range spec.ComponentTexts[1:] {
		fullComponent += "." + component
		level := r.levels[fullComponent]
		found := false
		lower := strings.ToLower(component)
		for j := range level {
			if lower == level[j] {
				found = true
				break
			}
		}
		if !found {
			r.levels[fullComponent] = append(level, lower)
			spaces := strings.Repeat("  ", i+1)
			if i == size {
				if r.prefix {
					fmt.Println(fn(fmt.Sprintf("%s%s%s", spaces, prefix, component)))
				} else {
					fmt.Println(fn(fmt.Sprintf("%s%s", spaces, component)))
				}
			} else {
				fmt.Println(fmt.Sprintf("%s%s", spaces, component))
			}
		}
	}
}

func (r *verboseReporter) PrintSummary(summary *types.SuiteSummary, fn ColorFunc) {
	r.summarizer.printSummary(summary, r.skipped, fn)
}

func (r *verboseReporter) SummarizeFailures(failures []*types.SpecSummary, pendings []*types.SpecSummary) {
	r.summarizer.printFailures(failures, pendings)
}

func (r *verboseReporter) Skip(*types.SpecSummary) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skipped++
}
