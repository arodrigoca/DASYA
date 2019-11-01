package main

import(
  "testing"
  "os")

func TestCMDline(*testing.T){

  os.Args = []string{"cmd", "-file ../../bin/lang.fx", "-debug"}
  filename, Dflag := parseArguments()
  t.Log(filename)
  t.Log(Dflag)


}
