# goddns

Your IP address changes. Your DNS record does not. goddns fixes that.

It watches your public IP address and updates your DNS records whenever it changes — so your homelab, VPN, or self-hosted whatever stays reachable, even when your ISP decides to hand you a new address at 3am for no particular reason.

## How it works

goddns is built around three concepts:

- **Retrievers** fetch your current public IP address (e.g. by asking your UniFi controller, or a public API like ifconfig.co)
- **Providers** push that IP to your DNS records (e.g. Hetzner Cloud DNS)
- **Instances** tie one retriever to one or more providers and run on a configurable interval

You can run multiple instances simultaneously — useful if you have several sites or networks whose IPs need to be tracked independently.

## Configuration

Copy `config.example.toml` to `config.toml` and adjust to taste:

```toml
interval = "5m"

[retriever.home]
type = "ifconfigco"

[provider.my-provider]
type      = "hetzner_cloud"
api_token = "your-token-here"
zone      = "example.com"
rr_name   = "@"

[instance.homelab]
[instance.homelab.retriever]
name = "home"

[[instance.homelab.provider]]
name = "my-provider"
```

Provider and retriever parameters can be overridden per instance, so one provider definition can serve multiple zones:

```toml
[[instance.homelab.provider]]
name = "my-provider"

[[instance.homelab.provider]]
name    = "my-provider"
zone    = "another-zone.com"
rr_name = "subdomain"
```

See `config.example.toml` for the full reference and [`docs/PARAMETERS.md`](docs/PARAMETERS.md) for all available parameters per plugin.

## Running

The recommended way to run goddns is via Docker:

```sh
docker run -d \
  -v /path/to/config.toml:/config.toml \
  -e GODDNS_CONFIG=/config.toml \
  ghcr.io/tinyoverflow/goddns:latest
```

Images are published to the GitHub Container Registry for every release. Replace `latest` with a specific version tag for reproducible deployments.

## Adding a plugin

Create a file in `internal/plugins/`, name it `provider_<name>.go` or `retriever_<name>.go`, implement the `plugin.Provider` or `plugin.Retriever` interface, and register it in `init()`:

```go
package plugins

import "goddns/internal/plugin"

type myConfig struct {
    Endpoint string `json:"endpoint" required:"true" doc:"API endpoint"`
}

func init() {
    plugin.RegisterProvider("my_provider", newFromConfig, myConfig{})
}

func newFromConfig(params map[string]any) (plugin.Provider, error) {
    cfg, err := plugin.Decode[myConfig](params)
    // ...
}
```

Run `go generate` to update the parameter documentation.

## Building from source

```sh
go build ./cmd/goddns
```

---

> **AI Disclaimer:** This project used AI assistance (Claude by Anthropic) for architectural design decisions and to generate `cmd/gendoc`, the tool that produces the plugin parameter documentation.
