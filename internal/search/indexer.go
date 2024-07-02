package search

import "github.com/emarifer/search-engine/internal/services"

// Index is an in-memory inverted index. It maps tokens to url IDs.
type Index map[string][]string

// Add adds documents to the Index.
func (idx Index) Add(docs []services.CrawledUrl) {
	for _, doc := range docs {
		docString := doc.Url + " " +
			doc.PageTitle + " " +
			doc.PageDescription + " " +
			doc.Headings
		for _, token := range analyze(docString) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Don't add same ID twice.
				continue
			}

			idx[token] = append(ids, doc.ID)
		}
	}
}

/* REFERENCES.- AKHIL SHARMA REPOSITORY AND VIDEO:
https://github.com/AkhilSharma90/go-full-text-search
https://www.youtube.com/watch?v=BPLpzpgp79A
*/
