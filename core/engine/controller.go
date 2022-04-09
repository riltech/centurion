package engine

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// Starts HTTP serving
func GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	return router
}
