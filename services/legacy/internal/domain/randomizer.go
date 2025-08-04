package domain

import (
	"math/rand"
	"strings"
)

type Randomizer struct {
	// Гласные
	a []string
	// Согласные
	b []string
	// Не используются
	c []string
}

func NewRandomizer() *Randomizer {
	return &Randomizer{
		a: []string{"а", "е", "ё", "и", "о", "у", "ы", "э", "ю", "я"},
		b: []string{"б", "в", "г", "д", "ж", "з", "й", "к", "л", "м", "н", "п", "р", "с", "т", "ф", "х", "ц", "ш", "щ", "ч"},
		c: []string{"ъ", "ь"},
	}
}

func (r *Randomizer) syllable() string {
	// Слог состоит из согласной и гласной
	syllable := r.stringChoice(r.b) + r.stringChoice(r.a)

	// С некоторым шансом гласная будет двойной
	if rand.Intn(100) < 35 {
		syllable += r.stringChoice(r.a)
	}

	return syllable
}

func (r *Randomizer) word(minLen, maxLen int, firstUpper bool) string {
	word := ""

	syllableCount := minLen + rand.Intn(maxLen)

	for range syllableCount {
		word += r.syllable()
	}

	if firstUpper {
		tmp := []rune(word)
		word = strings.ToTitle(string(tmp[0])) + string(tmp[1:])
	}

	return word
}

func (g *Randomizer) Name() string {
	return g.word(2, 4, true)
}

func (g *Randomizer) stringChoice(src []string) string {
	if len(src) == 0 {
		return ""
	}

	return src[rand.Intn(len(src))]
}
