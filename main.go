package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ASCII Colors
const (
	CLRG = "\x1b[32;1m" // green
	CLRN = "\x1b[0m"    // no color
)

var reCSF = regexp.MustCompile(`(\w+)(?:=|\()([^,)]+)`)
var reDefTerm = regexp.MustCompile(`\([^,)]+\){1}$`)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	s := bufio.NewScanner(os.Stdin)

	// // read from file for local testing
	// f, err := os.Open("f1.txt")
	// check(err)
	// s := bufio.NewScanner(f)
	// // end read from file

	for i := 1; s.Scan(); i++ {
		fmt.Printf("%s### flow %v ###%s\n", CLRG, i, CLRN)
		parseDumpFlow(i, s.Text())
		fmt.Printf("\n\n")
	}
}

// combineComplexFlowEntries takes the split by comma flow entry and combines the
// cmplex elements like ipv4(src=10.100.1.5,dst=92.222.79.96,proto=6,tos=0,ttl=64,frag=no)
// returns the []string slice with flow terms
func combineComplexFlowEntries(flow *string) []string {
	flowSplit := strings.Split(*flow, ",")
	var terms []string
	var buf []string
	var curDepth int
	for _, term := range flowSplit {
		// fmt.Println(term)
		buf = append(buf, term)
		if strings.Contains(term, "(") {
			curDepth = curDepth + strings.Count(term, "(")
		}
		if strings.Contains(term, ")") {
			curDepth = curDepth - strings.Count(term, ")")
		}

		if curDepth == 0 {
			terms = append(terms, strings.Join(buf, ","))
			buf = []string{}
		}

	}
	// fmt.Printf("debug from combineComplexFlowEntries: %v\n", terms)
	return terms
}

func parseTerms(terms []string) {
	actions := false
	curActionType := ""
	for i, t := range terms {
		t = strings.Trim(t, " ") // strip blanks
		switch {
		// last element of the flow term
		case reDefTerm.MatchString(t):
			parseDefTerm(t, 0)
		case i == len(terms)-1:
			parseLastTerm(t)
		// terms related to actions
		case strings.HasPrefix(t, "actions") || actions:
			if actions == false {
				actions = true
				fmt.Println("actions:")
			}
			if strings.HasPrefix(t, "actions") {
				t = strings.TrimPrefix(t, "actions:")
			}
			parseActionTerm(t, &curActionType)
		case reCSF.MatchString(t) && !actions:
			parseDefComplexTerm(t, 0)
		}
	}
}

// parseDefTerm applies a default parsing logic for the flow term
// that does not need a specific parsing algo
// it strips the ending ')' char and replaces leading whitespaces
// also changing '(' to ': '
// i.e. in_port(12) -> in_port: 12
// or packets:11 -> packets: 11
func parseDefTerm(t string, indent int) {
	// t = t[:len(t)-1]
	// fmt.Println(t)
	t = strings.Trim(t, " )")
	t = strings.Replace(t, ":", ": ", 1)
	t = strings.Replace(t, "(", ": ", 1)
	fmt.Printf("%s%s\n", strings.Repeat(" ", indent), t)
}

// parseDefComplexTerm parses fields of a default complex term
// which comes as a set of fields with nested key:values
// i.e. eth(dst=68:54:ed:00:7f:c6) or eth(src=5a:eb:de:fe:6b:44,dst=00:00:0a:0a:01:0b)
func parseDefComplexTerm(t string, indent int) {
	idx := strings.Index(t, "(") // position of the first bracket
	fmt.Printf("%s%s:\n", strings.Repeat(" ", indent), t[:idx])
	t = t[idx+1 : len(t)-1] // strip term from the enclosure (i.e. strips 'eth()')
	m := reCSF.FindAllStringSubmatch(t, -1)
	// log.Printf("submatch: %v", m)
	for i := range m {
		for j := 1; j < len(m[i]); j++ {
			switch j % 2 {
			case 1:
				{
					fmt.Printf("%s%s: ", strings.Repeat(" ", indent+2), m[i][j])
				}
			case 0:
				{
					fmt.Printf("%s", m[i][j])
					if m[i][j-1] == "tun_id" {
						i, _ := strconv.ParseInt(m[i][j], 0, 64)
						fmt.Printf(" (%v)\n", i)
					} else {
						fmt.Println()
					}

				}
			}
		}

	}
}

func parseActionTerm(t string, curActionType *string) {
	newActionType := t[:strings.Index(t, "(")]
	if newActionType != *curActionType {
		fmt.Printf("  %s:\n", newActionType)
		*curActionType = newActionType
	}
	t = t[strings.Index(t, "(")+1:]
	if reCSF.MatchString(t) {
		parseDefComplexTerm(t, 4)
	} else {
		parseDefTerm(t, 4)
	}

}

func parseLastTerm(t string) {
	fmt.Printf("out_port: %s\n", t)
}

func parseDumpFlow(i int, flow string) {
	terms := combineComplexFlowEntries(&flow)
	// fmt.Println(terms)
	parseTerms(terms)
}
