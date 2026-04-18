package main

import (
	"fmt"
	"os"

	"github.com/user/stashctl/internal/config"
	"github.com/user/stashctl/internal/store"
)

func main() {
	cfg, err := config.Load(config.DefaultPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	st, err := store.New(cfg.StorePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening store: %v\n", err)
		os.Exit(1)
	}

	root := newRootCmd(cfg, st)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
