package amigo

import "fmt"

var filter = map[string]func(string){
	"SuccessfulAuth":    succAuth,
	"DeviceStateChange": devStateChange,
}

func succAuth(val string) {
	fmt.Println("From f1: ", val)
}

func devStateChange(val string) {

}
