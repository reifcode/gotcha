package gotcha

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/ginkgo/reporters"
)

func RunSpecs(t ginkgo.GinkgoTestingT, description string) bool {
	reporters := []ginkgo.Reporter{NewReporter()}
	return ginkgo.RunSpecsWithCustomReporters(t, description, reporters)
}

func NewReporter() reporters.Reporter {
	stenographer := NewGotcha(config.DefaultReporterConfig)
	return reporters.NewDefaultReporter(config.DefaultReporterConfig, stenographer)
}
