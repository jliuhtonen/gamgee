package blocklist

import (
	"io"
	"net/http"
	"strings"
)

const domainList = "https://cdn.jsdelivr.net/gh/hagezi/dns-blocklists@latest/wildcard/pro.mini-onlydomains.txt"

func FetchList() (*DomainList, error) {
	response, err := http.Get(domainList)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsedDomains := parseList(string(responseContent))
	domainList := New(parsedDomains)

	return domainList, nil
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
