package main
import (
"fmt"
)

var g = 0 
func lingxing(d int,n int)() {
    if d < 0  {
    	return
    }
    tmp := n - d
    for i := 0; i < tmp; i++ {
    	fmt.Print(" ")
    }
    tmp = d * 2 - 1
    for i := 0; i < tmp; i++ {
    	fmt.Print("*")
    }
    fmt.Println("")
    
    if d >= n {
    	g = 1
    }
    if g == 0 {
    	lingxing(d + 1 ,n)
    }else{
    	lingxing(d - 1 , n)
    }
}


func main() {
	lingxing(1,10)
}