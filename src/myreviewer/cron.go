package myreviewer

import(
	"log"
	"time"
	"github.com/robfig/cron"
)

func InitializeCron() {
	log.Println("cron initialized")
	c := cron.New()
	c.AddFunc("@every 30m", func() {reNotifyAllPending()})
	c.Start()
}

func reNotifyAllPending() {
	//check weekday and working hour
	now := time.Now()
	today := now.Weekday()
	hour := now.Hour()
	if today != time.Saturday && today != time.Sunday {
		if hour > 9 && hour < 18 {
			allActiveReview, err := getAllActiveReview()
			if err != nil {
				log.Println(err)
			}
			for _, review := range allActiveReview {
				for _, reviewer := range review.Reviewer {
					if reviewer.Status == 1 {
						err := reNotifyReview(review, "devel-go.tkpd:7412")
						if err != nil {
							log.Println(err)
						}
						break
					}
				}
			}
		} 
	}
}