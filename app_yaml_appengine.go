//+build appengine

package gaeenv

import "log"

func applyAppengineEnvironments(parsedYaml map[string]interface{}) error {
	log.Println("Skip applyAppengineEnvironments()")
	return nil
}
