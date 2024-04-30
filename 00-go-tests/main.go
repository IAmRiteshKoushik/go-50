package main

import "strconv"

func Converter(input int) string {

    hasConverted := (input % 3) == 0

    if hasConverted{
        return "Converted"
    }
    return strconv.Itoa(input)
}
