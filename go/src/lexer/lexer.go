package main

import( "fmt"
        "os"
        "flag"
        "bufio"
)

type RuneScanner interface{
    ReadRune() (r rune, size int, err error)
    UnreadRune() error
}

type Lexer struct{
    file string
    line int
    rs RuneScanner
    lastrune rune
}


func parseArguments() string{

	filenamePtr := flag.String("file", "", "filename to read")
	flag.Parse()
	if *filenamePtr == ""{
			fmt.Println("Error: at least argument -file is necessary.")
			os.Exit(1)
	}

	return *filenamePtr
}

func NewLexer(rs RuneScanner, filename string) (l *Lexer){
    l = &Lexer{line: 1}
    l.file = filename
    l.rs = rs
    return l
}

func (l *Lexer) get() (r rune){

    var err error

    r, _, err = l.rs.ReadRune()

    if err == nil{
        l.lastrune = r
        if r == '\n'{
            l.line++
        }
    }

    /*
    if err == io.EOF{
        l.lastrune ==
    }
    */

    return r
}


func main(){

    filename := parseArguments()
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: No such file or directory.")
		return
	}
	defer file.Close()

    reader := bufio.NewReader(file)

}
