package search

import (
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseBody(t *testing.T) {
	// Create a sample HTML body
	body := strings.NewReader(`
		<html>
			<head>
				<meta name="description" content="Page Description">
				<title>Page Title</title>
			</head>
			<body>
				<h1>Heading 1</h1>
				<a href="https://example.com">Internal Link</a>
				<a href="https://external.com">External Link</a>>
				<a href="/internal">Internal Link</a>
			</body>
		</html>
	`)

	baseUrl, _ := url.Parse("https://example.com")

	expectedPageTitle := "Page Title"
	expectedPageDesc := "Page Description"
	expectedHeadings := "Heading 1"
	expectedInternalLinks := []string{
		"https://example.com", "https://example.com/internal",
	}
	expectedExternalLinks := []string{"https://external.com"}

	// Call the function `parseBody`
	result, err := parseBody(body, baseUrl)

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Compare the page title result with the expected value
	if result.PageTitle != expectedPageTitle {
		t.Errorf(
			"Expected page title '%s', but got '%s'",
			expectedPageTitle,
			result.PageTitle,
		)
	}

	// Compare the page description result with the expected value
	if result.PageDescription != expectedPageDesc {
		t.Errorf(
			"Expected page description '%s', but got '%s'",
			expectedPageDesc,
			result.PageDescription,
		)
	}

	// Compare the headings result with the expected value
	if result.Headings != expectedHeadings {
		t.Errorf(
			"Expected headings '%s', but got '%s'",
			expectedHeadings,
			result.Headings,
		)
	}

	// Compare the internal links result with the expected value
	if !equalSlices(result.Links.Internal, expectedInternalLinks) {
		t.Errorf(
			"Expected internal links '%v', but got '%v'",
			expectedInternalLinks,
			result.Links.Internal,
		)
	}

	// Compare the external links result with the expected value
	if !equalSlices(result.Links.External, expectedExternalLinks) {
		t.Errorf(
			"Expected external links '%v', but got '%v'",
			expectedExternalLinks,
			result.Links.External,
		)
	}
}

func TestGetLinks(t *testing.T) {
	// Create a sample HTML node
	doc, _ := html.Parse(strings.NewReader(`
		<html>
			<body>
				<a href="https://example.com">Internal Link</a>
				<a href="https://external.com">External Link</a>
				<div><a href="/internal">Internal Link</a></div>
				<a href="#section">Anchor Link</a>
				<a href="mailto:info@example.com">Mail Link</a>
				<a href="tel:+1234567890">Telephone Link</a>
				<a href="javascript:void(0)">JavaScript Link</a>
				<a href="document.pdf">PDF Link</a>
				<a href="document.md">MD Link</a>
			</body>
		</html>
	`))

	baseUrl, _ := url.Parse("https://example.com")

	expectedInternal := []string{
		"https://example.com",
		"https://example.com/internal",
	}
	expectedExternal := []string{"https://external.com"}

	linksch := make(chan Links)

	// Call the function `getLinks`
	go getLinks(doc, baseUrl, linksch)

	result := <-linksch

	// Compare the internal links result with the expected value
	if !equalSlices(result.Internal, expectedInternal) {
		t.Errorf(
			"Expected internal links '%v', but got '%v'",
			expectedInternal,
			result.Internal,
		)
	}

	// Compare the external links result with the expected value
	if !equalSlices(result.External, expectedExternal) {
		t.Errorf(
			"Expected external links '%v', but got '%v'",
			expectedExternal,
			result.External,
		)
	}
}

// Helper function to check if two string slices are equal
// (or reflect.DeepEqual(a, b))
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestIsSameHost(t *testing.T) {
	// Define test cases
	testCases := []struct {
		absoluteUrl string
		baseUrl     string
		expected    bool
	}{
		{"https://example.com/path", "https://example.com", true},
		{"https://example.com/path", "https://www.example.com", false},
		{"https://example.com", "https://example.com", true},
		{"https://example.com", "https://example.org", false},
		{"https://example.com", "http://example.com", true},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		result := isSameHost(tc.absoluteUrl, tc.baseUrl)

		// Compare the result with the expected value
		if result != tc.expected {
			t.Errorf(
				"For absoluteUrl '%s' and baseUrl '%s', expected '%t', but got '%t'",
				tc.absoluteUrl,
				tc.baseUrl,
				tc.expected,
				result,
			)
		}
	}
}

func TestGetPageData(t *testing.T) {
	// Create a sample HTML node
	doc, _ := html.Parse(strings.NewReader(`
		<html>
			<head>
				<meta name="description" content="Page Description">
				<title>Page Title</title>
			</head>
			<body>
				<h1>Heading 1</h1>
				<p>Some content</p>
			</body>
		</html>
	`))

	expectedTitle := "Page Title"
	expectedDesc := "Page Description"

	titlech, descch := make(chan string), make(chan string)

	// Call the function `getPageData`
	go getPageData(doc, titlech, descch)

	resultTitle, resultDesc := <-titlech, <-descch

	// Compare the title result with the expected value
	if resultTitle != expectedTitle {
		t.Errorf(
			"Expected title '%s', but got '%s'", expectedTitle, resultTitle,
		)
	}

	// Compare the descrition result with the expected value
	if resultDesc != expectedDesc {
		t.Errorf(
			"Expected description '%s', but got '%s'", expectedDesc, resultDesc,
		)
	}
}

func TestGetPageHeadings(t *testing.T) {
	// Create a sample HTML node
	doc, _ := html.Parse(strings.NewReader(`
		<html>
			<body>
				<h1>Heading 1</h1>
				<div>
					<h1>Heading 2</h1>
				</div>
				<h2>Not a heading</h2>
				<h1></h1>
			</body>
		</html>
	`))

	expected := "Heading 1, Heading 2"

	headingsch := make(chan string)

	// Call the function `getPageHeadings`
	go getPageHeadings(doc, headingsch)

	result := <-headingsch

	// Compare the result with the expected value
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

/* 3 WAYS TO COMPARE SLICES:
https://yourbasic.org/golang/compare-slices/
https://go.dev/play/p/tzl9Z3ofn3W

COMMAND FOR RUN ALL TESTS:
go test ./...
*/
