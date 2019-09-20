package gaeenv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseAppYaml(t *testing.T) {
	err := applyAppYaml(true, "testdata/app.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "service-name-v1", Getenv("GAE_SERVICE"))
	assert.NotEqual(t, "", Getenv("GAE_VERSION"))
	assert.NotEqual(t, "", Getenv("GAE_INSTANCE"))
}

func Test_parseAppYaml_noServiceName(t *testing.T) {
	err := applyAppYaml(true, "testdata/app-no-service.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "default", Getenv("GAE_SERVICE"))
	assert.NotEqual(t, "", Getenv("GAE_VERSION"))
	assert.NotEqual(t, "", Getenv("GAE_INSTANCE"))
}
