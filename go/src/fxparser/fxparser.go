package fxparser

import(
        "fxlex"
        "strings"
        "fmt"
        "os"
        "errors"
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

func (p *Parser) Finside() error{
//<FINSIDE> :: = <FDECARGS> ')' |')'

  p.pushTrace("FINSIDE")
  _, _ = p.l.Lex()
  return nil

}

func (p *Parser) Fsig() error{
  //<FSIG> :: = 'func' ID '(' <FINSIDE>
  p.pushTrace("FSIG")
  _, err, isFunc := p.match(fxlex.TokFunc)

  if err != nil || !isFunc{
    err = errors.New("Missing 'func' token on function definition")
    return err
  }

  _, err, isId := p.match(fxlex.TokId)

  if err != nil || !isId{
    err = errors.New("Missing function id on function definition")
    return err
  }

  _, err, isLpar := p.match(fxlex.TokType('('))

  if err != nil || !isLpar{
    err = errors.New("Missing '(' token on function definition")
    return err
  }

  if err := p.Finside(); err != nil{
    return err
  }

  return nil
}

func (p *Parser) Body() error{
  p.pushTrace("BODY")
  _, _ = p.l.Lex()
  return nil

}

func (p *Parser) Func() error{
  //<FUNC> ::= <FSIG> '{' <BODY> '}'
  p.pushTrace("FUNC")
  defer p.popTrace()

  if err := p.Fsig(); err != nil{
    return err
  }

  _, err, isLbra := p.match(fxlex.TokType('{'))

  if err != nil || !isLbra{
    err = errors.New("Missing '{' token on function")
    return err
  }

  if err := p.Body(); err != nil{
    return err
  }

  _, err, isRbra := p.match(fxlex.TokType('}'))

  if err != nil || !isRbra{
    err = errors.New("Missing '}' token on function")
    return err
  }
  fmt.Println(isRbra)

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
