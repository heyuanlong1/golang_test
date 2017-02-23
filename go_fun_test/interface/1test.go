package main
import "fmt"
import "math"

type geometr interface{
	area() float64
	perim() float64
} 

type rect struct{
	width ,height float64
}
type circle struct{
	radius float64
}

func (r rect) area() float64{
	return r.width * r.height
}
func (r rect) perim() float64{
	return 2*r.width + 2*r.height
}
	
func (c circle) area() float64 {
    return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
    return 2 * math.Pi * c.radius
}

func measure(g geometr) {
    fmt.Println(g)
    fmt.Println(g.area())
    fmt.Println(g.perim())
}

func main() {
	r := rect{width:2,height:4}
	c := circle{radius:5}

	measure(r)
	measure(c)
}









