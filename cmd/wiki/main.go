package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/walle/wiki"
)

// Version of the tool.
const Version = "1.0.0"

// Exit statuses for the tool.
const (
	UsageErrorExitStatus   = 1
	NoSuchPageExitStatus   = 2
	RequestErrorExitStatus = 3
	ParsingErrorExitStatus = 4
	SuccessExitStatus      = 0
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `wiki is a tool used to fetch exerpts from wikipedia
Usage: wiki [options...] query
Options:
`)
		flag.PrintDefaults()
	}

	language := flag.String("l", "en", "The language to use")
	url := flag.String("u", "https://%s.wikipedia.org/w/api.php", "The api url")
	noColor := flag.Bool("n", false, "If the output should not be colorized")
	simple := flag.Bool("s", false, "If simple output should be used")
	help := flag.Bool("h", false, "Print help information and exit.")
	version := flag.Bool("version", false, "Print version information and exit.")

	flag.Parse()

	// If version is requested, print info and exit
	if *version {
		fmt.Fprintf(os.Stdout, "wiki %s\n", Version)
		os.Exit(SuccessExitStatus)
	}

	// If help is requested, print info and exit
	if *help {
		flag.Usage()
		os.Exit(SuccessExitStatus)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(UsageErrorExitStatus)
	}

	page := getPage(url, language)

	if page.Content == "" {
		fmt.Fprintf(os.Stderr, "No such page\n")
		if !*simple {
			fmt.Printf("Create it on: %s\n", page.URL)
		}
		os.Exit(NoSuchPageExitStatus)
	}

	if *simple {
		printPageSimple(page)
	} else if *noColor {
		printPagePlain(page)
	} else {
		printPageColor(page)
	}

	os.Exit(SuccessExitStatus)
}

func getPage(url, language *string) *wiki.Page {
	query := strings.Title(strings.Join(flag.Args(), " "))
	req, err := wiki.NewRequest(*url, query, *language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create request %s\n", err)
		os.Exit(RequestErrorExitStatus)
	}

	resp, err := req.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not execute request %s\n", err)
		os.Exit(RequestErrorExitStatus)
	}

	page, err := resp.Page()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse data %s\n", err)
		os.Exit(ParsingErrorExitStatus)
	}

	return page
}

func printPagePlain(page *wiki.Page) {
	if page.Redirect != nil {
		fmt.Printf("Redirected from %s to %s\n\n",
			page.Redirect.From,
			page.Redirect.To,
		)
	}
	fmt.Println(page.Content)
	fmt.Printf("\nRead more: %s\n", page.URL)
}

func printPageSimple(page *wiki.Page) {
	fmt.Println(page.Content)
}

func printPageColor(page *wiki.Page) {
	if page.Redirect != nil {
		fmt.Printf(
			"\x1b[31m"+
				"Redirected from "+
				"\x1b[41;37m%s\x1b[49;31m to \x1b[41;37m%s"+
				"\x1b[0m"+
				"\n\n",
			page.Redirect.From,
			page.Redirect.To,
		)
	}
	fmt.Println(page.Content)
	fmt.Printf("\n\x1b[32mRead more: %s\x1b[0m\n", page.URL)
}
