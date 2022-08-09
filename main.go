package main

import (
	"log"

	"github.com/HUN220/magdam/internal/cli"
	"github.com/HUN220/magdam/internal/pull"
	"github.com/HUN220/magdam/internal/push"
)

func main() {
	// TODO: Add `magdam pull -continue` (continue on error)
	// TODO: Add ability to pull a single item

	var opts cli.CliOptions
	if err := cli.NewOptions(&opts); err != nil {
		log.Fatal(err)
	}

	switch opts.Command {
	case "pull":
		pull.PullCmd(opts.ApiKeyId, opts.ApiKey, opts.BaseUrl)

	case "push":
		push.PushCmd(opts.ApiKeyId, opts.ApiKey, opts.BaseUrl)

	default:
		log.Fatal("command not yet implemented: ", opts.Command)
	}
}
