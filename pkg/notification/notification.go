package notification

import (
	"fmt"

	"github.com/adierkens/expo-server-sdk-go"
	"github.com/vincentvella/fast-fwrd-api/pkg/supabase"
)

func SendNotification(userId string) {
	
	fmt.Println("sending notification to user", userId)

  // The client can be configured to use a different API location
  client := expo.NewExpo()
	deviceIds := supabase.GetUserDevices(userId)

	if (len(deviceIds) > 0) {
		notifications := []*expo.ExpoPushMessage{}
		
		for _, v := range deviceIds {
			m := expo.NewExpoPushMessage()
			
			m.To = v["device_id"].(string) // Your expo push token
			m.Title = "It's Fasting Time!" // The title of the notification
			m.Body = "Testing a push" // The body of the notification
			
			notifications = append(notifications, m)
		}
		
		response, err := client.SendPushNotifications(notifications)
		if err != nil {
			panic(err)
		}
		fmt.Println(response)
	}
}