/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package config

const (
	DefaultParser   = "kindle-html"
	DefaultExporter = "org-roam"
)

type Config struct {
	InputPath string
	OutputDir string

	Parser   string
	Exporter string

	SplitBook    bool
	AuthorSubDir bool

	PromptYesToAll bool
	PromptNoToAll  bool
}
