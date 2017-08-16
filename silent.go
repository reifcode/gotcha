package gotcha

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/types"
)

type silentReporter struct {
	summarizer summarizer
}

func newSilentReporter() *silentReporter {
	return &silentReporter{
		summarizer: summarizer{},
	}
}

func (r *silentReporter) AnnounceSuite(description string) {
	fmt.Println()
	color.New().Add(color.Bold).Println(description)
	fmt.Println()
}

func (r *silentReporter) PrintSingleSpec(spec *types.SpecSummary, prefix string, fn ColorFunc) {
	fmt.Print(fn("."))
}

func (r *silentReporter) PrintSummary(spec *types.SuiteSummary, fn ColorFunc) {
	r.summarizer.printSummary(spec, fn)
}

func (r *silentReporter) SummarizeFailures(failures []*types.SpecSummary, pendings []*types.SpecSummary) {
	r.summarizer.printFailures(failures, pendings)
}
