// Package router is a micro plugin for defining HTTP routes
package router

import (
	"errors"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/micro/go-os/log"
	"github.com/micro/micro/plugin"
)

type Option func(o *Options)

type router struct {
	opts Options

	// TODO: optimise for concurrency
	sync.RWMutex
	routes Routes
}

var (
	// Default config source file
	DefaultFile   = "routes.json"
	DefaultPath   = []string{"api"}
	DefaultLogger = log.NewLog(log.WithOutput(log.NewOutput(log.OutputName("/dev/stderr"))))
)

func (r *router) update(routes Routes) {
	// sort routes
	sort.Sort(sortedRoutes{routes})
	// update
	r.Lock()
	r.routes = routes
	r.Unlock()
}

func (r *router) run(c config.Config) {
	var routes Routes

	// load routes immediately if possible
	if err := c.Get(DefaultPath...).Scan(&routes); err != nil {
		DefaultLogger.Error("[router] Failed to get routes", err)
	} else {
		r.update(routes)
	}

	var watcher config.Watcher

	// try to get a watch
	for i := 0; i < 100; i++ {
		w, err := c.Watch(DefaultPath...)
		if err != nil {
			DefaultLogger.Error("[router] Failed to get watcher", err)
			time.Sleep(time.Second)
			continue
		}
		watcher = w
		break
	}

	// if the watch is nil we exit
	if watcher == nil {
		DefaultLogger.Fatal("[router] Failed to get watcher in 100 attempts")
	}

	// watch and update routes
	for {
		// get next
		v, err := watcher.Next()
		if err != nil {
			DefaultLogger.Error("[router] Watcher Next() Error", err)
			time.Sleep(time.Second)
			continue
		}

		var routes Routes

		// scan into routes
		if err := v.Scan(&routes); err != nil {
			DefaultLogger.Error("[router] Failed to scan routes... skipping update", err)
			continue
		}

		// update the routes
		r.update(routes)
	}
}

func (r *router) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "config_source",
			EnvVar: "CONFIG_SOURCE",
			Usage:  "Source to read the config from e.g file:path/to/file, platform",
		},
	}
}

func (r *router) Commands() []cli.Command {
	return nil
}

func (r *router) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// get routes
			r.RLock()
			routes := r.routes
			r.RUnlock()

			// routes are ordered on update
			for _, route := range routes.Routes {
				// route matched write and return
				if route.Match(req) {
					route.Write(w, req)
					return
				}
			}

			// serve the default handler
			h.ServeHTTP(w, req)
		})
	}
}

func (r *router) Init(ctx *cli.Context) error {
	// TODO: Make this more configurable and add more sources
	var conf config.Config

	if c := ctx.String("config_source"); len(c) == 0 && r.opts.Config == nil {
		return errors.New("config source must be defined")
	} else if len(c) > 0 {
		var source config.Source

		switch c {
		case "platform":
			source = config.NewSource()
		case "file":
			fileName := DefaultFile

			parts := strings.Split(c, ":")

			if len(parts) > 1 {
				fileName = parts[1]
			}

			source = file.NewSource(config.SourceName(fileName))
		default:
			return errors.New("Unknown config source " + c)
		}

		conf = config.NewConfig(config.WithSource(source))
	} else {
		conf = r.opts.Config
	}

	go r.run(conf)

	return nil
}

func (r *router) String() string {
	return "router"
}

func NewRouter(opts ...Option) plugin.Plugin {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	r := &router{
		opts: options,
	}

	if options.Config != nil {
		var routes Routes
		if err := options.Config.Get(DefaultPath...).Scan(&routes); err == nil {
			r.update(routes)
		}
	}

	return r
}
