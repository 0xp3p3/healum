// Package ip_whitelist is a micro plugin for whitelisting ip addresses
package ip_whitelist

import (
	"net"
	"net/http"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/micro/plugin"
)

type whitelist struct {
	cidrs map[string]*net.IPNet
	ips   map[string]bool
}

func (w *whitelist) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "ip_whitelist",
			Usage:  "Comma separated whitelist of allowed IPs",
			EnvVar: "IP_WHITELIST",
		},
	}
}

func (w *whitelist) load(ips ...string) {
	for _, ip := range ips {
		parts := strings.Split(ip, "/")

		switch len(parts) {
		// assume just an ip
		case 1:
			w.ips[ip] = true
		case 2:
			// parse cidr
			_, ipnet, err := net.ParseCIDR(ip)
			if err != nil {
				log.Fatalf("[ip_whitelist] failed to parse %v: %v", ip, err)
			}
			w.cidrs[ipnet.String()] = ipnet
		default:
			log.Fatalf("[ip_whitelist] failed to parse %v", ip)
		}
	}

}

func (w *whitelist) match(ip string) bool {
	// make ip
	nip := net.ParseIP(ip)

	// check ips
	if ok := w.ips[nip.String()]; ok {
		return true
	}

	// check cidrs
	for _, cidr := range w.cidrs {
		if cidr.Contains(nip) {
			return true
		}
	}

	// no match
	return false
}

func (w *whitelist) Commands() []cli.Command {
	return nil
}

func (w *whitelist) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// check remote addr; if we can't parse it passes through
			if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				// reject if no match
				if !w.match(ip) {
					http.Error(rw, "forbidden", 403)
					return
				}
			}

			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

func (w *whitelist) Init(ctx *cli.Context) error {
	if whitelist := ctx.String("ip_whitelist"); len(whitelist) > 0 {
		w.load(strings.Split(whitelist, ",")...)
	}
	return nil
}

func (w *whitelist) String() string {
	return "ip_whitelist"
}

func NewIPWhitelist(ips ...string) plugin.Plugin {
	// create plugin
	w := &whitelist{
		cidrs: make(map[string]*net.IPNet),
		ips:   make(map[string]bool),
	}

	// load ips
	w.load(ips...)

	return w
}
