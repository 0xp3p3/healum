package db

import (
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type Options struct {
	Database string
	Table    string

	Client client.Client

	// For alternative options
	Context context.Context
}

type SearchOptions struct {
	Metadata Metadata
	Limit    int64
	Offset   int64
}

func Database(d string) Option {
	return func(o *Options) {
		o.Database = d
	}
}

func Table(t string) Option {
	return func(o *Options) {
		o.Table = t
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Search Options

func WithMetadata(md Metadata) SearchOption {
	return func(o *SearchOptions) {
		o.Metadata = md
	}
}

func WithLimit(l int64) SearchOption {
	return func(o *SearchOptions) {
		o.Limit = l
	}
}

func WithOffset(ot int64) SearchOption {
	return func(o *SearchOptions) {
		o.Offset = ot
	}
}
