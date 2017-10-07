package myreviewer

import (
	"log"
	"math/rand"
	"time"
)

func randomMember(team *Team, developerId string) []*Member {
	members, err := getAvailableSEMemberList(team.TeamId, developerId)
	if err != nil {
		log.Println(err)
	}
	reviewers := []*Member{}
	length := len(members)

	if length <= 2 {
		for idx, _ := range members {
			reviewers = append(reviewers, members[idx])
		}
	} else {
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		var int1, int2 int
		var reviewer1, reviewer2 string
		int1 = random.Intn(length)
		reviewer1 = members[int1].Username
		reviewers = append(reviewers, members[int1])

		loopCount := 0
		for {
			int2 = random.Intn(length)
			reviewer2 = members[int2].Username
			if reviewer1 != reviewer2 {
				reviewers = append(reviewers, members[int2])
				break
			} else if loopCount > 10 {
				break
			}
			loopCount++
		}
	}
	return reviewers
}
