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

func main() {
	config, err := config.ReadConfig()
	listenAddr := ":" + strconv.Itoa(config.Port)

	if err != nil {
		panic(err)
	}

	blockList, err := blocklist.FetchLists(config.BlocklistURIs)

	if err != nil {
		panic(err)
	}

	client := dns.NewClient()

	dns.ListenAndServe(listenAddr, "udp", dns.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) {
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
		respMsg, _, err := client.Exchange(ctx, msg, "udp", config.UpstreamDns)
		if err != nil {
			fmt.Println("ERROR" + err.Error())
			return
		}
		respMsg.WriteTo(w)
	}))
}
