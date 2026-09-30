package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/miekg/dns"
	"github.com/jliuhtonen/gamgee/internal/blocklist"
	"github.com/jliuhtonen/gamgee/internal/config"
)

func replyWithError(w dns.ResponseWriter, msg *dns.Msg, rcode uint16) {
	reply := msg.Copy()
	reply.Response = true
	reply.Rcode = rcode
	reply.Data = nil
	reply.WriteTo(w)
}

func cacheKey(rr *dns.RR) string {
	header := (*rr).Header()
	return strings.ToLower(header.Name) + "|" +
		strconv.Itoa(int(dns.RRToType(*rr))) + "|" +
		strconv.Itoa(int(header.Class))
}

func main() {
	config, err := config.ReadConfig()

	if err != nil {
		panic(err)
	}

	blockList, err := blocklist.FetchLists(config.BlocklistURIs)

	if err != nil {
		panic(err)
	}

	listenAddr := ":" + strconv.Itoa(config.Port)

	client := dns.NewClient()

	dns.ListenAndServe(listenAddr, "udp", dns.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) {
		if msg.Opcode != dns.OpcodeQuery {
			replyWithError(w, msg, dns.RcodeNotImplemented)
			return
		}

		if len(msg.Question) != 1 {
			replyWithError(w, msg, dns.RcodeFormatError)
			return
		}

		q := msg.Question[0]
		key := cacheKey(&q)
		fmt.Println("Cache key", key)

		domain, _ := strings.CutSuffix(q.Header().Name, ".")
		if blockList.Contains(domain) {
			replyWithError(w, msg, dns.RcodeNameError)
			return
		}

		respMsg, _, err := client.Exchange(ctx, msg, "udp", config.UpstreamDns)
		if err != nil {
			replyWithError(w, msg, dns.RcodeServerFailure)
			return
		}
		respMsg.WriteTo(w)
	}))
}
