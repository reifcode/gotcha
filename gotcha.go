package gotcha

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/ginkgo/types"
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
	stenographer := NewGotcha(config.DefaultReporterConfig)
	return reporters.NewDefaultReporter(config.DefaultReporterConfig, stenographer)
}

type Gotcha struct {
	prefix  bool
	current int
	levels  map[int][]string
}

func NewGotcha(cfg config.DefaultReporterConfigType) *Gotcha {
	prefix := false
	if cfg.NoColor {
		color.NoColor = true
		prefix = true
	}

	return &Gotcha{
		levels: make(map[int][]string),
		prefix: prefix,
	}
}

func (g *Gotcha) AnnounceSuite(description string, randomSeed int64, randomizingAll bool, quiet bool) {
	fmt.Println()
	fmt.Println(description)
	fmt.Println()
}

func (g *Gotcha) AnnounceSuccesfulSpec(spec *types.SpecSummary) {
	g.printSingleSpec(spec, prefixOK, color.Green)
}

func (g *Gotcha) AnnounceSuccesfulSlowSpec(spec *types.SpecSummary, quiet bool) {
	g.printSingleSpec(spec, prefixOK, color.Green)
}

func (g *Gotcha) AnnounceSuccesfulMeasurement(spec *types.SpecSummary, quiet bool) {
	g.printSingleSpec(spec, prefixOK, color.Green)
}

func (g *Gotcha) AnnouncePendingSpec(spec *types.SpecSummary, verbose bool) {
	g.printSingleSpec(spec, prefixPending, color.Yellow)
}

func (g *Gotcha) AnnounceSpecTimedOut(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	g.printSingleSpec(spec, prefixFail, color.Red)
}

func (g *Gotcha) AnnounceSpecPanicked(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	g.printSingleSpec(spec, prefixFail, color.Red)
}

func (g *Gotcha) AnnounceSpecFailed(spec *types.SpecSummary, quiet bool, fullTrace bool) {
	g.printSingleSpec(spec, prefixFail, color.Red)
}

func (g *Gotcha) AnnounceSpecRunCompletion(summary *types.SuiteSummary, quiet bool) {
	fmt.Println()
	fmt.Println(fmt.Sprintf("Finished in %.4f seconds", summary.RunTime.Seconds()))

	var fn colorFunc
	if summary.NumberOfFailedSpecs > 0 {
		fn = color.Red
	} else if summary.NumberOfPassedSpecs == summary.NumberOfTotalSpecs {
		fn = color.Green
	} else {
		fn = color.Yellow
	}
	g.printSummary(summary, fn)

	fmt.Println()
}

func (g *Gotcha) SummarizeFailures(summaries []*types.SpecSummary) {
	var failures []*types.SpecSummary
	for i := range summaries {
		if summaries[i].HasFailureState() {
			failures = append(failures, summaries[i])
		}
	}
	if len(failures) > 0 {
		fmt.Println()
		fmt.Println("Failures:")
		fmt.Println()

		for i, failed := range failures {
			failure := failed.Failure
			fmt.Printf("  %d) %s\n", i+1, strings.Join(failed.ComponentTexts[1:], " "))
			lines := strings.Split(failure.Message, "\n")
			for i := range lines {
				lines[i] = fmt.Sprintf("     %s", lines[i])
			}
			color.Red(strings.Join(lines, "\n"))
			color.Cyan(fmt.Sprintf("     %s", failure.Location.String()))
			fmt.Println()
		}
	}

	var pendings []*types.SpecSummary
	for i := range summaries {
		if summaries[i].Pending() {
			pendings = append(pendings, summaries[i])
		}
	}
	if len(pendings) > 0 {
		fmt.Println()
		fmt.Println("Pending:")
		fmt.Println()

		for i, pending := range pendings {
			fmt.Println(fmt.Sprintf("  %d) %s", i+1, strings.Join(pending.ComponentTexts[1:], " ")))
			loc := pending.ComponentCodeLocations[len(pending.ComponentCodeLocations)-1]
			color.Cyan(fmt.Sprintf("     %s", loc))
		}
	}
}

type colorFunc func(string, ...interface{})

func (g *Gotcha) printSingleSpec(spec *types.SpecSummary, prefix string, fn colorFunc) {
	size := len(spec.ComponentTexts[1:]) - 1
	for i, component := range spec.ComponentTexts[1:] {
		level := g.levels[i]
		found := false
		for _, c := range level {
			if component == c {
				found = true
				break
			}
		}
		if i > len(level) || !found {
			g.levels[i] = append(level, component)
			spaces := strings.Repeat("  ", i)
			if i == size {
				if g.prefix {
					fn(fmt.Sprintf("%s%s%s", spaces, prefix, component))
				} else {
					fn(fmt.Sprintf("%s%s", spaces, component))
				}
			} else {
				fmt.Println(fmt.Sprintf("%s%s", spaces, component))
			}
		}
	}
}

func (g *Gotcha) renderStatPartial(s string, n int, pluralize bool) string {
	if pluralize && n != 1 {
		s += "s"
	}
	return fmt.Sprintf("%d %s", n, s)
}

func (g *Gotcha) printSummary(summary *types.SuiteSummary, fn colorFunc) {
	var out []string
	out = append(out, g.renderStatPartial("example", summary.NumberOfTotalSpecs, true))
	out = append(out, g.renderStatPartial("failure", summary.NumberOfFailedSpecs, true))

	if summary.NumberOfPendingSpecs > 0 {
		out = append(out, g.renderStatPartial("pending", summary.NumberOfPendingSpecs, false))
	}

	fn(strings.Join(out, ", "))
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
