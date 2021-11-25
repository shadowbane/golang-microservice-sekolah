package router

import (
	"github.com/shadowbane/golang-microservice-sekolah/cmd/api/controllers"

	"github.com/julienschmidt/httprouter"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/application"
)

func Get(app *application.Application) *httprouter.Router {
	mux := httprouter.New()

	// index
	mux.GET("/api/v1/school", controllers.SchoolIndex(app))

	// store
	mux.POST("/api/v1/school", controllers.SchoolCreate(app))

	// post
	//mux.POST("/whatsapp", _default.Store(app))
	//mux.GET("/users/:id", getuser.runIndex(app))
	return mux
}
