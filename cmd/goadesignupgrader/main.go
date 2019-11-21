package main

import (
	"github.com/tchssk/goadesignupgrader"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(goadesignupgrader.Analyzer) }