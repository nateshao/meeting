package test

import (
	"meeting/internal/helper"
	"testing"
)

func TestName(t *testing.T) {
	println(helper.GetMd5("123456"))
}
