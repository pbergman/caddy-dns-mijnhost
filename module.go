package mijn_host

import (
	"strconv"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/mijnhost"
	"github.com/pbergman/provider"
	"go.uber.org/zap"
)

// Provider lets Caddy read and manipulate DNS records hosted by this DNS provider.
type Provider struct{ *mijnhost.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.mijnhost",
		New: func() caddy.Module { return &Provider{new(mijnhost.Provider)} },
	}
}

// Provision sets up the module. Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {

	p.Provider.ApiKey = caddy.NewReplacer().ReplaceAll(p.Provider.ApiKey, "")

	if p.Provider.DebugLevel > 0 {
		p.Provider.DebugOut = zap.NewStdLog(ctx.Logger()).Writer()
	}

	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	mijnhost [<api_key>] {
//	    api_key			<api_key>
//		debug_level 	<debug level for client 0, 1, 2 or 3>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {

		if d.NextArg() {
			p.Provider.ApiKey = d.Val()
		}

		if d.NextArg() {
			return d.ArgErr()
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_key":
				if p.Provider.ApiKey != "" {
					return d.Err("Api Key already set")
				}
				if d.NextArg() {
					p.Provider.ApiKey = d.Val()
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			case "debug_level":

				if d.NextArg() {

					if val, err := strconv.Atoi(d.Val()); err == nil && (val > 0 && val < 4) {
						p.Provider.DebugLevel = provider.OutputLevel(val)
					}
				}

				if d.NextArg() {
					return d.ArgErr()
				}

			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.ApiKey == "" {
		return d.Err("missing API Key")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
