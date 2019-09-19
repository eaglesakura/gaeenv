//+build !appengine

package gaeenv

import (
	"fmt"
	"os"
)

func applyAppengineEnvironments(parsedYaml map[string]interface{}) error {

	service, ok := parsedYaml["service"]
	if !ok {
		service = "default"
	}

	fmt.Printf("Apply: GAE_SERVICE=%v\n", service)

	envCache["GAE_VERSION"] = "__GAE_VERSION__"
	envCache["GAE_INSTANCE"] = "__GAE_INSTANCE__"
	envCache["GAE_SERVICE"] = fmt.Sprintf("%v", service)

	_ = os.Setenv("GAE_VERSION", "__GAE_VERSION__")
	_ = os.Setenv("GAE_INSTANCE", "__GAE_INSTANCE__")
	_ = os.Setenv("GAE_SERVICE", fmt.Sprintf("%v", service))
	return nil
}
