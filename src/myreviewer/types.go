package myreviewer

import(
	"time"
)

type (
	Team struct {
		TeamId 	string
		Name 	string
		Webhook string
		Channel string
		Member 	[]*Member
	}

	Member struct {
		MemberId 	string
		TeamId		string
		Name 		string
		Username 	string
		Role 		string
	}

	Review struct {
		ReviewId 		string
		TeamId 			string
		QA 				string // QA user id
		Developer		string // SE user id
		PullRequest		string // Pull request
		Status 			int
		Reviewer 		[]*Reviewer

		TeamName 		string
	}

	Reviewer struct {
		Id 					string
		ReviewerId 			string
		LastNotify 			time.Time
		Status 				int
		NotifyCount 		int
		ReviewerName 		string
		ReviewerUsername 	string
		StatusStr 			string
	}
)