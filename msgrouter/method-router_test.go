package msgrouter

import (
	"testing"
)

func TestForEachMethodsMap(t *testing.T) {
	for k, v := range MethodsMap {
		t.Log(k.String())
		for k1, v1 := range v {
			t.Log("\r", k1.String(), v1)
		}
	}
}

func TestForEachMethodNames(t *testing.T) {
	for k, v := range MethodNameList {
		t.Log(k, v)
	}
}

func TestStringCMP(t *testing.T) {
	t.Log(94 & 8)

}
