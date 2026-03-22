# Plugin Parameters

- [Retrievers](#retrievers)
  - [ifconfigco](#ifconfigco)
  - [unifi](#unifi)

- [Providers](#providers)
  - [hetznercloud](#hetznercloud)
  - [spaceship](#spaceship)

## Retrievers

### `ifconfigco`

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `base_url` | `string` | No | "https://ifconfig.co" | API base URL |

### `unifi`

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `base_url` | `string` | Yes |  | Unifi controller base URL (e.g. https://192.168.1.1) |
| `api_token` | `string` | Yes |  | API authentication token |
| `site_id` | `string` | No | "default" | Site ID |
| `verify_tls` | `bool` | No | false | Enable TLS certificate verification |

## Providers

### `hetznercloud`

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `api_token` | `string` | Yes |  | Hetzner Cloud API token |
| `zone` | `string` | Yes |  | DNS zone (e.g. homelab.com) |
| `rr_name` | `string` | Yes |  | Resource record name (e.g. @) |

### `spaceship`

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `api_key` | `string` | Yes |  | Spaceship API key |
| `api_secret` | `string` | Yes |  | Spaceship API secret |
| `domain` | `string` | Yes |  | The domain name (e.g. example.com) |
| `subdomain` | `string` | Yes |  | The (sub)domain name |
| `ttl` | `int` | No | 3600 | The TTL (Time to live) for the entry |

