package main

import "fmt"

func main(){
	
	var x uint = 23 // explicit typing
	y := 80 // implicit typing
	z := "hii"

	fmt.Println(y)

	fmt.Println(x)
	fmt.Printf("%T ", z) // %T type
	fmt.Printf("%v ", y) // %v value
	fmt.Printf("%b ", y) // %v binary
	fmt.Printf("%b ", z) // %v binary

}