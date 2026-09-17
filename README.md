# Gamgee 🪴🪏🧑‍🌾

Gamgee is a tiny DNS sinkhole written in Golang. It supports Wildcard Domains types of domain lists, you can find some examples from e.g. [HaGeZi's DNS-Blocklists](https://github.com/hagezi/dns-blocklists).

## Configuring

Gamgee reads config.toml file.

Example config file:
```toml
port = 53530
upstreamDns = "10.0.0.1:53"
blockListURIs = [
        "https://example.com/dns-blocklist-wildcard-domains.txt",
        "https://example.com/another-cool-domain-list.txt"
]
```

