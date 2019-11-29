package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kr/text"
	"github.com/mattn/go-colorable"

	"github.com/walle/wiki"
)

var Version = "1.4.1" // nolint:gochecknoglobals

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
		fmt.Fprintln(os.Stderr, `wiki is a tool used to fetch excerpts from wikipedia
Usage: wiki [options...] query
Options:`)
		flag.PrintDefaults()
	}

	defLang, defURL := handleDefaults()

	language := flag.String("l", defLang, "The language to use")
	url := flag.String("u", defURL, "The api url")
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

	handleFlags(version, help)

	page := getPage(url, language, noCheckCert)

	if page.Content == "" {
		fmt.Fprintln(os.Stderr, "No such page")
		if !*simple {
			fmt.Printf("Create it on: %s\n", page.URL)
		}
		os.Exit(NoSuchPageExitStatus)
	}

	var buf bytes.Buffer
	switch {
	case *simple:
		printPageSimple(page, &buf)
	case *short:
		printPageShort(page, &buf)
	case *noColor:
		printPagePlain(page, &buf)
	default:
		printPageColor(page, &buf)
	}

	out := colorable.NewColorableStdout()
	if *wrap > 0 {
		paragraphs := strings.Split(buf.String(), "\n\n")
		for _, p := range paragraphs {
			fmt.Fprintln(out, text.Wrap(p, *wrap))
		}
		os.Exit(SuccessExitStatus)
	}

	fmt.Fprintln(out, buf.String())
	os.Exit(SuccessExitStatus)
}

func handleDefaults() (string, string) {
	defLang := os.Getenv("WIKI_LANG")
	defURL := os.Getenv("WIKI_URL")
	if defLang == "" {
		defLang = "en"
	}
	if defURL == "" {
		defURL = "https://%s.wikipedia.org/w/api.php"
	}
	return defLang, defURL
}

func handleFlags(version *bool, help *bool) {
	if *version {
		fmt.Fprintf(os.Stdout, "wiki %s\n", Version)
		os.Exit(SuccessExitStatus)
	}

	if *help {
		flag.Usage()
		os.Exit(SuccessExitStatus)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(UsageErrorExitStatus)
	}
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

func printPagePlain(page *wiki.Page, w io.Writer) {
	if page.Redirect != nil {
		fmt.Fprintf(w, "Redirected from %s to %s\n\n",
			page.Redirect.From,
			page.Redirect.To,
		)
	}
	fmt.Fprintln(w, page.Content)
	fmt.Fprintf(w, "\nRead more: %s\n", page.URL)
}

func printPageSimple(page *wiki.Page, w io.Writer) {
	fmt.Fprintln(w, page.Content)
}

func printPageShort(page *wiki.Page, w io.Writer) {
	fmt.Fprintln(w, page.Content[:strings.Index(page.Content, ".")+1])
}

func printPageColor(page *wiki.Page, w io.Writer) {
	if page.Redirect != nil {
		fmt.Fprintf(w,
			"\x1b[31m"+
				"Redirected from "+
				"\x1b[41;37m%s\x1b[49;31m to \x1b[41;37m%s"+
				"\x1b[0m"+
				"\n\n",
			page.Redirect.From,
			page.Redirect.To,
		)
	}
	fmt.Fprintln(w, page.Content)
	fmt.Fprintf(w, "\n\x1b[32mRead more: %s\x1b[0m\n", page.URL)
}
