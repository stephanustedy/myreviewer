package myreviewer

import(
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stephanustedy/myreviewer/src/util"
)

func ReviewHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	content := homeTemplateHandler()
	util.RenderHTML(w, content, 200)
}

func ManageHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	content := manageTemplateHandler()
	util.RenderHTML(w, content, 200)
}

func AddTeamHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	name := req.PostFormValue("team_addteam")
	webhook := req.PostFormValue("slack_webhook_addteam")
	channel := req.PostFormValue("slack_channel_addteam")

	if name != "" && webhook != "" && channel != "" {

	}

	team := &Team{
		Name : name,
		Webhook: webhook,
		Channel : channel,
	}

	err := addTeam(team)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("success")
	}

	content := manageTemplateHandler()
	util.RenderHTML(w, content, 200)
}