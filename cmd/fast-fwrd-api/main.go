package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/vincentvella/fast-fwrd-api/pkg/notification"
	"github.com/vincentvella/fast-fwrd-api/pkg/supabase"
)


func getFasts(now time.Time) ([]map[string]interface{}, []map[string]interface{}) {
	// Create timestamp used to query supabase
	timestamp := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second()) 
	fmt.Println(timestamp)
	return supabase.GetFastsAt(timestamp, "start_at"), supabase.GetFastsAt(timestamp, "finish_at")
}

func notifyUserFastEnd(fasts []map[string]interface{}) {
	for _, fast := range fasts {
		notification.SendNotification(fast["uid"].(string), notification.SendPlannedEndNotification)
	}
}
func notifyUserFastStart(fasts []map[string]interface{}) {
	for _, fast := range fasts {
		notification.SendNotification(fast["uid"].(string), notification.SendPlannedStartNotification)
	}
}

func runPoll(now time.Time, ch chan struct{}) {
	fastsStarting, fastsEnding := getFasts(now)
	if len(fastsEnding) > 0 {
		notifyUserFastEnd(fastsEnding)
	}
	if len(fastsStarting) > 0 {
		notifyUserFastStart(fastsStarting)
	}
	ch <- struct{}{}
}

func PollForFasts() {
	ticker := time.NewTicker(1 * time.Second)
	// TODO: catch within poll to quit polling?
	quit := make(chan struct{})
	
	for {
			select {
			case <-quit:
				fmt.Println("Stopped")
				ticker.Stop()
				return
			case t := <-ticker.C:
				c := make(chan struct{})
				go runPoll(t, c)
			}
	}
}


func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Ok")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
			return value
	}
	return fallback
}

func main() {
	// supabase.GetFasts()
	// Start background poll
	go PollForFasts()

	// Start server
	router := httprouter.New()
	router.GET("/status", Status)
	port := getEnv("PORT", "8080")

	log.Fatal(http.ListenAndServe(":" + port, router))
}


