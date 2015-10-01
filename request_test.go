package wiki_test

import (
	"testing"

	"github.com/walle/wiki"
)

func Test_NewRequestLanguageInterpolation(t *testing.T) {
	r, err := wiki.NewRequest("example.com", "Test", "sv")
	if err != nil {
		t.Errorf("Could not create request: %s\n", err)
	}

	expected := "example.com?action=query&converttitles=&exintro=&explaintext=" +
		"&format=json&inprop=url&prop=extracts%7Cinfo&redirects=&titles=Test"
	if r.String() != expected {
		t.Errorf("Expected %s got %s", expected, r.String())
	}

	r, err = wiki.NewRequest("%s.example.com", "Test", "sv")
	if err != nil {
		t.Errorf("Could not create request: %s\n", err)
	}

	expected = "sv.example.com?action=query&converttitles=&exintro=" +
		"&explaintext=&format=json&inprop=url&prop=extracts%7Cinfo" +
		"&redirects=&titles=Test"
	if r.String() != expected {
		t.Errorf("Expected %s got %s", expected, r.String())
	}
}
