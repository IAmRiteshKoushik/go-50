package main

import "testing"

func TestConverter(t *testing.T) {
    result := Converter(3);
    if result != "Converted" {
        t.Errorf("Result was incorrect, got: %s, want :%s", result, "Converted")
    }
}
