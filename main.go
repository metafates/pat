package main

import (
	"github.com/metafates/pat/cmd"
	"github.com/metafates/pat/config"
	"github.com/metafates/pat/log"
	"github.com/samber/lo"
)

func main() {
	lo.Must0(config.Init())
	lo.Must0(log.Init())
	cmd.Execute()
}
