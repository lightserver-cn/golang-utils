package string

import (
	"crypto/rand"
	"io"
	"regexp"

	"github.com/sirupsen/logrus"
)

var characterString = "1234567890" + "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var numberString = "1234567890"
var numberSequence = regexp.MustCompile(`([a-zA-Z](\d+)([a-zA-Z]))`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

func RandNumString(length int) string {
	return EncodeToString(numberString, length)
}

func RandString(length int) string {
	return EncodeToString(characterString, length)
}

func EncodeToString(str string, length int) string {
	table := []byte(str)
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		logrus.Error(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
