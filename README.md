# Gotcha

Gotcha is an RSpec-like [ginkgo](https://github.com/onsi/ginkgo) reporter.

```
import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/ruggi/gotcha"
)

func TestExample(t *testing.T) {
    RegisterFailHandler(Fail)
    gotcha.RunSpecs(t, "Example Suite")
}
```
