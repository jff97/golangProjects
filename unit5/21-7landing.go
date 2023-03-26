package main
import (
	"fmt"
	"encoding/json"
)

func main() {
	type location struct {
		name string 'json:"name"'
		lat  float64 'json:latitude'
		long float64 'json:longitude'
  	}
  
  	locations := []location{
		{name: "Bradbury Landing", lat: -4.5895, long: 137.4417},
		{name: "Columbia Memorial Station", lat: -14.5684, long: 175.472636},
		{name: "Challenger Memorial Station", lat: -1.9462, long: 354.4734},
  	}

	bytes := json.Marshal(locations)
	fmt.Println(string(bytes))
}