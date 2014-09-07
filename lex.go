package main

import (
	"bufio" // lectura de archivos
	"fmt"   // biblioteca IO
	"log"   // log de errores
	"os"    // administración de SO
	// "reflect"
	"strconv" // conversor
	"strings" // manejo de strings
)

type lexemes interface{}

const (
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numerals = "0123456789"
)

var tokens = map[string]lexemes{
	"KEYWORD": []string{
		"auto",
		"break",
		"case",
		"char",
		"const",
		"continue",
		"default",
		"do",
		"double",
		"else",
		"enum",
		"extern",
		"float",
		"for",
		"goto",
		"if",
		"int",
		"long",
		"register",
		"return",
		"short",
		"signed",
		"sizeof",
		"static",
		"struct",
		"switch",
		"typedef",
		"union",
		"unsigned",
		"void",
		"volatile",
		"while",
	},
	"ALPHABET":     Alphabet + "_",
	"NUMERALS":     Numerals,
	"ALPHANUMERIC": Alphabet + Numerals + "_",
	"ARTOPERATOR":  []string{"+", "-", "*", "/", "%", "+=", "-=", "*=", "/=", "%="},
	"LOGOPERATOR":  []string{"!=", ">", "<", ">=", "<=", "&&", "||", "!", "=="},
	"BITOPERATOR":  []string{"&", "|", "^", "~", ">>", "<<", "<<=", ">>=", "&=", "^=", "|="},
	"DELIMITERS":   []string{"[", "]", "(", ")", "{", "}"},
	"STRUCT":       []string{"->", "."},
	"STRING":       []string{"\u0027", "\u0022"},
	"SEPARATORS":   []string{" ", ",", ";"},
	// TODO: agregar todos los demás elementos.
}

type Lex struct {
	file  *os.File
	line  int
	colmn int
}

// Inicializador de type Lex
func NewLex() *Lex {
	self := new(Lex)

	self.line = 0

	return self
}

// Abre un archivo o manda un error si este no se encontró
func (self *Lex) OpenFile(path string) {
	var err error

	self.file, err = os.Open(path)

	if err != nil {
		fmt.Println("Error: no se encontró el archivo.")
		log.Fatal(err)
	}
}

// Cierra un archivo
func (self *Lex) CloseFile() {
	self.file.Close()
}

// Escaneando el archivo ingresado
func (self *Lex) Scanning() {
	read := bufio.NewScanner(self.file)
	combuf := new(bool) // se activa /* y se desactiva */

	for read.Scan() {
		self.line++
		sizeline := len(strings.TrimSpace(read.Text()))

		if sizeline != 0 &&
			self.elementIsNotComment(string(read.Text()), combuf) &&
			string(read.Text()[0]) != "#" {
			chunk := string(strings.TrimSpace(read.Text()))

			self.Analyze(strings.Split(chunk, ""))
		}
	}

	if err := read.Err(); err != nil {
		log.Fatal(err)
	}
}

// Indentifica si la línea es un comentario
func (self *Lex) elementIsNotComment(text string, combuf *bool) bool {
	switch {
	case text == "/*":
		*combuf = true
		return false
	case text == "*" && *combuf == true:
		return false
	case text == "*/":
		*combuf = false
		return false
	case text == "//":
		return false
	}

	return true
}

// Analiza caracter por caracter y encuentra el token correspondiente
func (self *Lex) Analyze(char []string) {
	tam := len(char)

	for i := 0; i < tam; i++ {

		switch {
		case self.lexemeInToken(char[i], strings.Split(tokens["ALPHABET"].(string), "")):
			self.isLiteral(&i, char)

		case char[i] == tokens["STRING"].([]string)[0] ||
			char[i] == tokens["STRING"].([]string)[1]:
			self.isString(&i, char)

		case self.lexemeInToken(char[i], tokens["DELIMITERS"].([]string)):
			self.isDelimiter(char[i])

		case self.lexemeInToken(char[i], tokens["SEPARATORS"].([]string)):
			self.isSeparator(char[i])

		case self.lexemeInToken(char[i], tokens["LOGOPERATOR"].([]string)):
			self.isLogOperator(&i, char)

		case self.lexemeInToken(char[i], tokens["BITOPERATOR"].([]string)):
			self.isBitOperator(&i, char)

		case char[i] == "?":
			self.printLexeme("?", "TEROPERATOR")

		case char[i] == "=":
			if char[i+1] == "=" || char[i-1] == "=" {
				self.printLexeme("==", "LOGOPERATOR")
				i += 2
			} else {
				self.printLexeme("=", "EQUAL")
			}

		default:
		}
	}
}

// Encuentra en los tokens arreglos si encuentra el caracter
func (self *Lex) lexemeInToken(char string, tokens []string) bool {
	for _, val := range tokens {
		if char == val {
			return true
		}
	}
	return false
}

// Obtiene solamente los valores entre separatores y delimitadores
func (self *Lex) noLimiters(i *int, lexeme *string, char []string) {
	for {
		*lexeme += char[*i]
		*i++

		if self.lexemeInToken(char[*i], tokens["DELIMITERS"].([]string)) ||
			self.lexemeInToken(char[*i], tokens["SEPARATORS"].([]string)) {
			break
		}
	}
	*i--
}

// Encuentra si el lexeme es un keyword o identificador
func (self *Lex) isLiteral(i *int, char []string) {
	lexeme := ""

	self.noLimiters(i, &lexeme, char)

	switch {
	case self.lexemeInToken(lexeme, tokens["KEYWORD"].([]string)):
		self.isKeyword(lexeme)
	default:
		fmt.Println(lexeme)
	}
}

// Imprime el lexeme y el token para operadores aritméticos
func (self *Lex) isArtOperator(i *int, char []string) {}

// Imprime el lexeme y el token para operadores de bits
func (self *Lex) isBitOperator(i *int, char []string) {
	lexeme := ""

	self.noLimiters(i, &lexeme, char)
	if self.lexemeInToken(lexeme, tokens["BITOPERATOR"].([]string)) {
		self.printLexeme(lexeme, "BITOPERATOR")
	} else if self.lexemeInToken(lexeme, tokens["LOGOPERATOR"].([]string)) {
		self.printLexeme(lexeme, "LOGOPERATOR")
	}
}

// Imprime el lexeme y el token para operadores lógicos
func (self *Lex) isLogOperator(i *int, char []string) {
	lexeme := ""

	self.noLimiters(i, &lexeme, char)
	if self.lexemeInToken(lexeme, tokens["LOGOPERATOR"].([]string)) {
		self.printLexeme(lexeme, "LOGOPERATOR")
	} else if self.lexemeInToken(lexeme, tokens["BITOPERATOR"].([]string)) {
		self.printLexeme(lexeme, "BITPERATOR")
	}
}

func (self *Lex) isKeyword(char string) {
	self.printLexeme(char, "KEYWORD")
}

func (self *Lex) isAlphabet(char string) {
}

// Imprime un delimitador
func (self *Lex) isDelimiter(char string) {
	self.printLexeme(char, "DELIMITER")
}

// Imprime separadores
func (self *Lex) isSeparator(char string) {
	if char != " " {
		self.printLexeme(char, "SEPARATOR")
	}
}

// Armar el string e imprimirlo
func (self *Lex) isString(i *int, char []string) {
	lexeme := ""

	for {
		lexeme += char[*i]
		*i++

		if char[*i] == tokens["STRING"].([]string)[0] ||
			char[*i] == tokens["STRING"].([]string)[1] {
			lexeme += char[*i]
			break
		}
	}
	self.printLexeme(lexeme, "STRING")
}

// Imprime un lexeme y su token
func (self *Lex) printLexeme(lexeme string, token string) {
	fmt.Printf("(l:%d) %s     %s\n", self.line, lexeme, token)
}

// Manda un error con columna y fila
func (self *Lex) error() {
	err := "Error en (l:" + strconv.Itoa(self.line) + ", c:" + strconv.Itoa(self.colmn) + ")"
	log.Fatal(err)
}

// Función principal
func main() {
	lex := NewLex()

	lex.OpenFile("c.c")
	defer lex.CloseFile()
	lex.Scanning()

	// fmt.Println(tokens["COMMENT"].([]string)[0])
}
