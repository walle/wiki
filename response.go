package wiki

import (
	"net/url"
	"time"
)

// Page contains the parsed data.
type Page struct {
	ID       int
	Title    string
	Content  string
	Language string
	URL      string
	Redirect *redirect
}

// Response contains the raw data the API returns.
type Response struct {
	Batchcomplete string
	Query         query
}

// Page parses the raw data and returns a Page with the relevant data.
func (r *Response) Page() (*Page, error) {
	page := &Page{}

	if len(r.Query.Redirects) > 0 {
		page.Redirect = &r.Query.Redirects[0]
	}

	for _, p := range r.Query.Pages {
		url, err := url.QueryUnescape(p.Canonicalurl)
		if err != nil {
			url = p.Canonicalurl
		}
		page.ID = p.Pageid
		page.Title = p.Title
		page.Content = p.Extract
		page.Language = p.Pagelanguage
		page.URL = url

		break
	}

	return page, nil
}

type query struct {
	Redirects []redirect
	Pages     map[string]page
}

type redirect struct {
	From string
	To   string
}

type page struct {
	Pageid       int
	Ns           int
	Title        string
	Extract      string
	Contentmodel string
	Pagelanguage string
	Touched      time.Time
	Fullurl      string
	Canonicalurl string
}
