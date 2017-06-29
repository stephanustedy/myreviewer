package main

import(
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/paytm/grace.v1"

	"github.com/stephanustedy/myreviewer/src/myreviewer"
)

func main() {
	myreviewer.Initialize()

	router := httprouter.New()

	router.GET("/myreviewer", myreviewer.ReviewHandler)
	router.GET("/myreviewer/manage", myreviewer.ManageHandler)

	router.POST("/myreviewer/add", myreviewer.AddTeamHandler)

	//Assets
	router.ServeFiles("/assets/scripts/*filepath", http.Dir("files/scripts"))
	router.ServeFiles("/assets/css/*filepath", http.Dir("files/css"))

	grace.Serve(":7412", router)
}