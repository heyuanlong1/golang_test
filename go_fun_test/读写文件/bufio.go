package main
import (
"io"
"os"
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
	r := bufio.NewReader(fi)

	//open output file
	fo,err := os.Create("1output.txt")
	check(err)
	defer func() {
		err:=fo.Close()
		check(err)
	}()
	w := bufio.NewWriter(fo)

	buf:=make([]byte,1024)
	for{
		n,err := r.Read(buf)
		if err != nil && err != io.EOF{
			panic(err)
		}
		if n == 0{
			break
		}

		// write a chunk
        if _, err := w.Write(buf[:n]); err != nil {
            panic(err)
        }
	}
	err = w.Flush()
	check(err)

}