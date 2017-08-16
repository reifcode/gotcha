package gotcha

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/types"
)

type verboseReporter struct {
	levels     map[int][]string
	prefix     bool
	summarizer summarizer
}

func newVerboseReporter(prefix bool) *verboseReporter {
	return &verboseReporter{
		levels:     make(map[int][]string),
		prefix:     prefix,
		summarizer: summarizer{},
	}
}

func (r verboseReporter) AnnounceSuite(description string) {
	fmt.Println()
	color.New().Add(color.Bold).Println(description)
	fmt.Println()
}

func (r verboseReporter) PrintSingleSpec(spec *types.SpecSummary, prefix string, fn ColorFunc) {
	size := len(spec.ComponentTexts[1:]) - 1
	for i, component := range spec.ComponentTexts[1:] {
		level := r.levels[i]
		found := false
		for _, c := range level {
			if component == c {
				found = true
				break
			}
		}
		if i > len(level) || !found {
			r.levels[i] = append(level, component)
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

func (r verboseReporter) PrintSummary(summary *types.SuiteSummary, fn ColorFunc) {
	r.summarizer.printSummary(summary, fn)
}

func (r verboseReporter) SummarizeFailures(failures []*types.SpecSummary, pendings []*types.SpecSummary) {
	r.summarizer.printFailures(failures, pendings)
}
