# [gqlgen](https://github.com/99designs/gqlgen)

[Getting Started](https://gqlgen.com/getting-started/)

`schema.graphqls` -> generate `graph/model/*` files + `graph/schema.resolvers.go` -> implement methods in `graph/schema.resolvers.go`.


1. Create `tools.go`

    ```
    printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go
    ```

1. Create the project skeleton

    ```
    go run github.com/99designs/gqlgen init
    ```

    <details>

    ```
    Creating gqlgen.yml
    Creating graph/schema.graphqls
    Creating server.go
    Generating...

    Exec "go run ./server.go" to start GraphQL server
    ```

    </details>

    Files are generated:

    ```
    .
    ├── README.md
    ├── gqlgen.yml
    ├── graph
    │   ├── generated.go
    │   ├── model
    │   │   └── models_gen.go
    │   ├── resolver.go
    │   ├── schema.graphqls
    │   └── schema.resolvers.go
    ├── server.go
    └── tools.go

    2 directories, 9 files
    ```

1. Define your schema (update the generated [graph/schema.graphqls](graph/schema.graphqls))

    > gqlgen is a schema-first library

    <details><summary>graphql/schema.graphqls</summary>

    ```gql
    # GraphQL schema example
    #
    # https://gqlgen.com/getting-started/

    type Todo {
      id: ID!
      text: String!
      done: Boolean!
      user: User!
    }

    type User {
      id: ID!
      name: String!
    }

    type Query {
      todos: [Todo!]!
    }

    input NewTodo {
      text: String!
      userId: String!
    }

    type Mutation {
      createTodo(input: NewTodo!): Todo!
    }
    ```

    </details>

1. Implement your resolver `graph/schema.resolvers.go`

    `graph/resolver.go`: generated with empty struct by gqlgen. track our state.

    ```go
    type Resolver struct {
        todos []*model.ToDo
    }
    ```

    `graph/schema.resolvers.go`: generated with `panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))` by gqlgen.

    ```go
    func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
        todo := &model.Todo{
            Text: input.Text,
            ID:   fmt.Sprintf("T%d", rand.Int()),
            User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
        }
        r.todos = append(r.todos, todo)
        return todo, nil
    }

    func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
        return r.todos, nil
    }
    ```


1. Run server

    ```
    go run ./server.go
    ```

1. localhost:8080

    createTodo

    ```
    mutation createTodo {
      createTodo(input: { text: "todo", userId: "1" }) {
        user {
          id
        }
        text
        done
      }
    }
    ```

    findTodos

    ```
    query findTodos {
      todos {
        text
        done
        user {
          name
        }
      }
    }
    ```

1. More config with `gqlgen.yml` (enable to use your custom model)

    1. `autobind`: allow gqlgen to use your custom models if it can find them rather than generating them.
        ```yaml
        autobind:
          - "github.com/[username]/gqlgen-todos/graph/model"
        ```
    1. add Todo fields resolver
        ```
        models:
          ...
          Todo:
            fields:
              user:
                resolver: true
        ```

    1. Create a new model file [graph/model/todo.go](graph/model/todo.go)

        ```go
        package model

        type Todo struct {
        	ID     string `json:"id"`
        	Text   string `json:"text"`
        	Done   bool   `json:"done"`
        	UserID string `json:"userId"`
        	User   *User  `json:"user"`
        }
        ```

    1. Run `go run github.com/99designs/gqlgen generate`
        `Todo` will be removed from `models_gen.go`
    1. Update `CreateTodo` in [graph/schema.resolvers.go](graph/schema.resolvers.go)

        Add the following line to `&model.Todo{}`

        ```go
            UserID: input.UserID,
        ```
1. Enable to update codes

    Put the following code between package and import in `resourver.go`
    ```go
    //go:generate go run github.com/99designs/gqlgen generate
    ```

    You can now regenerate the code with the following command:

    ```
    go generate ./...
    ```

    This is same as `go run github.com/99designs/gqlgen generate`
