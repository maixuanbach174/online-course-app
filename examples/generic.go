package main

import (
	"fmt"
)

type Person[T any] interface {
	Name(t T) string
}

type Teacher struct {
	name string
}

func (teacher Teacher) Name(t string) string {
	return teacher.name + " " + t
}

type Student struct {
	name string
}

func (s Student) Name(t string) string {
	return s.name
}

func main() {
	var p Person[string] = &Teacher{name: "John"}
	fmt.Println(p.Name("Mr."))
}
