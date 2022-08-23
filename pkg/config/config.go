/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package config

const (
	DefaultParser   = "kindle-html"
	DefaultExporter = "json"
	DefaultStorage  = "mongodb"
)

var GlobalCfg GlobalConfig

type GlobalConfig struct {
	PromptYesToAll bool
	PromptNoToAll  bool
}

type ConvertConfig struct {
	ListParsers   bool
	ListExporters bool

	InputPath string
	OutputDir string

	Parser   string
	Exporter string
}

type StorageConfig struct {
	ListStorages bool

	Storage string

	Filter string
}
