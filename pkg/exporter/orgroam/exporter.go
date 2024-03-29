/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package orgroam

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/exporter/orgroam/db"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/util"
)

const (
	defaultRoamDBPath   = "~/.emacs.d/.local/etc/org-roam.db"
	defaultSqlDriver    = "sqlite3"
	defaultTemplateType = 0
)

type Book struct {
	Title  string
	Author string
	Marks  []Mark
	UUID   uuid.UUID
}

type Mark struct {
	Type     string
	Section  string
	Location Location
	Data     string
	UserNote string
	Pos      int
	UUID     uuid.UUID
}

type Location struct {
	Chapter  string
	Page     *int
	Location *int
}

func (l Location) String() string {
	if l.Chapter != "" {
		return fmt.Sprintf("Chapter: %s Page: %d Loc: %d", l.Chapter, l.Page, l.Location)
	}
	return fmt.Sprintf("Page: %d Loc: %d", l.Page, l.Location)
}

func convertFromModelBook(book *model.Book) *Book {
	bk := &Book{
		Title:  book.Title,
		Author: book.Author,
	}
	for _, mk := range book.Marks {
		mark := Mark{
			Type:     mk.Type,
			Section:  mk.Section,
			Data:     mk.Data,
			UserNote: mk.UserNote,
			Location: Location{
				Chapter:  mk.Location.Chapter,
				Page:     mk.Location.Page,
				Location: mk.Location.Location,
			},
		}
		bk.Marks = append(bk.Marks, mark)
	}
	return bk
}

func writeRunesToFile(fullpath string, runes []rune) error {
	f, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to open or create file %s", fullpath))
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	defer buf.Flush()

	for i := range runes {
		_, err = fmt.Fprintf(buf, "%c", runes[i])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to write to file %s", fullpath))
		}
	}
	return nil
}

type OrgRoamExporter struct {
	updateRoamDB   bool
	roamDBPath     string
	dbDriver       string
	templateType   int
	insertRoamLink bool
	authorSubDir   bool
}

func (e *OrgRoamExporter) Name() string {
	return "org-roam"
}

func (e *OrgRoamExporter) LoadConfigs(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&e.updateRoamDB, "org-roam.update-db", false, "automatically update the roam sqlite db for links")
	cmd.PersistentFlags().StringVarP(&e.roamDBPath, "org-roam.db-path", "d", defaultRoamDBPath, "path to the org-roam sqlite3 database")
	cmd.PersistentFlags().StringVar(&e.dbDriver, "org-roam.db-driver", defaultSqlDriver, "the database driver to use")
	cmd.PersistentFlags().BoolVarP(&e.insertRoamLink, "org-roam.insert-roam-link", "l", true, "insert the roam links")
	cmd.PersistentFlags().IntVar(&e.templateType, "org-roam.template-type", defaultTemplateType, "the type of the template to use")
	cmd.PersistentFlags().BoolVar(&e.authorSubDir, "org-roam.author-subdir", true, "create sub-directory with the name of the author")
}

func (e *OrgRoamExporter) Export(cfg *config.ConvertConfig, books []*model.Book) error {
	for _, bk := range books {
		if err := e.exportBook(cfg, bk); err != nil {
			return errors.Wrap(err, "")
		}
	}
	return nil
}

func (e *OrgRoamExporter) exportBook(cfg *config.ConvertConfig, book *model.Book) error {
	bk := convertFromModelBook(book)

	sq, err := db.NewSqlInterface(e.roamDBPath, e.dbDriver)
	if err != nil {
		return err
	}
	defer sq.Close()

	fullpath, err := util.ResolvePath(e.generateOutputPath(bk, cfg))
	if err != nil {
		return err
	}
	dir := filepath.Dir(fullpath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		confirm, err := util.PromptExportOverrideConfirmation(fmt.Sprintf("directory %s doesn't exit, create?", dir))
		if err != nil {
			return err
		}
		if confirm {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create dir %q", dir))
			}
		}
	}

	if _, err := os.Stat(fullpath); err == nil || !os.IsNotExist(err) {
		confirm, err := util.PromptExportOverrideConfirmation(fmt.Sprintf("file %s already exits, replace?", fullpath))
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	sp := newSqlPlanner(sq, e.updateRoamDB)
	b, err := e.exportOrgRoam(bk, sp, cfg)
	if err != nil {
		return err
	}
	// Workaround the unicode encoding.
	if err := writeRunesToFile(fullpath, []rune(string(b))); err != nil {
		return err
	}

	if err := sp.InsertFileEntry(bk, fullpath); err != nil {
		return err
	}

	if err := sp.CommitSql(); err != nil {
		return err
	}

	util.Log("Successfully created:", fullpath)

	return nil
}

func (e *OrgRoamExporter) generateOutputPath(b *Book, cfg *config.ConvertConfig) string {
	filename := fmt.Sprintf("《%s》 by %s.org", b.Title, b.Author)
	if e.authorSubDir {
		return filepath.Join(cfg.OutputDir, b.Author, filename)
	}
	return filepath.Join(cfg.OutputDir, filename)
}

func (e *OrgRoamExporter) exportOrgRoam(b *Book, sp SqlPlanner, cfg *config.ConvertConfig) ([]byte, error) {
	var orgTitleTpl, orgEntryTpl string
	if e.templateType < 0 || e.templateType > len(OrgTemplates) {

		orgTitleTpl = OrgTemplates[1].TitleTemplate
		orgEntryTpl = OrgTemplates[1].EntryTemplate
	}
	orgTitleTpl = OrgTemplates[e.templateType].TitleTemplate
	orgEntryTpl = OrgTemplates[e.templateType].EntryTemplate

	b.UUID = uuid.New()
	buf := new(bytes.Buffer)
	titleT := template.Must(template.New("template").Parse(orgTitleTpl))
	if err := titleT.Execute(buf, b); err != nil {
		return nil, fmt.Errorf("failed to execute org template for title: %v", err)
	}

	outputPath := e.generateOutputPath(b, cfg)
	if err := sp.InsertNodeLinkTitleEntry(b, outputPath); err != nil {
		return nil, err
	}

	for _, mk := range b.Marks {
		mk.UUID = uuid.New()
		mk.Pos = len([]rune(buf.String())) + len("\n*")

		if err := sp.InsertNodeLinkMarkEntry(b, &mk, outputPath); err != nil {
			return nil, err
		}

		entryT := template.Must(template.New("template").Parse(orgEntryTpl))
		if err := entryT.Execute(buf, mk); err != nil {
			return nil, fmt.Errorf("failed to execute org template for entries: %v", err)
		}
	}

	return buf.Bytes(), nil
}
