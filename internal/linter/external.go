package linter

import (
	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
)

func GetExternalAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		bodyclose.Analyzer,
		rowserr.NewAnalyzer(
			"github.com/jmoiron/sqlx",
		),
	}
}
