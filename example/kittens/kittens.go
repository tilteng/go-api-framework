package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tilteng/go-api-framework/api_framework"
	"github.com/tilteng/go-errors/errors"
)

// Our controller! It embeds a TiltController.
type KittensController struct {
	*api_framework.TiltController
}

// Track our created kittens in memory for this example
var kittens = map[string]*Kitten{}

// ErrorClasses
var ErrInvalidKittenID = &errors.ErrorClass{
	Name:    "ErrInvalidKittenID",
	Code:    "ERR_ID_INVALID_KITTEN_ID",
	Status:  400,
	Message: "Invalid kitten id specified",
}
var ErrKittenNotFound = &errors.ErrorClass{
	Name:    "ErrKittenNotFound",
	Code:    "ERR_ID_KITTEN_NOT_FOUND",
	Status:  404,
	Message: "No kitten found with that id",
}

// Used for deserializing POST data. jsonapi spec says you should use
// { "data": { "attributes": { ... } } }
type createKittenBody struct {
	Data kittenData `json:"data"`
}
type kittenData struct {
	Kitten Kitten `json:"attributes"`
}

// Used for API responses
type Kitten struct {
	Id    *api_framework.UUID `json:"id,omitempty"`
	Name  string              `json:"name"`
	Color string              `json:"color,omitempty"`
}

func (self *KittensController) AddKitten(ctx context.Context) {
	body_obj := &createKittenBody{}
	// ReadBody() is a method on the TiltController struct. It handles
	// deserializing the body into whatever object you pass. If you're
	// using the json schema middleware, the body has already been validated
	// against the schema by this point.
	self.ReadBody(ctx, &body_obj)

	kitten := &body_obj.Data.Kitten
	kitten.Id = self.GenUUID()
	if kitten.Id == nil {
		panic("uuid generation failed")
	}
	kittens[kitten.Id.String()] = kitten

	// WriteResponse() is a method on the TiltController struct. It handles
	// serializing your data according to Accept: header and returing the
	// response. POST routes automatically send back a 201 status code.
	// See GET example below to see how you can return a differnt code.
	self.WriteResponse(ctx, body_obj)
}

func (self *KittensController) GetKitten(ctx context.Context) {
	// RequestContext() is a method on the base go-api-router.Router struct
	// It handles taking a generic context and turning it into the
	// RequestContext struct that was created by go-api-router. Alternatively,
	// go-api-router exports a function 'RequestContextFromContext' that you
	// can use, instead. But I made the convenience method on the Router,
	// because it eliminates the need to import go-api-router. I think I
	// may move the method to the TiltController, though, instead of having
	// it on go-api-router.Router.
	rctx := self.RequestContext(ctx)

	// The route was defined as /kittens/{id}. This is how you pull the
	// ID. I'm considering making a convenience method on TiltController that
	// combines the above call with this call.
	id, _ := rctx.RouteVar("id")

	uuid := self.UUIDFromString(id)
	if uuid == nil {
		self.WriteResponse(
			ctx,
			ErrInvalidKittenID.New(0),
		)
		return
	}
	kitten, ok := kittens[uuid.String()]
	if !ok {
		self.WriteResponse(
			ctx,
			ErrKittenNotFound.New(0),
		)
		return
	}
	self.WriteResponse(ctx, &createKittenBody{
		Data: kittenData{
			Kitten: *kitten,
		},
	})
}

func registerKittens(c *api_framework.TiltController) (err error) {
	kittens := &KittensController{c}
	c.POST("/kittens", kittens.AddKitten,
		// Optional arguments. If you're using the json schema middleware,
		// it will look for a map[string]string argument containing a key
		// of "jsonschema". The value should be the name of a file under
		// the JSONSchemaFilePath (see controller_opts) excluding its
		// .json suffix. When this route is called, the body of data will
		// be validated against the schema found in the json file.
		map[string]string{"jsonschema": "create-kitten"},
	)
	// For more efficient routing, you can create a sub-path
	kittens_router := c.SubRouterForPath("/kittens")
	// ...and then define routes on that subpath.. This is actually
	// "/kittens/{id}"
	kittens_router.GET("/{id}", kittens.GetKitten)
	return
}

func main() {
	port := 31337

	controller_opts := api_framework.NewTiltControllerOpts()
	// BaseAPIURL is used to specify the real externally reachable URL. This
	// is used for returning paths to json schemas via the Link: header
	controller_opts.BaseAPIURL = fmt.Sprintf("http://localhost:%d", port)
	// Directory containing json schema files to load. Must end in .json
	controller_opts.JSONSchemaFilePath = "./schemas"
	// HTTP path where to make json schemas available
	controller_opts.JSONSchemaRoutePath = "/schemas"
	// If set, where output for apache-style logging goes
	controller_opts.ApacheLogWriter = os.Stderr

	controller := api_framework.NewTiltController(controller_opts)

	if err := controller.Init(); err != nil {
		log.Fatal(err)
	}

	if err := registerKittens(controller); err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		controller,
	))
}
