/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package json

import (
	jsonenc "encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
)

type JSONExporter struct {
	prettyPrint bool
	indent      string
}

func (e *JSONExporter) Name() string {
	return "json"
}

func (e *JSONExporter) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&e.prettyPrint, "json.pretty", false, "print the json with indent")
	cmd.PersistentFlags().StringVar(&e.indent, "json.indent", "  ", "sets the json indent")
}

func (e *JSONExporter) Export(cfg *config.ConvertConfig, books []*model.Book) error {
	var b []byte
	var err error
	if e.prettyPrint {
		b, err = jsonenc.MarshalIndent(books, "", e.indent)
	} else {
		b, err = jsonenc.Marshal(books)
	}
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}
	util.Log(string(b))
	return nil
}
