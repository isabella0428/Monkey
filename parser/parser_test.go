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

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
               literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input 			string
		operator 		string
		integerValue	int64
	} {
		{"!5;","!", 5},
		{"-15;", "-", 15},
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

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
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
		leftValue 	int64
		operator	string
		rightValue 	int64
	} {
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
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
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
                   tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
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