package design

import (
	. "github.com/goadesign/goa/design"        // want `"github.com/goadesign/goa/design" should be removed`
	. "github.com/goadesign/goa/design/apidsl" // want `"github.com/goadesign/goa/design/apidsl" should be replaced with "goa.design/goa/v3/dsl"`
)

var UserMedia = MediaType("application/vnd.user+json", func() { // want `MediaType should be replaced with ResultType`
	Attribute("id", Integer)          // want `Integer should be replaced with Int`
	Attribute("created_at", DateTime) // want `DateTime should be replaced with String`
})

var _ = Resource("user", func() { // want `Resource should be replaced with Service`
	BasePath("/users")      // want `BasePath should be replaced with Path and move it into HTTP`
	Action("show", func() { // want `Action should be replaced with Method`
		Routing(GET("/:user_id")) // want `colons in HTTP routing DSLs should be replaced with curly braces`
	})
})
