package transform

import (
	"github.com/speps/go-hashids"
)

type UrlShortener struct {
	hid *hashids.HashID
}

func shuffle(alphabet []string, salt []string) {
	if len(salt) == 0 {
		return
	}

	for i, v, p := len(alphabet)-1, 0, 0; i > 0; i-- {
		p += int(salt[v][0])
		j := (int(salt[v][0]) + v + p) % i
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
		v = (v + 1) % len(salt)
	}
	return
}

func NewShortener(salt string) *UrlShortener {
	hid := hashids.New()
	return &UrlShortener{hid: hid}
}

func NewShortenerWithAlphabet(salt, alphabet string) *UrlShortener {
	data := &hashids.HashIDData{Alphabet: alphabet, Salt: salt, MinLength: 3}
	hid := hashids.NewWithData(data)
	return &UrlShortener{hid: hid}
}

func (s UrlShortener) Encode(i int64) (string, error) {
	slice := make([]int64, 1)
	slice = append(slice, i)

	result, err := s.hid.EncodeInt64(slice)
	if err != nil {
		return "", err
	}

	return result, err
}
