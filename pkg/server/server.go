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
	typ, typOK := p.Args["type"].(string)
	title, titleOK := p.Args["title"].(string)
	author, authorOK := p.Args["author"].(string)
	data, dataOK := p.Args["data"].(string)
	note, noteOK := p.Args["note"].(string)
	tags, tagsOK := p.Args["tags"].([]interface{})
	createdBefore, createdBeforeOK := p.Args["createdBefore"].(int)
	createdAfter, createdAfterOK := p.Args["createdAfter"].(int)
	lastModifiedBefore, lastModifiedBeforeOK := p.Args["lastModifiedBefore"].(int)
	lastModifiedAfter, lastModifiedAfterOK := p.Args["lastModifiedAfter"].(int)

	if idOK {
		filter["_id"] = id
	}
	if typOK {
		filter["type"] = typ
	}
	if titleOK {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if authorOK {
		filter["author"] = bson.M{"$regex": author, "$options": "i"}
	}
	if dataOK {
		filter["data"] = bson.M{"$regex": data, "$options": "i"}
	}
	if noteOK {
		filter["note"] = bson.M{"$regex": note, "$options": "i"}
	}
	if tagsOK {
		for _, tag := range tags {
			tagVal, ok := tag.(string)
			if !ok {
				return nil, errors.New(fmt.Sprintf("Expect []string for tags, but got []%T", tag))
			}
			andCondition = append(andCondition, bson.M{"tags": bson.M{"$regex": tagVal, "$options": "i"}})
		}
	}
	if createdBeforeOK {
		andCondition = append(andCondition, bson.M{"createdAt": bson.M{"$lt": createdBefore}})
	}
	if createdAfterOK {
		andCondition = append(andCondition, bson.M{"createdAt": bson.M{"$gt": createdAfter}})
	}
	if lastModifiedBeforeOK {
		andCondition = append(andCondition, bson.M{"lastModifiedAt": bson.M{"$lt": lastModifiedBefore}})
	}
	if lastModifiedAfterOK {
		andCondition = append(andCondition, bson.M{"lastModifiedAt": bson.M{"$gt": lastModifiedAfter}})
	}

	if len(andCondition) > 0 {
		filter["$and"] = andCondition
	}
	return s.store.GetMarks(p.Context, filter, limit)
}
