package fxparser

import(
        "fxlex"
        "strings"
        "fmt"
        "os"
)

type Parser struct{
  l *fxlex.Lexer
  depth int
  DebugDesc bool
}

func NewParser(l *fxlex.Lexer) *Parser{
  return &Parser{l, 0, true}
}

func (p *Parser) pushTrace(tag string){

  if p.DebugDesc{
    tabs := strings.Repeat("\t", p.depth)
    fmt.Fprintf(os.Stderr, "%s%s\n", tabs, tag)
  }
  p.depth++
}

func (p *Parser) popTrace(){
  p.depth--
}

func (p *Parser) match(tT fxlex.TokType) (t fxlex.Token, e error, isMatch bool){

  t, err := p.l.Peek()
  //fmt.Print("Peek says: ")
  //t.PrintToken()
  if err != nil{
    return fxlex.Token{}, err, false
  }
  if t.Type != tT{
    return t, nil, false
  }
  t, err = p.l.Lex()
  //fmt.Print("LEXED: ")
  return t, nil, true

}

func (p *Parser) Func() error{
  //<FSIG> :: = 'func' ID '(' <FINSIDE>
  
}

func (p *Parser) Func() error{
  //<FUNC> ::= <FSIG> '{' <BODY> '}'
  p.pushTrace("FUNC")
  _, _ = p.l.Lex()
  return nil

}

func (p *Parser) End() error{
  //<END> ::= <PROG> | <EOF>
  p.pushTrace("END")

  _, err, isEOF := p.match(fxlex.TokEof)

  if err != nil{
    return err
  }

  if isEOF{
    p.pushTrace("EOF")
    return nil
  }

  return p.Prog()
}

func (p *Parser) Prog() error{
  //<PROG> ::= <FUNC> <END> | <EOF>
  p.pushTrace("PROG")

  _, err, isEOF := p.match(fxlex.TokEof)

  if err != nil{
    return err
  }

  if isEOF{
    p.pushTrace("EOF")
    return nil
  }

  if err := p.Func(); err != nil{
    return err
  }else{
    return p.End()
  }

}

func (p *Parser) Parse() error{
  p.pushTrace("Parse")
  defer p.popTrace()
  if err := p.Prog(); err != nil{
    return err
  }

  return nil
}
