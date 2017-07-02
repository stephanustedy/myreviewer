package main

import(
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/paytm/grace.v1"

	"github.com/stephanustedy/myreviewer/src/myreviewer"
)

func main() {
	log.SetFlags(log.Llongfile)

	myreviewer.Initialize()

	router := httprouter.New()

	router.GET("/myreviewer", myreviewer.HomeHandler)
	router.GET("/myreviewer/manage", myreviewer.ManageHandler)

	// Manage Team
	router.POST("/myreviewer/add_team", myreviewer.AddTeamHandler)
	router.POST("/myreviewer/update_team", myreviewer.AddTeamHandler)
	router.POST("/myreviewer/delete_team", myreviewer.DeleteTeamHandler)

	// Manage Member
	router.POST("/myreviewer/add_member", myreviewer.AddMemberHandler)
	router.POST("/myreviewer/delete_member", myreviewer.DeleteMemberHandler)

	// Review
	router.POST("/myreviewer/review", myreviewer.ReviewHandler)
	router.GET("/myreviewer/review/renotify", myreviewer.ReNotifyHandler)
	router.GET("/myreviewer/review/close", myreviewer.CloseReviewHandler)
	router.GET("/myreviewer/review/action", myreviewer.ReviewActionHandler)

	//Assets
	router.ServeFiles("/assets/scripts/*filepath", http.Dir("files/scripts"))
	router.ServeFiles("/assets/css/*filepath", http.Dir("files/css"))

	grace.Serve(":7412", router)
}