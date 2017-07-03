package myreviewer

import(
	"log"
	"time"
	"strconv"
)

func getTeamList() ([]*Team, error) {
	teams := []*Team{}
	rows, err := db.Query(`
			SELECT
				team_id
				, name
				, coalesce(webhook, '')
				, channel
			FROM
				team
			WHERE
				status = 1
		`)
	if err != nil {
		log.Println(err)
		return teams, err
	}
	defer rows.Close()
	for rows.Next() {
		team := &Team{}
		err := rows.Scan(&team.TeamId, &team.Name, &team.Webhook, &team.Channel)
		if err != nil {
			log.Println(err)
			return teams, err
		}
		members, err := getMemberList(team.TeamId)
		if err != nil {
			log.Println(err)
			return teams, nil
		}
		team.Member = members
		teams = append(teams, team)
	}

	return teams, nil
}

func getTeamById(teamId string) (*Team, error) {
	team := &Team{}
	rows, err := db.Query(`
			SELECT
				team_id
				, name
				, coalesce(webhook, '')
				, channel
			FROM
				team
			WHERE
				status = 1
				AND team_id = ` + teamId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&team.TeamId, &team.Name, &team.Webhook, &team.Channel)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		members, err := getMemberList(team.TeamId)
		if err != nil {
			log.Println(err)
			return nil, nil
		}
		team.Member = members
	}

	return team, nil
}

func getMemberList(teamId string) ([]*Member, error) {
	members := []*Member{}
	rows, err := db.Query(`
			SELECT
				member_id
				, name
				, username
				, role
			FROM
				member
			WHERE
				status = 1
				AND team_id = ` + teamId)
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer rows.Close()
	for rows.Next() {
		member := &Member{}
		var role int
		err := rows.Scan(&member.MemberId, &member.Name, &member.Username, &role)
		if err != nil {
			log.Println(err)
			return members, err
		}
		if role == 2 {
			member.Role = "QA"
		} else {
			member.Role = "SE"
		}
		members = append(members, member)
	}

	return members, nil
}

func getMemberById(memberId string) (*Member, error) {
	member := &Member{}
	rows, err := db.Query(`
			SELECT
				member_id
				, name
				, username
				, role
			FROM
				member
			WHERE 
				status = 1
				AND member_id =` + memberId)
	if err != nil {
		log.Println(err)
		return member, err
	}
	defer rows.Close()
	for rows.Next() {
		var role int
		err := rows.Scan(&member.MemberId, &member.Name, &member.Username, &role)
		if err != nil {
			log.Println(err)
			return member, err
		}
		if role == 2 {
			member.Role = "QA"
		} else {
			member.Role = "SE"
		}
	}

	return member, nil
}

func getSEMemberList(teamId, excludeId string) ([]*Member, error) {
	members := []*Member{}
	rows, err := db.Query(`
			SELECT
				member_id
				, name
				, username
				, role
			FROM
				member
			WHERE
				status = 1
				AND role = 1
				AND member_id != `+excludeId+`
				AND team_id = ` + teamId)
	if err != nil {
		log.Println(err)
		return members, err
	}
	defer rows.Close()
	for rows.Next() {
		member := &Member{}
		var role int
		err := rows.Scan(&member.MemberId, &member.Name, &member.Username, &role)
		if err != nil {
			log.Println(err)
			return members, err
		}
		if role == 2 {
			member.Role = "QA"
		} else {
			member.Role = "SE"
		}
		members = append(members, member)
	}

	return members, nil
}

func deleteTeam(teamId string) error {
	stmt, err := db.Prepare("UPDATE team SET status = 0 where team_id = ?")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(teamId)
    if err != nil {
    	return err
    }

    return nil
}

func deleteMember(memberId string) error {
	stmt, err := db.Prepare("UPDATE member SET status = 0 where member_id = ?")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(memberId)
    if err != nil {
    	return err
    }

    return nil
}

func addMember(member *Member) error {
	stmt, err := db.Prepare("INSERT INTO member(team_id, role, name, username, status) values(?,?,?,?,?)")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(member.TeamId, member.Role, member.Name, member.Username, 1)
    if err != nil {
    	return err
    }

    return nil
}

func addTeam(team *Team) error {
	stmt, err := db.Prepare("INSERT INTO team(name, status, webhook, channel) values(?,?,?,?)")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(team.Name, 1, team.Webhook, team.Channel)
    if err != nil {
    	return err
    }

    return nil
}

func updateTeam() {

}

func updateMember() {

}

func addReview(review *Review) (*Review, error) {
	stmt, err := db.Prepare("INSERT INTO review(team_id, qa, developer, pull_request, status) values(?,?,?,?,?)")
    if err != nil {
    	return nil, err
    }

    res, err := stmt.Exec(review.TeamId, review.QA, review.Developer, review.PullRequest, 1)
    if err != nil {
    	return nil, err
    }

    id, err := res.LastInsertId()
    if err != nil {
    	return nil, err
    }

    review.ReviewId = strconv.FormatInt(id, 10)

    reviewer, err := addReviewer(review)
    if err != nil {
    	return nil, err
    }

    review.Reviewer = reviewer

	return review, nil
}

func addReviewer(review *Review) ([]*Reviewer, error) {
	reviewers := []*Reviewer{}
	team, err  := getTeamById(review.TeamId)
	reviewer := randomMember(team, review.Developer)

	for _, v := range reviewer {
		if err != nil {
			log.Println(err)
			return reviewers, err
		}

		stmt, err := db.Prepare("INSERT INTO reviewer(review_id, reviewer_id, status, last_notify, notify_count) values(?,?,?,?,?)")
	    if err != nil {
	    	return reviewers, err
	    }

	    res, err := stmt.Exec(review.ReviewId, v.MemberId, 1, "now()", 1)
	    if err != nil {
	    	return reviewers, err
	    }

	    id, err := res.LastInsertId()
	    if err != nil {
	    	return reviewers, err
	    }

	    temp := &Reviewer{}
	    temp.Id = strconv.FormatInt(id, 10)
	    temp.ReviewerId = v.MemberId
	    temp.LastNotify = time.Now()
	    temp.NotifyCount = 1
	    temp.Status = 1

	    reviewers = append(reviewers, temp)
	}

	return reviewers, nil
}

func getAllActiveReview() ([]*Review, error) {
	reviews := []*Review{}
	rows, err := db.Query(`
			SELECT
				team.Name
				, review_id
				, pull_request
				, review.team_id
			FROM
				review
				JOIN team ON review.team_id = team.team_id
			WHERE
				review.status = 1`)
	if err != nil {
		log.Println(err)
		return reviews, err
	}
	defer rows.Close()
	for rows.Next() {
		review := &Review{}
		err := rows.Scan(&review.TeamName, &review.ReviewId, &review.PullRequest, &review.TeamId)
		if err != nil {
			log.Println(err)
			return reviews, err
		}
		review.Reviewer, err = getReviewerByReviewId(review.ReviewId)
		if err != nil {
			return reviews, nil
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func getReviewerByReviewId(reviewId string) ([]*Reviewer, error) {
	reviewers := []*Reviewer{}
	rows, err := db.Query(`
			SELECT
				id
				, reviewer.status
				, member.Name
			FROM
				reviewer
				JOIN member on reviewer.reviewer_id = member.member_id
			WHERE
				review_id = ` + reviewId)
	if err != nil {
		log.Println(err)
		return reviewers, err
	}
	defer rows.Close()
	for rows.Next() {
		reviewer := &Reviewer{}
		err := rows.Scan(&reviewer.Id, &reviewer.Status, &reviewer.ReviewerName)
		if err != nil {
			log.Println(err)
			return reviewers, err
		}

		if reviewer.Status == 1 {
			reviewer.StatusStr = "No Response"
		} else if reviewer.Status == 2 {
			reviewer.StatusStr = "Approved"
		} else if reviewer.Status == -1 {
			reviewer.StatusStr = "Not Approved"
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

func getPendingReviewerByReviewId(reviewId string) ([]*Reviewer, error) {
	reviewers := []*Reviewer{}
	rows, err := db.Query(`
			SELECT
				id
				, reviewer.status
				, member.Name
				, member.Username
				, reviewer.reviewer_id
			FROM
				reviewer
				JOIN member on reviewer.reviewer_id = member.member_id
			WHERE
				reviewer.status = 1
				AND review_id = ` + reviewId)
	if err != nil {
		log.Println(err)
		return reviewers, err
	}
	defer rows.Close()
	for rows.Next() {
		reviewer := &Reviewer{}
		var reviewStatus int
		err := rows.Scan(&reviewer.Id, &reviewStatus, &reviewer.ReviewerName, &reviewer.ReviewerUsername, &reviewer.ReviewerId)
		if err != nil {
			log.Println(err)
			return reviewers, err
		}

		if reviewStatus == 1 {
			reviewer.StatusStr = "No Response"
		} else if reviewStatus == 2 {
			reviewer.StatusStr = "Approved"
		} else if reviewStatus == -1 {
			reviewer.StatusStr = "Not Approved"
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

func getReviewById(reviewId string) (*Review, error) {
	review := &Review{}
	rows, err := db.Query(`
			SELECT
				review_id
				, team_id
				, pull_request
				, developer
				, qa
			FROM
				review
			WHERE
				review_id = ` + reviewId)
	if err != nil {
		log.Println(err)
		return review, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&review.ReviewId, &review.TeamId, &review.PullRequest, &review.Developer, &review.QA)
		if err != nil {
			log.Println(err)
			return review, err
		}
	}

	return review, nil
}

func updateReview(reviewId string, status int) error {
	stmt, err := db.Prepare("UPDATE review SET status = ? where review_id = ?")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(status, reviewId)
    if err != nil {
    	return err
    }

    return nil
}

func updateReviewer(reviewId, status, reviewerId string) error {
	stmt, err := db.Prepare("UPDATE reviewer SET status = ? where review_id = ? AND reviewer_id = ?")
    if err != nil {
    	return err
    }

    _, err = stmt.Exec(status, reviewId, reviewerId)
    if err != nil {
    	return err
    }

    return nil
}

func countPendingReviewer(reviewId string) int {
	var count int
	err := db.QueryRow(`
			SELECT
				count(*)
			FROM
				reviewer
			WHERE
				status = 1
				AND review_id = ` + reviewId).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	return count
}

func getNotApproveCount(reviewId string) bool {
	var count int
	err := db.QueryRow(`
			SELECT
				count(*)
			FROM
				reviewer
			WHERE
				status != 2
				AND review_id = ` + reviewId).Scan(&count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return false
	}
	return true
}