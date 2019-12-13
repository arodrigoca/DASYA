package fxparser

//import("fxlex")


const(
  SProg = iota
  SFunc
  SFsig
  SFdecargs
  SBody
  SStatement
  SFuncall
  SFargs
  SExpr
  SAtom
  SIter
)

type Sym struct{

  sType int

}
