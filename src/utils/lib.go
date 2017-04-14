package utils

import (
	"bytes"
	"fmt"
	"os"

	"github.com/robertkrimen/otto"
)

func MakeGid(f *os.File) string {
	buff := bytes.NewBuffer(nil)
	if _, err := buff.ReadFrom(f); err != nil {
		panic(err)
	}
	runtime := otto.New()
	if _, err := runtime.Run(buff.String()); err != nil {
		panic(err)
	}

	result, err := runtime.Call("guidRandom", nil)
	if err != nil {
		panic(err)
	}
	v, err := result.ToString()
	if err != nil {
		panic(err)
	}
	return v
}

func MakePasswd(f *os.File, publicKey, passwd string) string {
	buff := bytes.NewBuffer(nil)
	if _, err := buff.ReadFrom(f); err != nil {
		panic(err)
	}
	runtime := otto.New()
	if _, err := runtime.Run(buff.String()); err != nil {
		panic(err)
	}
	pub, err := runtime.ToValue(publicKey)
	if err != nil {
		panic(err)
	}
	password, err := runtime.ToValue(passwd)
	if err != nil {
		panic(err)
	}
	result, err := runtime.Call("PasswordEncrypt", nil, password, pub)
	if err != nil {
		panic(err)
	}
	v, err := result.ToString()
	if err != nil {
		panic(err)
	}
	fmt.Println("passwd", v)
	return v
}
