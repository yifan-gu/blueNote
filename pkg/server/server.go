/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/yifan-gu/blueNote/pkg/config"
	"github.com/yifan-gu/blueNote/pkg/model"
	"github.com/yifan-gu/blueNote/pkg/storage"
	"github.com/yifan-gu/blueNote/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
)

type server struct {
	config *config.ServerConfig
	store  storage.Storage
}

func NewServer(config *config.ServerConfig, store storage.Storage) Server {
	return &server{config: config, store: store}
}

func (s *server) Run() {
	schema = s.graphqlSchema()
	http.HandleFunc("/marks", handleGraphqlMarks)
	util.Logf("Server is running on %v\n", s.config.ListenAddr)
	http.ListenAndServe(s.config.ListenAddr, nil)

}

func handleGraphqlMarks(w http.ResponseWriter, r *http.Request) {
	result := executeQuery(r.Context(), r.URL.Query().Get("query"))
	json.NewEncoder(w).Encode(result)
}

func (s *server) resolveMarksQuery(p graphql.ResolveParams) (interface{}, error) {
	// TODO(yifan): Maybe use a more generic filter type than bson.M{} ?
	filter := bson.M{}
	andCondition := []bson.M{}

	limit, _ := p.Args["limit"].(int)
	id, idOK := p.Args["id"].(string)
	if idOK {
		filter["_id"] = id
	}
	typ, typOK := p.Args["type"].(string)
	if typOK {
		filter["type"] = typ
	}
	title, titleOK := p.Args["title"].(string)
	if titleOK {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	author, authorOK := p.Args["author"].(string)
	if authorOK {
		filter["author"] = bson.M{"$regex": author, "$options": "i"}
	}
	data, dataOK := p.Args["data"].(string)
	if dataOK {
		filter["data"] = bson.M{"$regex": data, "$options": "i"}
	}
	note, noteOK := p.Args["note"].(string)
	if noteOK {
		filter["note"] = bson.M{"$regex": note, "$options": "i"}
	}
	tags, tagsOK := p.Args["tags"].([]interface{})
	if tagsOK {
		for _, tag := range tags {
			tagVal, ok := tag.(string)
			if !ok {
				return nil, errors.New(fmt.Sprintf("Expect []string for tags, but got []%T", tag))
			}
			andCondition = append(andCondition, bson.M{"tags": bson.M{"$regex": tagVal, "$options": "i"}})
		}
	}
	createdBefore, createdBeforeOK := p.Args["createdBefore"].(int)
	if createdBeforeOK {
		andCondition = append(andCondition, bson.M{"createdAt": bson.M{"$lt": createdBefore}})
	}
	createdAfter, createdAfterOK := p.Args["createdAfter"].(int)
	if createdAfterOK {
		andCondition = append(andCondition, bson.M{"createdAt": bson.M{"$gt": createdAfter}})
	}
	lastModifiedBefore, lastModifiedBeforeOK := p.Args["lastModifiedBefore"].(int)
	if lastModifiedBeforeOK {
		andCondition = append(andCondition, bson.M{"lastModifiedAt": bson.M{"$lt": lastModifiedBefore}})
	}
	lastModifiedAfter, lastModifiedAfterOK := p.Args["lastModifiedAfter"].(int)
	if lastModifiedAfterOK {
		andCondition = append(andCondition, bson.M{"lastModifiedAt": bson.M{"$gt": lastModifiedAfter}})
	}

	if len(andCondition) > 0 {
		filter["$and"] = andCondition
	}
	return s.store.GetMarks(p.Context, filter, limit)
}

func (s *server) createOneMark(p graphql.ResolveParams) (interface{}, error) {
	mark := &model.Mark{
		Type:   p.Args["type"].(string),
		Title:  p.Args["title"].(string),
		Author: p.Args["author"].(string),
	}
	section, sectionOK := p.Args["section"]
	if sectionOK {
		mark.Section = section.(string)
	}
	data, dataOK := p.Args["data"]
	if dataOK {
		mark.Data = data.(string)
	}
	note, noteOK := p.Args["note"]
	if noteOK {
		mark.UserNote = note.(string)
	}
	tags, tagsOK := p.Args["tags"].([]interface{})
	if tagsOK {
		for i := range tags {
			mark.Tags = append(mark.Tags, tags[i].(string))
		}
	}
	location, locationOK := p.Args["location"].(map[string]interface{})
	if locationOK {
		createLocationField(mark, location)
	}

	if err := model.ValidateMark(mark); err != nil {
		return nil, err
	}
	id, err := s.store.CreateMark(p.Context, mark)
	if err != nil {
		return nil, err
	}
	mark.ID = id
	return mark, nil
}

func (s *server) updateOneMarkByID(p graphql.ResolveParams) (interface{}, error) {
	id, idOK := p.Args["id"].(string)
	if !idOK {
		return nil, errors.New("No id is given")
	}
	marks, err := s.store.GetMarks(p.Context, bson.M{"_id": id}, 0)
	if err != nil {
		return nil, err
	}
	if len(marks) != 1 {
		return nil, errors.New(fmt.Sprintf("Expect 1 mark, got %d", len(marks)))
	}
	update := marks[0]
	typ, typOK := p.Args["type"].(string)
	if typOK {
		update.Type = typ
	}
	title, titleOK := p.Args["title"].(string)
	if titleOK {
		update.Title = title
	}
	author, authorOK := p.Args["author"].(string)
	if authorOK {
		update.Author = author
	}
	section, sectionOK := p.Args["section"]
	if sectionOK {
		update.Section = section.(string)
	}
	data, dataOK := p.Args["data"].(string)
	if dataOK {
		update.Data = data
	}
	note, noteOK := p.Args["note"].(string)
	if noteOK {
		update.UserNote = note
	}
	location, locationOK := p.Args["location"].(map[string]interface{})
	if locationOK {
		createLocationField(update, location)
	}
	tags, tagsOK := p.Args["tags"].([]interface{})
	if tagsOK {
		update.Tags = nil
		for i := range tags {
			update.Tags = append(update.Tags, tags[i].(string))
		}
	}

	if err := model.ValidateMark(update); err != nil {
		return nil, err
	}

	if err := s.store.UpdateOneMark(p.Context, id, update); err != nil {
		return nil, err
	}
	return update, nil
}

func createLocationField(mark *model.Mark, location map[string]interface{}) {
	mark.Location = &model.Location{}
	for k, v := range location {
		switch k {
		case "chapter":
			mark.Location.Chapter = v.(string)
		case "page":
			page := v.(int)
			mark.Location.Page = &page
		case "location":
			location := v.(int)
			mark.Location.Location = &location
		}
	}
}

func (s *server) deleteOneMarkByID(p graphql.ResolveParams) (interface{}, error) {
	id, idOK := p.Args["id"].(string)
	if !idOK {
		return nil, errors.New("No id is given")
	}

	marks, err := s.store.GetMarks(p.Context, bson.M{"_id": id}, 0)
	if err != nil {
		return nil, err
	}
	if len(marks) != 1 {
		return nil, errors.New(fmt.Sprintf("Expect 1 mark, got %d", len(marks)))
	}
	if err := s.store.DeleteOneMark(p.Context, id); err != nil {
		return nil, err
	}
	return marks[0], nil
}
