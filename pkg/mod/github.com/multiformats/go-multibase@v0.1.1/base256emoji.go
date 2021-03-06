package multibase

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

var base256emojiTable = [256]rune{
	// Curated list, this is just a list of things that *somwhat* are related to our comunity
	'๐', '๐ช', 'โ', '๐ฐ', '๐', // Space
	'๐', '๐', '๐', '๐', '๐', '๐', '๐', '๐', // Moon
	'๐', '๐', '๐', // Our Home, for now (earth)
	'๐',                // Dragon!!!
	'โ',                // Our Garden, for now (sol)
	'๐ป', '๐ฅ', '๐พ', '๐ฟ', // Computer
	// The rest is completed from https://home.unicode.org/emoji/emoji-frequency/ at the time of creation (december 2021) (the data is from 2019), most used first until we reach 256.
	// We exclude modifier based emojies (such as flags) as they are bigger than one single codepoint.
	// Some other emojies were removed adhoc for various reasons.
	'๐', 'โค', '๐', '๐คฃ', '๐', '๐', '๐', '๐ญ', '๐', '๐',
	'๐', '๐', '๐', '๐ฅ', '๐ฅฐ', '๐', '๐', '๐', '๐ข', '๐ค',
	'๐', '๐', '๐ช', '๐', 'โบ', '๐', '๐ค', '๐', '๐', '๐',
	'๐', '๐น', '๐คฆ', '๐', '๐', 'โ', 'โจ', '๐คท', '๐ฑ', '๐',
	'๐ธ', '๐', '๐', '๐', '๐', '๐', '๐', '๐', '๐', '๐คฉ',
	'๐', '๐', '๐ค', '๐', '๐ฏ', '๐', '๐', '๐ถ', '๐', '๐คญ',
	'โฃ', '๐', '๐', '๐', '๐ช', '๐', '๐ฅ', '๐', '๐', '๐ฉ',
	'๐ก', '๐คช', '๐', '๐ฅณ', '๐ฅ', '๐คค', '๐', '๐', '๐ณ', 'โ',
	'๐', '๐', '๐ด', '๐', '๐ฌ', '๐', '๐', '๐ท', '๐ป', '๐',
	'โญ', 'โ', '๐ฅบ', '๐', '๐', '๐ค', '๐ฆ', 'โ', '๐ฃ', '๐',
	'๐', 'โน', '๐', '๐', '๐ ', 'โ', '๐', '๐บ', '๐', '๐ป',
	'๐', '๐', '๐', '๐', '๐น', '๐ฃ', '๐ซ', '๐', '๐', '๐ต',
	'๐ค', '๐', '๐ด', '๐ค', '๐ผ', '๐ซ', 'โฝ', '๐ค', 'โ', '๐',
	'๐คซ', '๐', '๐ฎ', '๐', '๐ป', '๐', '๐ถ', '๐', '๐ฒ', '๐ฟ',
	'๐งก', '๐', 'โก', '๐', '๐', 'โ', 'โ', '๐', '๐ฐ', '๐คจ',
	'๐ถ', '๐ค', '๐ถ', '๐ฐ', '๐', '๐ข', '๐ค', '๐', '๐จ', '๐จ',
	'๐คฌ', 'โ', '๐', '๐บ', '๐ค', '๐', '๐', '๐ฑ', '๐', '๐ถ',
	'๐ฅด', 'โถ', 'โก', 'โ', '๐', '๐ธ', 'โฌ', '๐จ', '๐', '๐ฆ',
	'๐ท', '๐บ', 'โ ', '๐', '๐', '๐ต', '๐', '๐คฒ', '๐ค ', '๐คง',
	'๐', '๐ต', '๐', '๐ง', '๐พ', '๐', '๐', '๐ค', '๐', '๐คฏ',
	'๐ท', 'โ', '๐ง', '๐ฏ', '๐', '๐', '๐ค', '๐', '๐', 'โ',
	'๐ด', '๐ฃ', '๐ธ', '๐', '๐', '๐ฅ', '๐คข', '๐', '๐ก', '๐ฉ',
	'๐', '๐ธ', '๐ป', '๐ค', '๐คฎ', '๐ผ', '๐ฅต', '๐ฉ', '๐', '๐',
	'๐ผ', '๐', '๐ฃ', '๐ฅ',
}

var base256emojiReverseTable map[rune]byte

func init() {
	base256emojiReverseTable = make(map[rune]byte, len(base256emojiTable))
	for i, v := range base256emojiTable {
		base256emojiReverseTable[v] = byte(i)
	}
}

func base256emojiEncode(in []byte) string {
	var l int
	for _, v := range in {
		l += utf8.RuneLen(base256emojiTable[v])
	}
	var out strings.Builder
	out.Grow(l)
	for _, v := range in {
		out.WriteRune(base256emojiTable[v])
	}
	return out.String()
}

type base256emojiCorruptInputError struct {
	index int
	char  rune
}

func (e base256emojiCorruptInputError) Error() string {
	return "illegal base256emoji data at input byte " + strconv.FormatInt(int64(e.index), 10) + ", char: '" + string(e.char) + "'"
}

func (e base256emojiCorruptInputError) String() string {
	return e.Error()
}

func base256emojiDecode(in string) ([]byte, error) {
	out := make([]byte, utf8.RuneCountInString(in))
	var stri int
	for i := 0; len(in) > 0; i++ {
		r, n := utf8.DecodeRuneInString(in)
		in = in[n:]
		var ok bool
		out[i], ok = base256emojiReverseTable[r]
		if !ok {
			return nil, base256emojiCorruptInputError{stri, r}
		}
		stri += n
	}
	return out, nil
}
