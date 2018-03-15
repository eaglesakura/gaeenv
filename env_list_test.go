package gaeenv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEnvLine(t *testing.T) {
	// コメント行
	{
		key, value := parseEnvLine("#  test value")
		assert.Equal(t, key, "")
		assert.Equal(t, value, "")
	}
	// 前後にスペースの入ったコメント業
	{
		key, value := parseEnvLine("   #  test value   ")
		assert.Equal(t, key, "")
		assert.Equal(t, value, "")
	}
	// Keyのみ
	{
		key, value := parseEnvLine("KEY=")
		assert.Equal(t, key, "KEY")
		assert.Equal(t, value, "")
	}
	// 前後にスペースが入ってKeyのみ
	{
		key, value := parseEnvLine("  KEY=   ")
		assert.Equal(t, key, "KEY")
		assert.Equal(t, value, "")
	}
	// KEY=VALUE
	{
		key, value := parseEnvLine("KEY=VALUE")
		assert.Equal(t, key, "KEY")
		assert.Equal(t, value, "VALUE")
	}
	// 前後にスペースが入ってKEY=VALUE
	{
		key, value := parseEnvLine("   KEY=VALUE    ")
		assert.Equal(t, key, "KEY")
		assert.Equal(t, value, "VALUE")
	}
	// 前後にスペースが入ってKEY=VALUE
	{
		key, value := parseEnvLine("   KEY=VALUE    ")
		assert.Equal(t, key, "KEY")
		assert.Equal(t, value, "VALUE")
	}
}
