package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/walle/wiki"
)

// Version of the tool.
const Version = "1.4.0"

// Exit statuses for the tool.
const (
	UsageErrorExitStatus   = 1
	NoSuchPageExitStatus   = 2
	RequestErrorExitStatus = 3
	ParsingErrorExitStatus = 4
	SuccessExitStatus      = 0
)

var buf = bytes.NewBufferString("")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `wiki is a tool used to fetch excerpts from wikipedia
Usage: wiki [options...] query
Options:
`)
		flag.PrintDefaults()
	}

	def_lang := os.Getenv("WIKI_LANG")
	def_url := os.Getenv("WIKI_URL")

	if def_lang == "" {
		def_lang = "en"
	}
	if def_url == "" {
		def_url = "https://%s.wikipedia.org/w/api.php"
	}

	language := flag.String("l", def_lang, "The language to use")
	url := flag.String("u", def_url, "The api url")
	noColor := flag.Bool("n", false, "If the output should not be colorized")
	simple := flag.Bool("s", false, "If simple output should be used")
	short := flag.Bool("short", false, "If short output should be used")
	wrap := flag.Int("w", 0, "The width text should be wrapped at. 0 is no wrap.")
	noCheckCert := flag.Bool(
		"no-check-certificate",
		false,
		"Skip verification of certificates",
	)
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

	page := getPage(url, language, noCheckCert)

	if page.Content == "" {
		fmt.Fprintf(os.Stderr, "No such page\n")
		if !*simple {
			fmt.Printf("Create it on: %s\n", page.URL)
		}
		os.Exit(NoSuchPageExitStatus)
	}

	if *simple {
		printPageSimple(page)
	} else if *short {
		printPageShort(page)
	} else if *noColor {
		printPagePlain(page)
	} else {
		printPageColor(page)
	}

	out := colorable.NewColorableStdout()
	if *wrap > 0 {
		paragraphs := strings.Split(buf.String(), "\n\n")
		for _, p := range paragraphs {
			fmt.Fprintln(out, wiki.Wrap(p, *wrap), "\n")
		}
	} else {
		fmt.Fprintln(out, buf.String())
	}

	os.Exit(SuccessExitStatus)
}

func getPage(url, language *string, noCheckCert *bool) *wiki.Page {
	query := strings.Join(flag.Args(), " ")
	req, err := wiki.NewRequest(*url, query, *language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create request %s\n", err)
		os.Exit(RequestErrorExitStatus)
	}

	resp, err := req.Execute(*noCheckCert)
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
		fmt.Fprint(buf, "Redirected from %s to %s\n\n",
			page.Redirect.From,
			page.Redirect.To,
		)
	}
	fmt.Fprintln(buf, page.Content)
	fmt.Fprint(buf, "\nRead more: %s\n", page.URL)
}

func printPageSimple(page *wiki.Page) {
	fmt.Fprintln(buf, page.Content)
}

func printPageShort(page *wiki.Page) {
	fmt.Fprintln(buf, page.Content[:strings.Index(page.Content, ".")+1])
}

func printPageColor(page *wiki.Page) {
	if page.Redirect != nil {
		fmt.Fprintf(buf,
			"\x1b[31m"+
				"Redirected from "+
				"\x1b[41;37m%s\x1b[49;31m to \x1b[41;37m%s"+
				"\x1b[0m"+
				"\n\n",
			page.Redirect.From,
			page.Redirect.To,
		)
	}
	fmt.Fprintln(buf, page.Content)
	fmt.Fprintf(buf, "\n\x1b[32mRead more: %s\x1b[0m\n", page.URL)
}
