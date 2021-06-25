package main

import (
	"fmt"
	. "mbs/parser"
	. "mbs/typechecker"
)

func main() {
	code := `a = 123;
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
	println(input);`

	block, err := ParseCode(code)

	if err != nil {
		fmt.Println("ERROR parsing the code!")
	}

	valid := TypeCheckBlock(block)

	if !valid {
		fmt.Println("ERROR typechecking the code")
	}

	block.Eval() //Code generation/execution
}
