package main

import ("testing"
        "strings"
        "bufio"
)


func TestQuijote(t *testing.T){

  scanner := bufio.NewScanner(strings.NewReader("Hola. Saludos. fsadf;lkjsbdl;3214: qu . ggfßðfg®äåé ðf f : ßðfß"))
  var testDictionary map[string]wordInfo
	testDictionary = make(map[string]wordInfo)
  scanAndProcess(scanner, testDictionary)
  if len(testDictionary) < 4{
    t.Errorf("program failed, returned %v", testDictionary)
  }

}

func TestQuijoteVoid(t *testing.T){

  scanner := bufio.NewScanner(strings.NewReader(""))
  var testDictionary map[string]wordInfo
	testDictionary = make(map[string]wordInfo)
  scanAndProcess(scanner, testDictionary)
  if len(testDictionary) > 0{
    t.Errorf("program failed, returned %v", testDictionary)
  }

}
