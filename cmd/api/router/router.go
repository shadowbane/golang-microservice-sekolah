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

	// show
	mux.GET("/api/v1/school/:id", controllers.SchoolShow(app))

	// store
	mux.POST("/api/v1/school", controllers.SchoolCreate(app))

	// update
	mux.PUT("/api/v1/school/:id", controllers.SchoolUpdate(app))

	// delete
	mux.DELETE("/api/v1/school/:id", controllers.SchoolDelete(app))

	// post
	//mux.POST("/whatsapp", _default.Store(app))
	//mux.GET("/users/:id", getuser.runIndex(app))
	return mux
}
