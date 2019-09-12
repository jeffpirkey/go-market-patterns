package tools

import (
	"fmt"
	"github.com/pkg/errors"
	"go-market-patterns/mal"
	"go-market-patterns/model/core"
	"os"
	"sort"
	"strings"
)

const (
	mdPatternHeader = "| Value | Up | Down | No Change | Total |\r\n" +
		"|------:|---:|-----:|----------:|------:|\r\n"
)

func PrintMarkdownPatterns(repo *mal.Repos, symbol string, f *os.File) error {

	var patterns core.PatternSlice
	patterns, err := repo.PatternRepo.FindBySymbol(symbol)
	if err != nil {
		return errors.Wrapf(err, "problem printing markdown for %v", symbol)
	}

	sort.Sort(patterns)

	var sb strings.Builder
	sb.WriteString("# Patterns for " + strings.ToUpper(symbol) + "\r\n\r\n")
	sb.WriteString(mdPatternHeader)
	for _, pattern := range patterns {
		sb.WriteString(fmt.Sprintf("| %v | %v | %v | %v | %v |\r\n",
			pattern.Value, pattern.UpCount, pattern.DownCount, pattern.NoChangeCount, pattern.TotalCount))
	}

	_, err = f.WriteString(sb.String())
	if err != nil {
		return errors.Wrapf(err, "problem writing pattern to file for %v", symbol)
	}

	return nil
}
