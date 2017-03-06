package main  
  
import (  
    "crypto/md5"  
    "math/rand"  
    "fmt"
    "encoding/hex"


)  
  
func getmd5(s string) string {  
    md := md5.New()  
    md.Write([]byte(s))  
    cipherStr := md.Sum(nil)  
    return hex.EncodeToString(cipherStr)
}
  
func T() (name ,x string) {  
    name = rangdom_string()  
    x = getmd5(name) 
    return 
}  
  
func rangdom_string() string {  
    var x []byte  
    for i := 0; i < 10; i++ {  
        a := rand.Intn(100)  
        x = append(x, byte(a+33))  
    }  
    return hex.EncodeToString(x)
}  

func main() {
    name ,x :=T()
    fmt.Println(name)
    fmt.Println(x)
}