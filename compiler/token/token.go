// Package token defines lexical token types and the Token struct for moonBASIC.
package token

// TokenType classifies a scanned token.
type TokenType int

const (
	Illegal TokenType = iota
	EOF

	// Identifiers and literals
	IDENT
	INT
	FLOAT
	STRING

	NEWLINE

	// Delimiters
	LPAREN
	RPAREN
	COMMA
	COLON
	DOT

	// Variable suffix tokens (consumed by lexer as part of identifier)
	HASH     // # float suffix
	DOLLAR   // $ string suffix
	QUESTION // ? bool suffix

	// Operators (single-char and multi)
	ASSIGN // =
	PLUS
	MINUS
	STAR
	SLASH
	CARET
	EQ  // = in comparisons (same lexeme as ASSIGN; parser disambiguates)
	NEQ // <>
	LT
	GT
	LTE
	GTE
	PLUSEQ
	MINUSEQ
	STAREQ
	SLASHEQ

	// Keywords
	IF
	THEN
	ELSE
	ELSEIF
	ENDIF
	WHILE
	WEND
	ENDWHILE
	FOR
	TO
	DOWNTO
	STEP
	NEXT
	REPEAT
	UNTIL
	DO
	LOOP
	EXIT
	CONTINUE
	SELECT
	CASE
	DEFAULT
	ENDSELECT
	FUNCTION
	ENDFUNCTION
	RETURN
	TYPE
	FIELD
	ENDTYPE
	NEW
	DELETE
	EACH
	GOTO
	GOSUB
	AND
	OR
	NOT
	XOR
	MOD
	DIM
	AS
	REDIM
	PRESERVE
	LOCAL
	GLOBAL
	CONST
	STATIC
	SWAP
	ERASE
	INCLUDE
	TRUE
	FALSE
	NULL
	END   // program termination (bare END)
	PRINT // global built-in kept as keyword for fast recognition
)

// Token is one lexical unit with source position.
type Token struct {
	Type TokenType
	Lit  string
	Line int
	Col  int
}

var keywords = map[string]TokenType{
	"IF":          IF,
	"THEN":        THEN,
	"ELSE":        ELSE,
	"ELSEIF":      ELSEIF,
	"ENDIF":       ENDIF,
	"WHILE":       WHILE,
	"WEND":        WEND,
	"ENDWHILE":    ENDWHILE,
	"FOR":         FOR,
	"TO":          TO,
	"DOWNTO":      DOWNTO,
	"STEP":        STEP,
	"NEXT":        NEXT,
	"REPEAT":      REPEAT,
	"UNTIL":       UNTIL,
	"DO":          DO,
	"LOOP":        LOOP,
	"EXIT":        EXIT,
	"CONTINUE":    CONTINUE,
	"SELECT":      SELECT,
	"CASE":        CASE,
	"DEFAULT":     DEFAULT,
	"ENDSELECT":   ENDSELECT,
	"FUNCTION":    FUNCTION,
	"ENDFUNCTION": ENDFUNCTION,
	"RETURN":      RETURN,
	"TYPE":        TYPE,
	"FIELD":       FIELD,
	"ENDTYPE":     ENDTYPE,
	"NEW":         NEW,
	"DELETE":      DELETE,
	"EACH":        EACH,
	"GOTO":        GOTO,
	"GOSUB":       GOSUB,
	"AND":         AND,
	"OR":          OR,
	"NOT":         NOT,
	"XOR":         XOR,
	"MOD":         MOD,
	"DIM":         DIM,
	"AS":          AS,
	"REDIM":       REDIM,
	"PRESERVE":    PRESERVE,
	"LOCAL":       LOCAL,
	"GLOBAL":      GLOBAL,
	"CONST":       CONST,
	"STATIC":      STATIC,
	"SWAP":        SWAP,
	"ERASE":       ERASE,
	"INCLUDE":     INCLUDE,
	"TRUE":        TRUE,
	"FALSE":       FALSE,
	"NULL":        NULL,
	"END":         END,
}

// LookupKeyword returns the keyword token type for an uppercase ident, or IDENT.
func LookupKeyword(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return IDENT
}

// String returns a readable name for debugging and tests.
func (t TokenType) String() string {
	switch t {
	case Illegal:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case NEWLINE:
		return "NEWLINE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case DOT:
		return "DOT"
	case HASH:
		return "HASH"
	case DOLLAR:
		return "DOLLAR"
	case QUESTION:
		return "QUESTION"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case CARET:
		return "CARET"
	case EQ:
		return "EQ"
	case NEQ:
		return "NEQ"
	case LT:
		return "LT"
	case GT:
		return "GT"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case PLUSEQ:
		return "PLUSEQ"
	case MINUSEQ:
		return "MINUSEQ"
	case STAREQ:
		return "STAREQ"
	case SLASHEQ:
		return "SLASHEQ"
	case IF:
		return "IF"
	case THEN:
		return "THEN"
	case ELSE:
		return "ELSE"
	case ELSEIF:
		return "ELSEIF"
	case ENDIF:
		return "ENDIF"
	case WHILE:
		return "WHILE"
	case WEND:
		return "WEND"
	case ENDWHILE:
		return "ENDWHILE"
	case FOR:
		return "FOR"
	case TO:
		return "TO"
	case DOWNTO:
		return "DOWNTO"
	case STEP:
		return "STEP"
	case NEXT:
		return "NEXT"
	case REPEAT:
		return "REPEAT"
	case UNTIL:
		return "UNTIL"
	case DO:
		return "DO"
	case LOOP:
		return "LOOP"
	case EXIT:
		return "EXIT"
	case CONTINUE:
		return "CONTINUE"
	case SELECT:
		return "SELECT"
	case CASE:
		return "CASE"
	case DEFAULT:
		return "DEFAULT"
	case ENDSELECT:
		return "ENDSELECT"
	case FUNCTION:
		return "FUNCTION"
	case ENDFUNCTION:
		return "ENDFUNCTION"
	case RETURN:
		return "RETURN"
	case TYPE:
		return "TYPE"
	case FIELD:
		return "FIELD"
	case ENDTYPE:
		return "ENDTYPE"
	case NEW:
		return "NEW"
	case DELETE:
		return "DELETE"
	case EACH:
		return "EACH"
	case GOTO:
		return "GOTO"
	case GOSUB:
		return "GOSUB"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case XOR:
		return "XOR"
	case MOD:
		return "MOD"
	case AS:
		return "AS"
	case DIM:
		return "DIM"
	case REDIM:
		return "REDIM"
	case PRESERVE:
		return "PRESERVE"
	case LOCAL:
		return "LOCAL"
	case GLOBAL:
		return "GLOBAL"
	case CONST:
		return "CONST"
	case STATIC:
		return "STATIC"
	case SWAP:
		return "SWAP"
	case ERASE:
		return "ERASE"
	case INCLUDE:
		return "INCLUDE"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case END:
		return "END"
	default:
		return "UNKNOWN"
	}
}
