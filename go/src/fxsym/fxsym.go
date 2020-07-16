package fxsym

const (
  SNone = iota
  SKey
  SStr
  SConst
  SType
  SVar
  SUnary
  SBinary
  SProc
  SFunc
  SFCall
)

type Sym struct {
  name string
  sType int
  DataType ∗Type
  //lex.Place
  FloatVal float64
  IntVal int64
  Expr ∗Expr
  Prog ∗Prog
  Body ∗Body
  Asign ∗Asign
  Iter ∗Iter
  val ∗Sym
}
