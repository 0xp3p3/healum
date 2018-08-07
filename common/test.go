package common

import (
	"context"
	"github.com/micro/go-micro/metadata"
	"strings"
	"net/http"
)

var testMeta = metadata.Metadata{
	"test_meta": "test_meta",
}

func IsTestContext(ctx context.Context) bool {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return false
	}
	val, ok := md["test_meta"]
	if !ok || val != "test_meta"{
		return false
	}
	return true

}

func NewContextByHeader(ctx context.Context, header http.Header) context.Context{
	if header.Get("WWW-Testing") == "test_meta"{
		return NewTestContext(ctx)
	}
	return ctx
}

func NewTestContext(ctx context.Context) context.Context {
	return metadata.NewContext(ctx, testMeta)
}

func SetTestHeader(header http.Header) {
	header.Set("WWW-Testing", "test_meta")
}

func TestingName(name string) string {
	if !strings.HasSuffix(name, "_test"){
		return name + "_test"
	}
	return name
}
