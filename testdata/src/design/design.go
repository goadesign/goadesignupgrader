package design

import ( // want `\Aimport declarations should be fixed\z`
	. "github.com/goadesign/goa/design"        // want `\A"github.com/goadesign/goa/design" should be removed\z`
	. "github.com/goadesign/goa/design/apidsl" // want `\A"github.com/goadesign/goa/design/apidsl" should be replaced with "goa.design/goa/v3/dsl"\z`
	"net/http"
)

var _ = API("api", func() { // want `\Avariable declarations should be fixed\z`
	BasePath("/:version")                           // want `\ABasePath should be replaced with Path and wrapped by HTTP\z`
	Consumes("application/json", "application/xml") // want `\AConsumes should be wrapped by HTTP\z`
	Produces("application/json", "application/xml") // want `\AProduces should be wrapped by HTTP\z`
	Params(func() {                                 // want `\AParams should be wrapped by HTTP\z`
		Param("version")
	})
})

var User = Type("user", func() { // want `\Avariable declarations should be fixed\z`
	Attribute("permissions", HashOf(String, Boolean)) // want `\AHashOf should be replaced with MapOf\z`
})

var UserMedia = MediaType("application/vnd.user+json", func() { // want `\Avariable declarations should be fixed\z` `\AMediaType should be replaced with ResultType\z`
	Attribute("id", Integer)                          // want `\AInteger should be replaced with Int\z`
	Attribute("permissions", HashOf(String, Boolean)) // want `\AHashOf should be replaced with MapOf\z`
	Attribute("interests", HashOf(String, Integer,    // want `\AHashOf should be replaced with MapOf\z` `\AInteger should be replaced with Int\z`
		func() { // want `\Aoptional DSL for key of HashOf should be set by Key\z`
			MinLength(1)
			MaxLength(16)
		}, func() { // want `\Aoptional DSL for value of HashOf should be set by Elem\z`
			Minimum(1)
			Maximum(5)
		},
	))
	Attribute("created_at", DateTime) // want `\ADateTime should be replaced with String \+ Format\(FormatDateTime\)\z`
})

var _ = Resource("user", func() { // want `\Avariable declarations should be fixed\z` `\AResource should be replaced with Service\z`
	BasePath("/users")          // want `\ABasePath should be replaced with Path and wrapped by HTTP\z`
	CanonicalActionName("show") // want `\ACanonicalActionName should be replaced with CanonicalMethod and wrapped by HTTP\z`
	DefaultMedia(UserMedia)     // want `\ADefaultMedia should be removed\z`
	Headers(func() {            // want `\AHeaders should be wrapped by HTTP\z`
		Header("Time-Zone")
	})
	Params(func() { // want `\AParams should be wrapped by HTTP\z`
		Param("token")
	})
	Action("show", func() { // want `\AAction should be replaced with Method\z`
		Routing(GET("/:user_id")) // want `\ARouting should be replaced with HTTP\z` `\Acolons in HTTP routing DSLs should be replaced with curly braces\z`
		Headers(func() {          // want `\AHeaders should be wrapped by HTTP\z`
			Header("Link")
		})
		Response(OK, func() { // want `\AResponse should be wrapped by HTTP\z` `\AOK should be replaced with StatusOK\z`
			Media(UserMedia)      // want `\AMedia for a non-error response should be replaced with Result and wrapped by HTTP in the parent\z`
			Status(http.StatusOK) // want `\AStatus should be replaced with Code\z`
		})
		Response(NotFound, func() { // want `\AResponse should be wrapped by HTTP\z` `\ANotFound should be replaced with StatusNotFound\z`
			Media(ErrorMedia)           // want `\AMedia for an error response should be removed\z`
			Status(http.StatusNotFound) // want `\AStatus should be replaced with Code\z`
		})
		Metadata("swagger:summary", "Show users") // want `\AMetadata should be replaced with Meta\z`
	})
	Action("list", func() { // want `\AAction should be replaced with Method\z`
		Routing(GET("/")) // want `\ARouting should be replaced with HTTP\z`
		Params(func() {   // want `\AParams should be wrapped by HTTP\z`
			Param("page")
		})
		Response(OK, CollectionOf(UserMedia)) // want `\AResponse should be wrapped by HTTP\z` `\AOK should be replaced with StatusOK\z`
		Response(BadRequest, ErrorMedia)      // want `\AResponse should be wrapped by HTTP\z` `\ABadRequest should be replaced with StatusBadRequest\z` `\AErrorMedia should be removed\z`
	})
	Action("create", func() { // want `\AAction should be replaced with Method\z`
		Routing(POST("/")) // want `\ARouting should be replaced with HTTP\z`
		Payload(User)
		Response(Created, UserMedia, func() { // want `\AResponse should be wrapped by HTTP\z` `\ACreated should be replaced with StatusCreated\z`
			Headers(func() { // No need to fix Header in Response.
				Header("Location")
			})
		})
	})
})

var _ = Resource("post", func() { // want `\Avariable declarations should be fixed\z` `\AResource should be replaced with Service\z`
	Parent("user") // want `\AParent should be wrapped by HTTP\z`
})
