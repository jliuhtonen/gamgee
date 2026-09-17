package blocklist

import (
	"context"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/jliuhtonen/gamgee/internal/domainmatcher"
	"golang.org/x/sync/errgroup"
)

func FetchLists(uris []string) (*domainmatcher.DomainMatcher, error) {
	results := make([][]string, len(uris))

	grp, _ := errgroup.WithContext(context.Background())

	for i, uri := range uris {
		grp.Go(func() error {
			result, err := fetchList(uri)
			results[i] = result
			return err
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, err
	}

	return domainmatcher.New(slices.Concat(results...)), nil
}

func fetchList(domainListURI string) ([]string, error) {
	response, err := http.Get(domainListURI)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsedDomains := parseList(string(responseContent))

	return parsedDomains, nil
}

func parseList(listStr string) []string {
	var items []string

	for _, line := range strings.Split(listStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		items = append(items, strings.ToLower(line))
	}

	return items
}
