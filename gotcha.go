package gotcha

import (
	"github.com/fatih/color"
	"github.com/reifcode/ginkgo"
	"github.com/reifcode/ginkgo/config"
	"github.com/reifcode/ginkgo/reporters"
	"github.com/reifcode/ginkgo/types"
)

const (
	prefixOK      = "[OK] "
	prefixFail    = "[FAIL] "
	prefixPending = "[PENDING] "
)

func RunSpecs(t ginkgo.GinkgoTestingT, description string) bool {
	reporters := []ginkgo.Reporter{NewReporter()}
	return ginkgo.RunSpecsWithCustomReporters(t, description, reporters)
}

func NewReporter() reporters.Reporter {
	gotcha := NewGotcha(config.DefaultReporterConfig)
	return reporters.NewDefaultReporter(config.DefaultReporterConfig, gotcha)
}

type Reporter interface {
	AnnounceSuite(description string)
	PrintSingleSpec(spec *types.SpecSummary, prefix string, fn ColorFunc)
	PrintSummary(spec *types.SuiteSummary, fn ColorFunc)
	SummarizeFailures(failures []*types.SpecSummary, pendings []*types.SpecSummary)
}

type ColorFunc func(string, ...interface{}) string

type Gotcha struct {
	prefix   bool
	verbose  bool
	reporter Reporter
}

func NewGotcha(cfg config.DefaultReporterConfigType) *Gotcha {
	gotcha := &Gotcha{}
	if cfg.Succinct {
		gotcha.reporter = newSilentReporter()
	} else {
		gotcha.reporter = newVerboseReporter(cfg.NoColor)
	}
	color.NoColor = cfg.NoColor
	return gotcha
}

func (g *Gotcha) AnnounceSuite(description string, randomSeed int64, randomizingAll bool, quiet bool) {
	g.reporter.AnnounceSuite(description)
}

func (g *Gotcha) AnnounceSuccesfulSpec(spec *types.SpecSummary) {
	g.reporter.PrintSingleSpec(spec, prefixOK, color.GreenString)
}

func (g *Gotcha) AnnounceSuccesfulSlowSpec(spec *types.SpecSummary, quiet bool) {
	g.reporter.PrintSingleSpec(spec, prefixOK, color.GreenString)
}

func (g *Gotcha) AnnounceSuccesfulMeasurement(spec *types.SpecSummary, quiet bool) {
	g.reporter.PrintSingleSpec(spec, prefixOK, color.GreenString)
}

func (g *Gotcha) AnnouncePendingSpec(spec *types.SpecSummary, verbose bool) {
	g.reporter.PrintSingleSpec(spec, prefixPending, color.YellowString)
}

func (g *Gotcha) AnnounceSpecTimedOut(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	g.reporter.PrintSingleSpec(spec, prefixFail, color.RedString)
}

func (g *Gotcha) AnnounceSpecPanicked(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	g.reporter.PrintSingleSpec(spec, prefixFail, color.RedString)
}

func (g *Gotcha) AnnounceSpecFailed(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	g.reporter.PrintSingleSpec(spec, prefixFail, color.RedString)
}

func (g *Gotcha) AnnounceSpecRunCompletion(summary *types.SuiteSummary, quiet bool) {
	var fn ColorFunc
	if summary.NumberOfFailedSpecs > 0 {
		fn = color.RedString
	} else if summary.NumberOfPassedSpecs == summary.NumberOfTotalSpecs {
		fn = color.GreenString
	} else {
		fn = color.YellowString
	}
	g.reporter.PrintSummary(summary, fn)
}

func (g *Gotcha) SummarizeFailures(summaries []*types.SpecSummary) {
	var failures []*types.SpecSummary
	for i := range summaries {
		if summaries[i].HasFailureState() {
			failures = append(failures, summaries[i])
		}
	}

	var pendings []*types.SpecSummary
	for i := range summaries {
		if summaries[i].Pending() {
			pendings = append(pendings, summaries[i])
		}
	}

	g.reporter.SummarizeFailures(failures, pendings)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (g *Gotcha) AnnounceAggregatedParallelRun(nodes int, quiet bool) {
	// nothing
}

func (g *Gotcha) AnnounceParallelRun(a int, b int, quiet bool) {
	// nothing
}

func (g *Gotcha) AnnounceNumberOfSpecs(specsToRun int, total int, quiet bool) {
	// nothing
}

func (g *Gotcha) AnnounceTotalNumberOfSpecs(specs int, quiet bool) {
	// nothing
}

func (g *Gotcha) AnnounceSpecWillRun(spec *types.SpecSummary) {
	// nothing
}

func (g *Gotcha) AnnounceSkippedSpec(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	// nothing
}

func (g *Gotcha) AnnounceCapturedOutput(output string) {
	// nothing
}

func (g *Gotcha) AnnounceBeforeSuiteFailure(summary *types.SetupSummary, quiet bool, fullTrace bool) {
	// nothing
}

func (g *Gotcha) AnnounceAfterSuiteFailure(summary *types.SetupSummary, quiet bool, fullTrace bool) {
	// nothing
}
