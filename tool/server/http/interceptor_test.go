package http

import "testing"

func TestCamelToSnake(t *testing.T) {
	t.Log(CamelToSnake("HelloWorld"))
	t.Log(CamelToSnake("helloWorld"))
	t.Log(CamelToSnake("Hello_world"))
	t.Log(CamelToSnake("hello_world"))

	t.Log(CamelToSnake("_HelloWorld"))
	t.Log(CamelToSnake("_helloWorld"))
	t.Log(CamelToSnake("_Hello_world"))
	t.Log(CamelToSnake("_hello_world"))
}
