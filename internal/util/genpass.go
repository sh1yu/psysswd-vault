package util

import (
    "bytes"
    "errors"
    "math/rand"
    "time"
)

var baseMap map[string]string

func init() {
    baseMap = make(map[string]string)
    baseMap["base1"] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    baseMap["base2"] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"
    baseMap["base3"] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_:@{}*()"
}

func GenPass(baseStr string, length int) (string, error) {

    var base string
    switch baseStr {
    case "base1", "base2", "base3":
        base = baseMap[baseStr]
    default:
        return "", errors.New("invalid base: " + baseStr)
    }

    rand.Seed(time.Now().UnixNano())
    var buf bytes.Buffer
    for i := 0; i < length; i++ {
        number := rand.Intn(len(base))
        buf.WriteByte(base[number])
    }



    return buf.String(), nil
}
