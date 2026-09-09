package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

const domainList = "https://cdn.jsdelivr.net/gh/hagezi/dns-blocklists@latest/wildcard/pro.mini-onlydomains.txt"

func fetchList() (string, error) {
	response, err := http.Get(domainList)
	if err != nil {
		return "", nil
	}
	defer response.Body.Close()
	responseContent, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(responseContent), nil
}

type BlockList struct {
	items []string
}

func (bl *BlockList) Contains(domain string) bool {
	return slices.Contains(bl.items, domain)
}

func parseList(listStr string) BlockList {
	blocked := strings.Split(listStr, "\n")
	return BlockList{items: blocked}
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: warden <domain>")
		os.Exit(0)
	}
	domain := args[0]
	listStr, err := fetchList()
	if err != nil {
		panic(err)
	}
	blockList := parseList(listStr)
	if blockList.Contains(domain) {
		fmt.Println("BLOCKED")
	} else {
		fmt.Println("That's not blocked!")
	}
}
