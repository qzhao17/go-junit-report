package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/jstemmer/go-junit-report/pkg/gtr"
	"github.com/jstemmer/go-junit-report/pkg/parser/gotest"
)

var (
	noXMLHeader   = flag.Bool("no-xml-header", false, "do not print xml header")
	packageName   = flag.String("package-name", "", "specify a package name (compiled test have no package name in output)")
	goVersionFlag = flag.String("go-version", "", "specify the value to use for the go.version property in the generated XML")
	setExitCode   = flag.Bool("set-exit-code", false, "set exit code to 1 if tests failed")
	printEvents = flag.Bool("print-events", false, "print events (for debugging)")
)

func main() {
	flag.Parse()

	if flag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "%s does not accept positional arguments\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	// Read input
	events, err := gotest.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %s\n", err)
		os.Exit(1)
	}

	if *printEvents {
		for i, ev := range events {
			fmt.Printf("%02d: %#v\n", i, ev)
		}
	}
	report := gtr.FromEvents(events)

	if !*noXMLHeader {
		fmt.Fprintf(os.Stdout, xml.Header)
	}

	// TODO: write xml header?
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	if err := enc.Encode(gtr.JUnit(report)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing XML: %s\n", err)
		os.Exit(1)
	}
	if err := enc.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flusing XML: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "\n")

	if *setExitCode && report.HasFailures() {
		os.Exit(1)
	}
}

// TODO: read/write + test
