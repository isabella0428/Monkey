# Writing An Interpreter In Go🐒
Monkey is a programming language that I wrote based on book  **Writing An Interpreter In Go**

Thx for him for providing the amazing book！！



#### Overview

Monkey supports six different primitive types, **Integer, Boolean, Array, String, Functions and Map**.

It can do mathematical calculations, variable bindings, functions and the application of those functions, conditionals, return statements and even advanced concepts like higher-order functions and closures.

I also add five builtin functions.

```
1. len()
Input Type: Supports String and array
Return Type: int
Usage: Return the length of the String or array
eg: let a = "123";
		len(3)	//  Print 3

2. first
Input Type:	Array
Return Type: All Types
Usage:	Return the first element of the array
eg:	let a = [1, 2, 3]
		first(a)	// Print 1
		
3. last
Input Type:	Array
Return Type: All Types
Usage:	Return the last element of the array
eg:	let a = [1, 2, 3]
		first(a)	// Print 3
		
4. rest
Input Type: Array
Return type: All Types
Usage: Return the array with its first element removed
eg:	let a = [1, 2, 3]
		rest(a)  // Print [2, 3], Notify that it will not modify the original array
		a				 // Print [1, 2, 3]
    
5. puts
Input Type: All Type
Return type: None
Usage: Print All elements in the array
eg:	let a = [1, 2, 3]
		puts(a)  // Print [1, 2, 3]
```



#### Error handling

When the input type is wrong, the repl will print clear error Message.

```
Hello isabella! This is the Monkey Programming Language!
Feel free to type in commands
>>let a = 1;
>>a[1];
ERROR: index operator not supported: INTEGER
```



When it comes parser Error, a Monkey face will be printed.

```
Hello isabella! This is the Monkey Programming Language!
Feel free to type in commands
>>a[123
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
Woops! We ran into some monkey business here!
 parser errors:
        expected next token to be ], got EOF instead
```



#### Environment

Golang 		(Tested on: go1.13.4 darwin/amd64)



#### Running the console

```
go run main.go
```



#### Tests

```
1. Test Lexer
go test ./lexer

2. Test Parser
go test ./parser
go test ./ast				// Test AST Tree building

3. Test Evaluator
go test ./evaluator
go test ./object		// Test object creation in evaluator proces
```



