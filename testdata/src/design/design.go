package design

import (
	. "github.com/goadesign/goa/design"        // want `\A"github.com/goadesign/goa/design" should be removed\z`
	. "github.com/goadesign/goa/design/apidsl" // want `\A"github.com/goadesign/goa/design/apidsl" should be replaced with "goa.design/goa/v3/dsl"\z`
	"net/http"
)

var UserMedia = MediaType("application/vnd.user+json", func() { // want `\AMediaType should be replaced with ResultType\z`
	Attribute("id", Integer)                          // want `\AInteger should be replaced with Int\z`
	Attribute("permissions", HashOf(String, Boolean)) // want `\AHashOf should be replaced with MapOf\z`
	Attribute("created_at", DateTime)                 // want `\ADateTime should be replaced with String \+ Format\(FormatDateTime\)\z`
})

var _ = Resource("user", func() { // want `\AResource should be replaced with Service\z`
	BasePath("/users")      // want `\ABasePath should be replaced with Path and move it into HTTP\z`
	Action("show", func() { // want `\AAction should be replaced with Method\z`
		Routing(GET("/:user_id"))        // want `\ARouting should be replaced with HTTP and colons in HTTP routing DSLs should be replaced with curly braces\z`
		Response(OK, UserMedia, func() { // want `\AResponse should be wrapped by HTTP and OK should be replaced with StatusOK and Status should be replaced with Code\z`
			Status(http.StatusOK)
		})
		Metadata("swagger:summary", "Show users") // want `\AMetadata should be replaced with Meta\z`
	})
	Action("list", func() { // want `\AAction should be replaced with Method\z`
		Routing(GET("/"))                     // want `\ARouting should be replaced with HTTP\z`
		Response(OK, CollectionOf(UserMedia)) // want `\AResponse should be wrapped by HTTP and OK should be replaced with StatusOK\z`
	})
})
