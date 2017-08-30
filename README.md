# Gotcha

Gotcha is an RSpec-like [ginkgo](https://github.com/reifcode/ginkgo) reporter.

```
import (
    "testing"

    . "github.com/reifcode/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/reifcode/gotcha"
)

func TestExample(t *testing.T) {
    RegisterFailHandler(Fail)
    gotcha.RunSpecs(t, "Example Suite")
}
```
