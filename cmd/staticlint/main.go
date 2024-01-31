package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/dsbasko/yandex-go-shortener/internal/lint"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers, lint.GetXPassesAnalyzers()...)
	analyzers = append(analyzers, lint.GetStaticCheckAnalyzers()...)
	analyzers = append(analyzers, lint.GetStyleCheckAnalyzers()...)
	analyzers = append(analyzers, lint.GetExternalAnalyzers()...)

	multichecker.Main(analyzers...)
}
