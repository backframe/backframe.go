package rest

type Generator struct{}

func RestGenerator() Generator {
	// some config should be passed here
	return Generator{}
}

func (g *Generator) Generate() {
	// load some config
	// modify bf.spec.yml
	// inject the routes
	// interface with the bf-dbs pkg
	// inject the models
}

/**
Sample file tree for rest
_________________________
bf.spec.yml
routes.go
server.go
models
	|_ user.schema
	|_ order.schema
resources
	|_ user
		|__ user.model.go
		|__ user.controller.go
		|__ user.router.go(optional)
	|_ order
		|__ order.model.go
		|__ order.controller.go
		|__ order.router.go(optional)

*/
