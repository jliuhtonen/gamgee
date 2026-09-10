package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jliuhtonen/warden/internal/blocklist"
	"github.com/miekg/dns"
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

	client := &dns.Client{
		Net: "udp",
	}

	conn, err := client.Dial(upstreamAddr)
	if err != nil {
		panic(err)
	}

	dns.ListenAndServe(":53530", "udp", dns.HandlerFunc(func(w dns.ResponseWriter, msg *dns.Msg) {
		fmt.Println(msg.Question)
		for _, q := range msg.Question {
			domain, _ := strings.CutSuffix(q.Name, ".")
			if blockList.Contains(domain) {
				fmt.Println("BLOCKING ", domain)
				reply := msg.SetReply(msg)
				reply.Rcode = dns.RcodeNameError
				w.WriteMsg(reply)
				return
			}
		}
		fmt.Println("Received message" + msg.String())
		respMsg, _, err := client.ExchangeWithConn(msg, conn)
		if err != nil {
			fmt.Println("ERROR" + err.Error())
		}
		w.WriteMsg(respMsg)
	}))
}
