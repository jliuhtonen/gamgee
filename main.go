package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"codeberg.org/miekg/dns"
	"github.com/jliuhtonen/warden/internal/blocklist"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: warden <upstream dns address>")
	}
	upstreamAddr := args[0]
	blockList, err := blocklist.FetchList()

	if err != nil {
		panic(err)
	}

	client := dns.NewClient()

	dns.ListenAndServe(":53530", "udp", dns.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) {
		fmt.Println(msg.Question)
		for _, q := range msg.Question {
			domain, _ := strings.CutSuffix(q.Header().Name, ".")
			if blockList.Contains(domain) {
				fmt.Println("BLOCKING ", domain)
				reply := msg.Copy()
				reply.Response = true
				reply.Rcode = dns.RcodeNameError
				reply.Data = nil
				reply.WriteTo(w)
				return
			}
		}
		fmt.Println("Received message" + msg.String())
		respMsg, _, err := client.Exchange(ctx, msg, "udp", upstreamAddr)
		if err != nil {
			fmt.Println("ERROR" + err.Error())
			return
		}
		respMsg.WriteTo(w)
	}))
}
