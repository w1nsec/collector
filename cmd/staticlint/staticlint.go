package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/alexkohler/nakedret"
	"github.com/gnieto/mulint/mulint"

	"github.com/w1nsec/golinters/linters"
)

func AddStaticchecks(additionals []string) []*analysis.Analyzer {

	m := make(map[string]bool)
	for _, v := range additionals {
		m[v] = true
	}

	checks := make([]*analysis.Analyzer, 0)
	for _, v := range staticcheck.Analyzers {
		// ALL SA checks
		if strings.Contains(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		}
		// other checks
		if m[v.Analyzer.Name] {
			checks = append(checks, v.Analyzer)
		}
	}

	return checks
}

func main() {
	var (
		//readConfig  = true
		additional = []string{
			// SA	static checks
			// ALL

			// S	simple checks
			"S1011", // Use a single append to concatenate two slices
			"S1008", // Simplify returning boolean expression

			// ST	style checks
			"ST1001", // Dot imports are discouraged
			"ST1019", // Importing the same package multiple times

			// QF	quickfix
			"QF1001", // Apply De Morgan’s law

			// SA 	checks
			"SA4006",
			"SA5000",
			"SA6000",
			"SA9004",
		}

		DefaultLines = uint(5)
	)

	// standard checks from "golang.org/x/tools/go/analysis/passes/..."
	checks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,

		// analyzer from example
		linters.ErrCheckAnalyzer,
		linters.OSExitCheckAnalyzer,

		// third-party analyzers
		nakedret.NakedReturnAnalyzer(DefaultLines),
		mulint.Mulint,
	}

	checks = append(checks, AddStaticchecks(additional)...)

	multichecker.Main(
		checks...,
	)

}
