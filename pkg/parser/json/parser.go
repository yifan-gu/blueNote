/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package json

import (
	jsonenc "encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yifan-gu/blueNote/pkg/model"
)

type JSONParser struct {
	author string
	title  string
	stdin  bool
}

func (p *JSONParser) Name() string {
	return "json"
}

func (p *JSONParser) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&p.stdin, "json.stdin", false, "Treat the input as a json object")
	cmd.PersistentFlags().StringVar(&p.author, "json.author", "", "override the book author name")
	cmd.PersistentFlags().StringVar(&p.title, "json.title", "", "override the book title name")
}

func (p *JSONParser) Parse(inputPath string) ([]*model.Book, error) {
	var books []*model.Book
	var data []byte
	var err error

	if p.stdin {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		if inputPath == "" {
			return nil, errors.New("Input file is missing!")
		}
		data, err = ioutil.ReadFile(inputPath)
	}
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if err := jsonenc.Unmarshal(data, &books); err != nil {
		return nil, errors.Wrap(err, "")
	}

	if p.author != "" {
		for _, bk := range books {
			bk.Author = p.author
			for i := range bk.Marks {
				bk.Marks[i].Author = p.author
			}
		}
	}
	if p.title != "" {
		for _, bk := range books {
			bk.Title = p.title
			for i := range bk.Marks {
				bk.Marks[i].Title = p.title
			}
		}
	}
	return books, nil
}
