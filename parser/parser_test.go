package parser

import (
	"fmt"
	"testing"
	"../ast"
	"../lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserError(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil!\n")	
	}
	
	if len(program.Statements) != 3 {
		t.Fatalf("ParseProgram does not contain 3 Statements, got: %d\n", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	} {
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testReturnStatements(t *testing.T) {
	input := `
   		return 5;
   		return 10;
   		return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",len(program.Statements))
	}

	for _, tt := range program.Statements {
		returnStmt, ok := tt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", returnStmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
    			returnStmt.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, expectedIdentifier string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
	}

	// Check if stmt stores *LetStatement
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
        t.Errorf("stmt not *ast.LetStatement. got=%T", stmt)
		return false
	}

	if letStmt.Name.Value != expectedIdentifier {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", 
			expectedIdentifier, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != expectedIdentifier {
		t.Errorf("s.Name not '%s'. got=%s", expectedIdentifier, letStmt.Name)
		return false
	}
	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "footbar;"

	l := lexer.New(input)
	p := New(l)
	program := p.parseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
               program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if ident.Value != "footbar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "footbar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
               ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)

	program := p.parseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
               len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
               program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	testLiteralExpression(t, literal, 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input 			string
		operator 		string
		value			interface{}
	} {
		{"!5;","!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.parseProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                   1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}	

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
               integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input		string
		leftValue 	interface{}
		operator	string
		rightValue 	interface{}
	} {
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.parseProgram()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                   1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                   program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		testInfixExpression(t, exp, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input 		string
		expected 	string
	} {
		{
			"-a * b",
			"((-a) * b)",
		}, 
		{
			"!-a",
        	"(!(-a))",
    	},
		{
			"a + b + c",
        	"((a + b) + c)",
    	},
		{
			"a + b - c",
        	"((a + b) - c)",
    	},
		{
			"a * b * c",
        	"((a * b) * c)",
    	},
		{
			"a * b / c",
        	"((a * b) / c)",
    	},
		{
			"a + b / c",
        	"(a + (b / c))",
    	},
    	{
        	"a + b * c + d / e - f",
        	"(((a + (b * c)) + (d / e)) - f)",
		}, 
		{
			"3 + 4; -5 * 5",
        	"(3 + 4)((-5) * 5)",
    	},
    	{
        	"5 > 4 == 3 < 4",
        	"((5 > 4) == (3 < 4))",
		}, 
		{
			"5 < 4 != 3 > 4",
        	"((5 < 4) != (3 > 4))",
		},
    	{
        	"3 + 4 * 5 == 3 * 1 + 4 * 5",
        	"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		}, 
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
        	"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true", },
		{
			"false",
			"false", },
		{
   			"3 > 5 == false",
    		"((3 > 5) == false)",
		}, 
		{
			"3 < 5 == true",
            "((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
           },
		{
			"(5 + 5) * 2",
            "((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
            "(2 / (5 + 5))",
        },
		{
			"-(5 + 5)",
            "(-(5 + 5))",
        },
        {
            "!(true == true)",
            "(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.parseProgram()
		checkParserError(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool{
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
               ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
		case int:
			return testIntegerLiteral(t, exp, int64(v))
		case int64:
			return testIntegerLiteral(t, exp, int64(v))
		case string:
			return testIdentifier(t, exp, v)
		case bool:
			return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, 
		operator string, right interface{}) bool {
			inExp, ok := exp.(*ast.InfixExpression)
			if !ok {
				t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
				return false
			}

			if !testLiteralExpression(t, inExp.Left, left) {
				return false
			}

			if inExp.Operator != operator {
				t.Errorf("exp.Operator is not '%s'. got=%q", operator, inExp.Operator)
				return false
			}

			if !testLiteralExpression(t, inExp.Right, right) {
				return false 
			}
			return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct{
		input 		string
		expected    bool
	} {{"true;",true},}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.parseProgram()
		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                   1, len(program.Statements))	
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                   program.Statements[0])
		}

		testLiteralExpression(t, stmt.Expression, true)
	}
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
        t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false 
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) { 
		t.Errorf("bo.TokenLiteral not %t. got=%s",
               value, bo.TokenLiteral())
		return false 
	}
	return true
}

// func TestIfExpression(t *testing.T) {
// 	input := `if (x < y) { x }`

// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.parseProgram()
// 	checkParserError(t, p)

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
//                1, len(program.Statements))
// 	}

// 	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
//                program.Statements[0])
// 	}

// 	exp, ok := stmt.Expression.(*ast.IfExpression)
// 	if !ok {
// 		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
//                stmt.Expression)
// 	}

// 	if !testInfixExpression(t, exp.Condition, "x", "<" , "y" ) {
// 		return
// 	}

// 	if len(exp.Consequence.Statements) != 1 {
// 		t.Errorf("consequence is not 1 statements. got=%d\n",
//                len(exp.Consequence.Statements))
// 	}

// 	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
//               exp.Consequence.Statements[0])
// 	}

// 	if !testIdentifier(t, consequence.Expression, "x") {
// 		return
// 	}

// 	if exp.Alternative != nil {
// 		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
// 	}
// }

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.parseProgram()
	checkParserError(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] is not *ast.ReturnStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statements[0] is not *ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("Alternative is not 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Statements[0] is not *ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func checkParserError(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}