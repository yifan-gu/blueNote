/*
Copyright © 2022 Yifan Gu <guyifan1121@gmail.com>

*/

package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/yifan-gu/blueNote/pkg/util"
)

type Int64 int64

func (i *Int64) Val() int64 {
	return int64(*i)
}

func NewInt64(val string) *Int64 {
	value, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		util.Fatal(fmt.Sprintf("Failed to parse %v: %v", val, err))
	}
	int64Val := Int64(value)
	return &int64Val
}

var int64Type = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64Type",
	Description: "The `Int64Type` scalar type represents an int64 Object.",
	// Serialize serializes the Int64Type
	Serialize: func(value interface{}) interface{} {
		switch val := value.(type) {
		case int64:
			Int64Val := Int64(val)
			return Int64Val.Val()
		case *int64:
			Int64Val := Int64(*val)
			return Int64Val.Val()
		default:
			return nil
		}
	},
	// ParseValue parses GraphQL variables from int64 to `Int64`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case int64:
			return Int64(value)
		case *int64:
			return Int64(*value)
		case string:
			return NewInt64(value)
		case *string:
			return NewInt64(*value)
		default:
			return nil
		}
	},
	// ParseLiteral parses GraphQL AST value to `Int64`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return NewInt64(valueAST.Value)
		default:
			return nil
		}
	},
})

var locationInputType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "LocationInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"chapter": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"page": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"location": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
		},
	},
)

var locationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Location",
		Fields: graphql.Fields{
			"chapter": &graphql.Field{
				Type: graphql.String,
			},
			"page": &graphql.Field{
				Type: graphql.Int,
			},
			"location": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var markType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mark",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: graphql.String,
			},
			"section": &graphql.Field{
				Type: graphql.String,
			},
			"location": &graphql.Field{
				Type: locationType,
			},
			"data": &graphql.Field{
				Type: graphql.String,
			},
			"note": &graphql.Field{
				Type: graphql.String,
			},
			"tags": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"createdAt": &graphql.Field{
				Type: int64Type,
			},
			"lastModifiedAt": &graphql.Field{
				Type: int64Type,
			},
		},
	},
)

func (s *server) graphqlQueryType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				// Get (read) single mark at
				//  http://localhost:11212/marks?query={marks(id:1){title,author,data,note...}}
				"marks": &graphql.Field{
					Type:        graphql.NewList(markType),
					Description: "Get one or more marks",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"type": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"title": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"author": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"data": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"note": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"tags": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
						"location": &graphql.ArgumentConfig{
							Type: locationInputType,
						},
						"createdBefore": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"createdAfter": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"lastModifiedBefore": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"lastModifiedAfter": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
						"limit": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},
					Resolve: s.resolveMarksQuery,
				},
			},
		},
	)
}

func (s *server) graphqlMutationType() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				// Create a new mark
				// http://localhost:11212/marks?query=mutation+_{createOne(type:"",title:"",author:"",data:"",note:"",tags:[]){type,title,author,data,note,tags}}
				"createOne": &graphql.Field{
					Type:        markType,
					Description: "Create a new mark",
					Args: graphql.FieldConfigArgument{
						"type": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String), // TODO(yifan): Use Enum?
						},
						"title": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"author": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"section": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"location": &graphql.ArgumentConfig{
							Type: locationInputType,
						},
						"data": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"note": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"tags": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: s.createOneMark,
				},
				// Update a mark by id
				// http://localhost:11212/marks?query=mutation+_{updateOne(id:1,type:"",title:"",author:"",data:"",note:"",tags:[]){type,title,author,data,note,tags}}
				"updateOne": &graphql.Field{
					Type:        markType,
					Description: "Update a mark by its ID",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						"type": &graphql.ArgumentConfig{
							Type: graphql.String, // TODO(yifan): Use Enum?
						},
						"title": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"author": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"section": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"location": &graphql.ArgumentConfig{
							Type: locationInputType,
						},
						"data": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"note": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"tags": &graphql.ArgumentConfig{
							Type: graphql.NewList(graphql.String),
						},
					},
					Resolve: s.updateOneMarkByID,
				},
				// Delete a mark by id
				// http://localhost:11212/marks?query=mutation+_{delete(id:1,){type,title,author,data,note,tags}}
				"deleteOne": &graphql.Field{
					Type:        markType,
					Description: "Delete a mark by its ID",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: s.deleteOneMarkByID,
				},
			},
		},
	)

}

func (s *server) graphqlSchema() graphql.Schema {
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    s.graphqlQueryType(),
			Mutation: s.graphqlMutationType(),
		},
	)
	if err != nil {
		util.Fatal(err)
	}
	return schema
}

var schema graphql.Schema

func executeQuery(ctx context.Context, query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})
	if len(result.Errors) > 0 {
		util.Logf("Errors running graphql query: %v\n", result.Errors)
	}
	return result
}
