package myreviewer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func postToSlack(text, snippet, channel, webhook string) error {
	payload := map[string]interface{}{
		"text":       text,
		"link_names": 1,
		"attachments": []map[string]interface{}{
			map[string]interface{}{
				"text": snippet,
			},
		},
	}

	payload["channel"] = channel
	payload["username"] = "review-bot"
	payload["icon_emoji"] = ":smiling_imp:"

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
	return "<http://" + host + "/myreviewer/review/action?reviewer_id=" + reviewerId + "&act=" + REVIEWER_STATUS_APPROVE + "&id=" + review.ReviewId + "|Approved>"
}

func generateNotApproveUrl(review *Review, host, reviewerId string) string {
	return "<http://" + host + "/myreviewer/review/action?reviewer_id=" + reviewerId + "&act=" + REVIEWER_STATUS_NOT_APPROVE + "&id=" + review.ReviewId + "|Not Approved>"
}

func generateBusyUrl(review *Review, host, reviewerId string) string {
	return "<http://" + host + "/myreviewer/review/action?reviewer_id=" + reviewerId + "&act=" + REVIEWER_STATUS_BUSY + "&id=" + review.ReviewId + "|Busy>"
}

func notifyReview(review *Review, host, specificReviewerId string) error {
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
		if specificReviewerId == "" || v.ReviewerId == specificReviewerId {
			reviewerIds = append(reviewerIds, v.ReviewerId)
			reviewerUsername = append(reviewerUsername, v.ReviewerUsername)
			reviewerName = append(reviewerName, v.ReviewerName)
		}
	}

	text := "Review in dong kak " + strings.Join(reviewerUsername, " ")
	snippet := "Pull Request : " + review.PullRequest + "\nReviewer:\n"
	for k, v := range reviewerName {
		snippet = snippet + strconv.Itoa(k+1) + ". " + v + " " + generateApproveUrl(review, host, reviewerIds[k]) + "/" + generateNotApproveUrl(review, host, reviewerIds[k]) + "/" + generateBusyUrl(review, host, reviewerIds[k]) + "\n"
	}

	for {
		err = postToSlack(text, snippet, team.Channel, team.Webhook)
		if err == nil {
			break
		}
		log.Println(err)
	}

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

	for {
		err = postToSlack(text, snippet, team.Channel, team.Webhook)
		if err == nil {
			break
		}
		log.Println(err)
	}

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

	for {
		err = postToSlack(text, snippet, team.Channel, team.Webhook)
		if err == nil {
			break
		}
		log.Println(err)
	}

	return err
}
