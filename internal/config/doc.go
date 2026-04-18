// Package config manages stashctl application configuration.
//
// Configuration is stored as JSON and supports loading from a file path,
// falling back to sensible defaults when no file exists. Use [DefaultPath]
// to resolve the canonical config location (~/.stashctl/config.json).
//
// Example usage:
//
//	path, _ := config.DefaultPath()
//	cfg, err := config.Load(path)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// modify and persist
//	cfg.DefaultFmt = "detail"
//	config.Save(path, cfg)
package config
