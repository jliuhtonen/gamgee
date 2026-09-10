package blocklist

import (
	"io"
	"net/http"
	"slices"
	"strings"
)

const domainList = "https://cdn.jsdelivr.net/gh/hagezi/dns-blocklists@latest/wildcard/pro.mini-onlydomains.txt"

func FetchList() (*BlockList, error) {
	response, err := http.Get(domainList)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	blockList := parseList(string(responseContent))

	return &blockList, nil
}

type BlockList struct {
	items []string
}

func (bl *BlockList) Contains(domain string) bool {
	domain = strings.ToLower(domain)
	return slices.ContainsFunc(bl.items, func(i string) bool {
		return domain == i || strings.HasSuffix(domain, "."+i)
	})
}

func parseList(listStr string) BlockList {
	var items []string
	for _, line := range strings.Split(listStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		items = append(items, strings.ToLower(line))
	}
	return BlockList{items: items}
}
