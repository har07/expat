package debug

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(o interface{}) string {
	str, err := json.Marshal(o)
	if err != nil {
		fmt.Printf("PrintJSON error: %s", err.Error())
	}
	return string(str)
}