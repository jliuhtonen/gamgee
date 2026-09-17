package domainmatcher

import (
	"github.com/hashicorp/go-immutable-radix/v2"
	"slices"
	"strings"
)

type domainMatcherValue struct{}

type DomainMatcher struct {
	domainLookupTree *iradix.Tree[domainMatcherValue]
}

func (dl *DomainMatcher) Contains(domain string) bool {
	key := reverseDomainKey(domain)
	_, _, ok := dl.domainLookupTree.Root().LongestPrefix(key)

	return ok
}

func New(domains []string) *DomainMatcher {
	tree := insertDomainsToTree(iradix.New[domainMatcherValue](), domains)
	return &DomainMatcher{
		domainLookupTree: tree,
	}
}

func (dl *DomainMatcher) Append(domains []string) *DomainMatcher {
	return &DomainMatcher{
		domainLookupTree: insertDomainsToTree(dl.domainLookupTree, domains),
	}
}

func insertDomainsToTree(tree *iradix.Tree[domainMatcherValue], domains []string) *iradix.Tree[domainMatcherValue] {
	tx := tree.Txn()
	for _, domain := range domains {
		tx.Insert(reverseDomainKey(strings.ToLower(domain)), domainMatcherValue{})
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
