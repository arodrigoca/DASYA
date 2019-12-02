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

func (p *Parser) Parse() error{
  p.pushTrace("Parse")
  defer p.popTrace()
  if err := p.Prog(); err != nil{
    return err
  }

  return nil
}
