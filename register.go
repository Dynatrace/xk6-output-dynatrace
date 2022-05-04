package dynatracewriter

import (
	"github.com/henrikrexed/xk6-dynatrace-output/pkg/dynatracewriter"
	"go.k6.io/k6/output"
)

func init() {
	output.RegisterExtension("output-dynatrace", func(p output.Params) (output.Output, error) {
		return dynatracewriter.New(p)
	})
}