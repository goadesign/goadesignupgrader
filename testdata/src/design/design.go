package design

import (
	. "github.com/goadesign/goa/design"        // want `"github.com/goadesign/goa/design" should be removed`
	. "github.com/goadesign/goa/design/apidsl" // want `"github.com/goadesign/goa/design/apidsl" should be replaced with "goa.design/goa/v3/dsl"`
)

var UserMedia = MediaType("application/vnd.user+json", func() { // want `MediaType should be replaced with ResultType`
	Attribute("id", Integer) // want `Integer should be replaced with Int`
})

var _ = Resource("user", func() { // want `Resource should be replaced with Service`
	Action("show", func() { // want `Action should be replaced with Method`
		Routing(GET("/users/:user_id")) // want `colons in HTTP routing DSLs should be replaced with curly braces`
	})
})
