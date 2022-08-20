/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package config

const (
	DefaultParser   = "kindle-html"
	DefaultExporter = "org-roam"
)

type ConvertConfig struct {
	ListParsers   bool
	ListExporters bool

	InputPath string
	OutputDir string

	Parser   string
	Exporter string

	Author string
	Title  string

	SplitBook    bool
	AuthorSubDir bool

	PromptYesToAll bool
	PromptNoToAll  bool
}
