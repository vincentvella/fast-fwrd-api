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

func notifyUser(fasts []map[string]interface{}) {
	for _, fast := range fasts {
		notification.SendNotification(fast["uid"].(string))
	}
}

func runPoll(now time.Time, ch chan struct{}) {
	fastsStarting, fastsEnding := getFasts(now)
	if len(fastsEnding) > 0 {
		notifyUser(fastsEnding)
	}
	if len(fastsStarting) > 0 {
		notifyUser(fastsStarting)
	}
	ch <- struct{}{}
}

func PollForFasts() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			go runPoll(now, quit)
		case <-quit:
			ticker.Stop()
			return	
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


