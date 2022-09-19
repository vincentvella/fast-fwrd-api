package main

import (
	"fmt"
	"time"

	"github.com/vincentvella/fast-fwrd-api/pkg/notification"
	"github.com/vincentvella/fast-fwrd-api/pkg/supabase"
)

func getFasts(now time.Time) []map[string]interface{} {
	// Create timestamp used to query supabase
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second())
	fmt.Println(timestamp)
	return supabase.GetFastsAt(timestamp)
}

func notifyUser(fasts []map[string]interface{}) {
	for _, fast := range fasts {
		notification.SendNotification(fast["uid"].(string))
	}
}

func runPoll(now time.Time, ch chan struct{}) {
	fasts := getFasts(now)
	if len(fasts) > 0 {
		notifyUser(fasts)
	}
	ch <- struct{}{}
}

func PollForFasts() {
	wait := make(chan struct{})
	for {
			time.Sleep(1 * time.Second)
			now := time.Now()
			go runPoll(now, wait)
			<-wait
	}
}

func main() {
	// supabase.GetFasts()
	PollForFasts()
}


