package assert

import "fmt"

func Equals[T comparable](expected, actual T) {
    if expected != actual {
        panic(fmt.Sprintf("values were unexpectedly not equal. expected=%v, got=%v", expected, actual))
    }
}

func NotEquals[T comparable](expected, actual T) {
    if expected == actual {
        panic(fmt.Sprintf("values were unexpectedly equal. expected=%v, got=%v", expected, actual))
    }
}

func True(value bool) {
    if !value {
        panic("value was unexpectedly false")
    }
}
