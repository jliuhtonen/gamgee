package blocklist

import (
	"github.com/hashicorp/go-immutable-radix/v2"
	"slices"
	"strings"
)

type domainListValue struct{}

type DomainList struct {
	domainLookupTree *iradix.Tree[domainListValue]
}

func (dl *DomainList) Contains(domain string) bool {
	key := reverseDomainKey(domain)
	_, _, ok := dl.domainLookupTree.Root().LongestPrefix(key)

	return ok
}

func New(domains []string) *DomainList {
	tree := insertDomainsToTree(iradix.New[domainListValue](), domains)
	return &DomainList{
		domainLookupTree: tree,
	}
}

func insertDomainsToTree(tree *iradix.Tree[domainListValue], domains []string) *iradix.Tree[domainListValue] {
	tx := tree.Txn()
	for _, domain := range domains {
		tx.Insert(reverseDomainKey(strings.ToLower(domain)), domainListValue{})
	}
	return tx.Commit()
}

func reverseDomainKey(domain string) []byte {
	return []byte((reverseDomain(domain) + "."))
}

func reverseDomain(domain string) string {
	splitDomain := strings.Split(domain, ".")
	slices.Reverse(splitDomain)
	return strings.Join(splitDomain, ".")
}
