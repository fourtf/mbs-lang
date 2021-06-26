package main

import (
	"fmt"
	. "mbs/parser"
	. "mbs/typechecker"
)

func main() {
	code := `
	a = 123;
	b = "abc";
	c = true;
	d = 4.2;
	if (c) {
		println("c is true");
	}
	if (a == 123) {
		println("a is 123");
	}
	
	if (c && true) {
		println("c && true");
	}
	
	if (b == "abc") {
		println("b is abc");
	}
	
	println(b + "123");
	for (;false;) {
	}
	
	for (e = 1; e < 4; e = e + 1) {
		println("e");
	}
	input = readln();
	println(input);
	for (e = 0; e < 10; e = e + 2) {
		println("*");
	}
	for (i = 0; i < 10; i = i+1) {
		a = a+++++++1;
	}
	println(input);
	`

	block, err := ParseCode(code)

	if err != nil {
		fmt.Println("ERROR parsing the code!")
	} else {

		valid := TypeCheckBlock(block)

		if !valid {
			fmt.Println("ERROR typechecking the code")
		} else {
			block.Eval() //Code generation/execution
		}
	}
}
