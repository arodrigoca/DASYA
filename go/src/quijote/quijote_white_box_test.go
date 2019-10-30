package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestCleanWord(t *testing.T) {

	result := CleanWord("``. .«»¡²³¤€¼½:hola")
	if result != "hola" {
		t.Logf("CleanWord failed, returned %s", result)
	} else {
		t.Logf("CleanWord suceeded, returned %s", result)
	}

}

func TestQuijote(t *testing.T) {

	scanner := bufio.NewScanner(strings.NewReader(".Hola. Saludos. .fsadf;lkjsbdl;3214: qu . :ggfßðfg®äåé: ðf f : ßðfß"))
	var testDictionary map[string]wordInfo
	testDictionary = make(map[string]wordInfo)
	ScanAndProcess(scanner, testDictionary)
	if len(testDictionary) < 4 {
		t.Errorf("program failed, returned %v", testDictionary)
	}

}
