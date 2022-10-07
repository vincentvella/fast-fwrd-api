package notification

import (
	"fmt"

	"github.com/adierkens/expo-server-sdk-go"
	"github.com/vincentvella/fast-fwrd-api/pkg/supabase"
)

type setPushData func(*expo.ExpoPushMessage) *expo.ExpoPushMessage


func SendNotification(userId string, setter setPushData) {
	
	fmt.Println("sending notification to user", userId)

  // The client can be configured to use a different API location
  client := expo.NewExpo()
	deviceIds := supabase.GetUserDevices(userId)

	if (len(deviceIds) > 0) {
		notifications := []*expo.ExpoPushMessage{}
		
		for _, v := range deviceIds {
			m := expo.NewExpoPushMessage()
			m.To = v["device_id"].(string) // Your expo push token
			var message = setter(m)

			notifications = append(notifications, message)
		}
		
		response, err := client.SendPushNotifications(notifications)
		if err != nil {
			panic(err)
		}
		fmt.Println(response)
	}
}

func SendPlannedStartNotification(message *expo.ExpoPushMessage) *expo.ExpoPushMessage {
	message.Title = "Time to Start Fasting!" // The title of the notification
	message.Body = "Be sure to log your last meal" // The body of the notification	
	return message
}

func SendPlannedEndNotification(message *expo.ExpoPushMessage) *expo.ExpoPushMessage {
	message.Title = "Time to Eat!" // The title of the notification
	message.Body = "Be sure to log when you begin eating" // The body of the notification	
	return message
}