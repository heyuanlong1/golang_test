package main  
  
import "testing"  
  
func Test_getmd5(t *testing.T) {  
    for i := 0; i < 100; i++ {  
        T()  
    }  
}  
func Test_Range(t *testing.T) {  
    for i := 0; i < 100; i++ {  
        rangdom_string()  
    }  
}  
func Benchmark_getmd5(b *testing.B) {  
    for i := 0; i < b.N; i++ {  
        T()  
    }  
}  
