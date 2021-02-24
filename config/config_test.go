package config

import (
	"os"
	"testing"
)

func TestAAA(t *testing.T) {
	cnf := new(Config)
	t.Log(cnf.DB)
}
func TestDecodeConfig(t *testing.T) {
	path := "./config.toml"
	cnf, err := DecodeConfig(path)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Log(cnf)
}

func TestSafeWriteConfig(t *testing.T) {
	path := "./config.toml"
	cnf, err := DecodeConfig(path)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	p2 := "./safeConfig.toml"
	c, err := os.Create(p2)
	if err!=nil{
		t.Fatal(err)
	}
	//t.Log(cnf.API, cnf.DB)
	barr, err := ConfigComment(cnf)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Write(barr)
	if err != nil {
		t.Fatal(err)
	}
}
