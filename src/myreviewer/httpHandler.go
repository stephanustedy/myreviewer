package myreviewer

import(
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/stephanustedy/myreviewer/src/util"
)

func HomeHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	message := req.FormValue("message")
	teams, err := getTeamList()
	if err != nil {
		log.Println(err)
	}
	reviews, err := getAllActiveReview()
	if err != nil {
		log.Println(err)
	}
	content := reviewTemplateHandler(message, teams, reviews)
	util.RenderHTML(w, content, 200)
}

func ManageHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	message := req.FormValue("message")

	teams, err := getTeamList()
	if err != nil {
		log.Println(err)
	}
	content := manageTemplateHandler(message, teams)
	util.RenderHTML(w, content, 200)
}

func AddTeamHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	name := req.PostFormValue("name")
	webhook := req.PostFormValue("webhook")
	channel := req.PostFormValue("channel")

	if name != "" && webhook != "" && channel != "" {

	}

	team := &Team{
		Name : name,
		Webhook: webhook,
		Channel : channel,
	}

	message := "success add team!"

	err := addTeam(team)
	if err != nil {
		log.Println(err)
		message = "failed add team!"
	}
	http.Redirect(w, req, "/myreviewer/manage?message=" + message, http.StatusSeeOther)
}

func DeleteTeamHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	teamId := req.PostFormValue("team_id")

	message := "success delete team!"

	err := deleteTeam(teamId)
	if err != nil {
		log.Println(err)
		message = "failed delete team!"
	}
	http.Redirect(w, req, "/myreviewer/manage?message=" + message, http.StatusSeeOther)
}

func AddMemberHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	team := req.PostFormValue("team")
	role := req.PostFormValue("role")
	name := req.PostFormValue("name")
	username := req.PostFormValue("username")

	member := &Member{
		TeamId : team,
		Role : role,
		Name : name,
		Username : username,
	}

	message := "success add member!"

	err := addMember(member)
	if err != nil {
		log.Println(err)
		message = "failed add member!"
	}
	http.Redirect(w, req, "/myreviewer/manage?message=" + message, http.StatusSeeOther)
}

func DeleteMemberHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	memberId := req.PostFormValue("member_id")

	message := "success delete member!"

	err := deleteMember(memberId)
	if err != nil {
		log.Println(err)
		message = "failed delete member!"
	}
	http.Redirect(w, req, "/myreviewer/manage?message=" + message, http.StatusSeeOther)
}

func ReviewHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	qa := req.PostFormValue("qa")
	se := req.PostFormValue("se")
	team := req.PostFormValue("team")
	pr := req.PostFormValue("pr")
	host := req.Host

	review := &Review {
		TeamId : team,
		QA : qa,
		Developer : se,
		PullRequest : pr,
	}
	message := "success notify review!"
	review, err := addReview(review)
	if err != nil {
		log.Println("failed add review", err)
		message = "failed notify review!"
	}

	err = notifyReview(review, host)
	if err != nil {
		log.Println("failed notify review", err)
		message = "failed notify review!"
	}

	http.Redirect(w, req, "/myreviewer?message=" + message, http.StatusSeeOther)
}

func ReNotifyHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reviewId := req.FormValue("id")
	host := req.Host

	review, err := getReviewById(reviewId)
	if err != nil {
		log.Println(err)
	}
	message := "success renotify review!"
	err = reNotifyReview(review, host)
	if err != nil {
		log.Println("failed notify review", err)
		message = "failed notify review!"
	}

	http.Redirect(w, req, "/myreviewer?message=" + message, http.StatusSeeOther)
}

func CloseReviewHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reviewId := req.FormValue("id")

	message := "success close review!"

	err := updateReview(reviewId, 0)
	if err != nil {
		log.Println(err)
		message = "failed close review!"
	}
	http.Redirect(w, req, "/myreviewer?message=" + message, http.StatusSeeOther)
}

func ReviewActionHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reviewId := req.FormValue("id")
	act := req.FormValue("act")
	reviewerId := req.FormValue("reviewer_id")


	err := updateReviewer(reviewId, act, reviewerId)
	if err != nil {
		log.Println(err)
	}

	count := countPendingReviewer(reviewId)
	if count == 0 {
		isApproved := getNotApproveCount(reviewId)
		review, err := getReviewById(reviewId)
		if err != nil {
			log.Println(err)
		}
		if isApproved {
			err := notifyApproved(review)
			if err != nil {
				log.Println(err)
			}
		} else {
			err := notifyNotApproved(review)
			if err != nil {
				log.Println(err)
			}
		}
		// set review to done
		updateReview(reviewId, 2)
	}

	content := closeTemplateHandler()
	util.RenderHTML(w, content, 200)
}