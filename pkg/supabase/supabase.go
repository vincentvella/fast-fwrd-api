package supabase

import (
	"fmt"
	"os"

	supabase "github.com/nedpals/supabase-go"
)

var supabaseClient = supabase.CreateClient(
	os.Getenv("SUPABASE_URL"),
	os.Getenv("SUPABASE_KEY"),
)

func QueryTable(tableName string) map[string]interface{} {
  var results map[string]interface{}
  err := supabaseClient.DB.From(tableName).Select("*").Single().Execute(&results)
  if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
  }
	return results
}

func GetUserDevices(userId string) [](map[string]interface{}) {
	var results [](map[string]interface{})
	err := supabaseClient.DB.From("devices").Select("*").Eq("uid", userId).Execute(&results)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
	return results
}

func GetFasts() [](map[string]interface{}) {
	var results [](map[string]interface{})
	err := supabaseClient.DB.From("fasts").Select("*").Execute(&results)
	fmt.Println(results)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
	// TODO possibly paginate
	return results
}

func GetFastsAt(timestamp string) [](map[string]interface{}) {
	var results [](map[string]interface{})
	err := supabaseClient.DB.From("fasts").Select("*").Eq("ends_at", timestamp).Execute(&results)
	fmt.Println(results)
	if err != nil {
		fmt.Println(fmt.Errorf(err.Error()))
	}
	// TODO possibly paginate
	return results
}