package gotcha

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/types"
)

type summarizer struct{}

func (s summarizer) printSummary(summary *types.SuiteSummary, fn ColorFunc) {
	fmt.Println()
	fmt.Printf("Finished in %.4f seconds\n", summary.RunTime.Seconds())

	out := []string{
		s.renderPartial("example", summary.NumberOfTotalSpecs, true),
		s.renderPartial("failure", summary.NumberOfFailedSpecs, true),
	}
	if summary.NumberOfPendingSpecs > 0 {
		out = append(out, s.renderPartial("pending", summary.NumberOfPendingSpecs, false))
	}

	fmt.Println(fn(strings.Join(out, ", ")))
	fmt.Println()
}

func (s summarizer) renderPartial(str string, n int, pluralize bool) string {
	if pluralize && n != 1 {
		str += "s"
	}
	return fmt.Sprintf("%d %s", n, str)
}

func (s summarizer) printFailures(failures []*types.SpecSummary, pendings []*types.SpecSummary) {
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
		fmt.Println()
	}

	if len(pendings) > 0 {
		fmt.Println("Pending:")
		fmt.Println()

		for i, pending := range pendings {
			fmt.Println(fmt.Sprintf("  %d) %s", i+1, strings.Join(pending.ComponentTexts[1:], " ")))
			loc := pending.ComponentCodeLocations[len(pending.ComponentCodeLocations)-1]
			color.Cyan(fmt.Sprintf("     %s", loc))
		}
	}
}
