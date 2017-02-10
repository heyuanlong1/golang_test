package main
import (
"os"
"io"
"fmt"
"bufio"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	//open input file
	fi,err := os.Open("1input.txt")
	check(err)
	defer func() {
		err:=fi.Close()
		check(err)
	}()
	scanner := bufio.NewScanner(fi)

	for scanner.Scan(){
		fmt.Print(scanner.Text())
	}
	fmt.Println()

	fi2,err := os.Open("1input.txt")
	check(err)
	defer func() {
		err:=fi2.Close()
		check(err)
	}()
	r := bufio.NewReader(fi2)
	for{
		//buf,err := r.ReadString('\n')
		buf,err := r.ReadBytes('\n')
		if err != nil && err != io.EOF{
			panic(err)
		}

		fmt.Print(buf)
		if err == io.EOF{
			break
		}		
	}
}