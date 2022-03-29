package richie

import (
	"log"
)

func (r *Rule) RateLimit(username string) {
	if r.AllowVisit(username) {
		log.Println(username, "访问1次,剩余:", r.RemainingVisits(username))
	} else {
		log.Println(username, "访问过多,稍后再试")
	}
}
