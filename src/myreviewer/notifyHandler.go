package myreviewer

import(
	"strings"
	"strconv"
	"net/http"
	"bytes"
	"log"
	"encoding/json"
	"io/ioutil"
	"time"
)

func notifyReview(review *Review, host string) error {
	team, err := getTeamById(review.TeamId)
	if err != nil {
		log.Println("error on get team by id", err)
		return err
	}

	reviewerUsername := []string{}
	reviewerId := []string{}
	reviewerName := []string{}
	for _, v := range review.Reviewer {
		member, err := getMemberById(v.ReviewerId)
		if err != nil {
			return err
		}
		reviewerUsername = append(reviewerUsername, member.Username)
		reviewerId = append(reviewerId, member.MemberId)
		reviewerName = append(reviewerName, member.Name)
	}

	text := "Review in dong kak " + strings.Join(reviewerUsername, " ")
	snippet := "Pull Request : " + review.PullRequest + "\nReviewer:\n"
	for k, v := range reviewerName {
		snippet = snippet + strconv.Itoa(k+1) + ". " + v + " " +generateApproveUrl(review, host, reviewerId[k]) + "/" + generateNotApproveUrl(review, host, reviewerId[k]) + "\n"
	}

	err = postToSlack(text, snippet, team.Channel, team.Webhook)
	return err
}

func postToSlack(text, snippet, channel, webhook string) error {
	payload := map[string]interface{}{
		"text": text,
		"link_names": 1,
		"attachments": []map[string]interface{}{
			map[string]interface{}{
				"text": snippet,
			},
		},
	}

	payload["channel"] = channel
	
	b, err := json.Marshal(payload)
	if err != nil {
		log.Println("[panics] marshal err", err)
		return err
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Post(webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("[panics] error on capturing error : %s\n", err.Error())
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[panics] error on capturing error : %s\n", err)
			return err
		}
		log.Printf("[panics] error on capturing error : %s\n", string(b))
		return err
	}
	return nil
}

func generateApproveUrl(review *Review, host, reviewerId string) string {
	return "<http://" + host + "/myreviewer/review/action?reviewer_id="+reviewerId+"&act=2&id=" + review.ReviewId+ "|Approved>"
}

func generateNotApproveUrl(review *Review, host, reviewerId string) string {
	return "<http://" + host + "/myreviewer/review/action?reviewer_id="+reviewerId+"&act=-1&id=" + review.ReviewId+ "|Not Approved>"
}

func reNotifyReview(review *Review, host string) error {
	team, err := getTeamById(review.TeamId)
	if err != nil {
		log.Println("error on get team by id", err)
		return err
	}

	pendingReviewer, err := getPendingReviewerByReviewId(review.ReviewId)
	if err != nil {
		log.Println(err)
		return err
	}

	reviewerUsername := []string{}
	reviewerIds := []string{}
	reviewerName := []string{}
	for _, v := range pendingReviewer {
		reviewerIds = append(reviewerIds, v.ReviewerId)
		reviewerUsername = append(reviewerUsername, v.ReviewerUsername)
		reviewerName = append(reviewerName, v.ReviewerName)
	}

	text := "Review in dong kak " + strings.Join(reviewerUsername, " ")
	snippet := "Pull Request : " + review.PullRequest + "\nReviewer:\n"
	for k, v := range reviewerName {
		snippet = snippet + strconv.Itoa(k+1) + ". " + v + " " +generateApproveUrl(review, host, reviewerIds[k]) + "/" + generateNotApproveUrl(review, host, reviewerIds[k]) + "\n"
	}

	err = postToSlack(text, snippet, team.Channel, team.Webhook)
	return err
}

func notifyApproved(review *Review) error {
	developer, err := getMemberById(review.Developer)
	if err != nil {
		return err
	}
	qa, err := getMemberById(review.QA)
	if err != nil {
		return err
	}
	team, err := getTeamById(review.TeamId)
	if err != nil {
		return err
	}

	text := "udah siap deploy nih kak " + qa.Username + " " + developer.Username
	snippet := "Pull Request : " + review.PullRequest

	err = postToSlack(text, snippet, team.Channel, team.Webhook)
	return err
}

func notifyNotApproved(review *Review) error {
	developer, err := getMemberById(review.Developer)
	if err != nil {
		return err
	}
	qa, err := getMemberById(review.QA)
	if err != nil {
		return err
	}
	team, err := getTeamById(review.TeamId)
	if err != nil {
		return err
	}

	text := "perlu ada yg di ubah nih kak " + qa.Username + " " + developer.Username
	snippet := "Pull Request : " + review.PullRequest

	err = postToSlack(text, snippet, team.Channel, team.Webhook)
	return err
}
