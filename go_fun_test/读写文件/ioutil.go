package main
import (
"io/ioutil"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	b , err:= ioutil.ReadFile("1input.txt")
	check(err)

	err= ioutil.WriteFile("1output.txt",b,0644)
	check(err)

	
}