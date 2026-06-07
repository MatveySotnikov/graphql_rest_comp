package main

import (
	"net/http"

	"task-api-comparison/store"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Тип Task в схеме GraphQL
var taskType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Task",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.ID},
		"title":       &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"done":        &graphql.Field{Type: graphql.Boolean},
	},
})

// Корневые запросы
var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"tasks": &graphql.Field{
			Type: graphql.NewList(taskType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return store.All(), nil
			},
		},
		"task": &graphql.Field{
			Type: taskType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				t, ok := store.ByID(id)
				if !ok {
					return nil, nil // в GraphQL принято возвращать nil при отсутствии
				}
				return t, nil
			},
		},
	},
})

// Входные типы для мутаций
var createTaskInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateTaskInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
		"description": &graphql.InputObjectFieldConfig{Type: graphql.String},
	},
})

var updateTaskInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateTaskInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"done": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Boolean)},
	},
})

// Мутации
var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"createTask": &graphql.Field{
			Type: taskType,
			Args: graphql.FieldConfigArgument{
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createTaskInput)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				input := p.Args["input"].(map[string]interface{})
				title := input["title"].(string)
				desc, _ := input["description"].(string)
				return store.Create(title, desc), nil
			},
		},
		"updateTask": &graphql.Field{
			Type: taskType,
			Args: graphql.FieldConfigArgument{
				"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(updateTaskInput)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Args["id"].(string)
				input := p.Args["input"].(map[string]interface{})
				done := input["done"].(bool)
				t, ok := store.UpdateDone(id, done)
				if !ok {
					return nil, nil
				}
				return t, nil
			},
		},
	},
})

func main() {
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true, // встроенная песочница
	})
	http.Handle("/graphql", h)
	http.ListenAndServe(":8083", nil)
}
