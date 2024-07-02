package search

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type CrawlData struct {
	Url          string
	Success      bool
	ResponseCode int
	CrawlBody    ParsedBody
}

type ParsedBody struct {
	CrawlTime       time.Duration
	PageTitle       string
	PageDescription string
	Headings        string
	Links           Links
}

type Links struct {
	Internal []string
	External []string
}

var userAgents = []string{
	"Mozilla/5.0 (Linux; U; Android 4.0.3; en-us; KFTT Build/IML74K) AppleWebKit/537.36 (KHTML, like Gecko) Silk/3.68 like Chrome/39.0.2171.93 Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 8_2 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12D508 Safari/600.1.4",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:39.0) Gecko/20100101 Firefox/39.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Xbox; Xbox One) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edge/44.18363.8131",
}

// randomUserAgent generates a random user agent
// to prevent the server from blocking us.
func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

// makeRequest creates the request to the given url
// and sets a random user agent for that request.
// Returns an http response reference or an error.
func makeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func runCrawl(inputUrl string) CrawlData {
	// resp, err := http.Get(inputUrl)
	resp, err := makeRequest(inputUrl)
	baseUrl, _ := url.Parse(inputUrl)
	if err != nil || resp == nil {
		log.Printf("something went wrong fetch the body: %s\n", err)

		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: 0,
			CrawlBody:    ParsedBody{},
		}
	}
	defer resp.Body.Close()

	// Check if response code is not 200
	if resp.StatusCode != 200 {
		log.Printf("non 200 code found: %d\n", resp.StatusCode)

		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlBody:    ParsedBody{},
		}
	}

	// Check the content type is text/html
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		// response is HTML
		data, err := parseBody(resp.Body, baseUrl)
		if err != nil {
			log.Printf(
				"something went wrong getting data from html body: %s\n", err,
			)

			return CrawlData{
				Url:          inputUrl,
				Success:      false,
				ResponseCode: resp.StatusCode,
				CrawlBody:    ParsedBody{},
			}
		}

		return CrawlData{
			Url:          inputUrl,
			Success:      true,
			ResponseCode: resp.StatusCode,
			CrawlBody:    data,
		}

	} else {
		// response is not HTML
		log.Println("non html response detected")

		return CrawlData{
			Url:          inputUrl,
			Success:      false,
			ResponseCode: resp.StatusCode,
			CrawlBody:    ParsedBody{},
		}
	}
}

func parseBody(body io.Reader, baseUrl *url.URL) (ParsedBody, error) {
	doc, err := html.Parse(body)
	if err != nil {
		log.Printf("something went wrong parsing body: %s\n", err)

		return ParsedBody{}, err
	}

	// Record timings
	start := time.Now()

	// Get the links from the doc
	links := getLinks(doc, baseUrl)

	// Get the page title & description
	title, desc := getPageData(doc)

	// Get the h1 tags for the page
	headings := getPageHeadings(doc)

	// Record timings
	end := time.Now()

	// Return the data
	return ParsedBody{
		CrawlTime:       end.Sub(start),
		PageTitle:       title,
		PageDescription: desc,
		Headings:        headings,
		Links:           links,
	}, nil
}

func isSameHost(absoluteUrl, baseUrl string) bool {
	absUrl, err := url.Parse(absoluteUrl)
	if err != nil {
		return false
	}

	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return false
	}

	return absUrl.Host == baseUrlParsed.Host
}

func checkUrlKind(url *url.URL) bool {
	uStr := url.String()

	return strings.HasPrefix(uStr, "#") ||
		strings.HasPrefix(uStr, "mail") ||
		strings.HasPrefix(uStr, "tel") ||
		strings.HasPrefix(uStr, "javascript") ||
		strings.HasSuffix(uStr, ".pdf") ||
		strings.HasSuffix(uStr, ".md")
}

// getLinks does a Depth First Search (DFS) of the html tree structure.
// This is a recursive function to scan the full tree.
// ↓ See note below about Depth-First Search (DFS) algorithm ↓
func getLinks(n *html.Node, baseUrl *url.URL) Links {
	links := Links{}
	if n == nil {
		return links
	}

	// Recursive search for "a" tags in the DOM tree
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		// Check if the current node is an `html.ElementNode`
		// and if it has a tag name of "a" (i.e., an anchor tag).
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url, err := url.Parse(attr.Val)
					// Check for errors or if url is:
					// 1) a hashtag/anchor,
					// 2) is mail link,
					// 3) is a telephone link,
					// 4)is a javascript link,
					// 5) is a PDF or MD file
					if err != nil || checkUrlKind(url) {
						continue
					}

					// If url is absolute then test if internal or extend
					// before append. Else add the baseUrl append as internal
					if url.IsAbs() {
						if isSameHost(url.String(), baseUrl.String()) {
							links.Internal = append(
								links.Internal, url.String(),
							)
						} else {
							links.External = append(
								links.External, url.String(),
							)
						}
					} else {
						// Check to see example:
						// https://go.dev/src/net/url/example_test.go
						rel := baseUrl.ResolveReference(url)
						links.Internal = append(
							links.Internal, rel.String(),
						)
					}
				}
			}
		}

		// Recursively call function to do Depth First Search of entire tree
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}

	// We call the previously defined function
	findLinks(n)

	return links
}

func getPageData(n *html.Node) (string, string) {
	title, desc := "", ""
	if n == nil {
		return title, desc
	}

	// Find the page title & description

	// Recursive function to search for `meta` elements
	// in the HTML tree and extracts their `name` and `content` attributes.
	var findMetaAndTitle func(*html.Node)
	findMetaAndTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			// Check if first child is empty
			if n.FirstChild == nil {
				title = ""
			} else {
				title = n.FirstChild.Data
			}
		} else if n.Type == html.ElementNode && n.Data == "meta" {
			var name, content string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "content" {
					content = attr.Val
				}
			}

			if name == "description" {
				desc = content
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findMetaAndTitle(c)
		}
	}

	// We call the previously defined function
	findMetaAndTitle(n)

	return title, desc
}

func getPageHeadings(n *html.Node) string {
	if n == nil {
		return ""
	}

	// Find all h1 elements and concatenate their content
	var (
		headings strings.Builder
		findH1   func(*html.Node)
	)

	// Recursive search for h1 headings
	findH1 = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			// Check if first child is empty
			if n.FirstChild != nil {
				headings.WriteString(n.FirstChild.Data)
				headings.WriteString(", ")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findH1(c)
		}
	}

	// We call the previously defined function
	findH1(n)

	// Remove the last comma and space
	// from the concatenated string & return
	return strings.TrimSuffix(headings.String(), ", ")
}

/* DEPTH-FIRST SEARCH ALGORITHM IN GO:
https://reintech.io/blog/depth-first-search-algorithm-in-go
https://www.tutorialspoint.com/golang-program-to-implement-depth-first-search

https://www.google.com/search?q=golang+Depth+First+Search+(DFS)&oq=golang+Depth+First+Search+(DFS)&aqs=chrome..69i57j33i160.3311j0j7&sourceid=chrome&ie=UTF-8

USER-AGENTS LIST:
https://gist.github.com/pzb/b4b6f57144aea7827ae4

*/
