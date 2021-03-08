package internal

import "math/rand"

// Some interesting names
var names = []string{
	"cygnet",
	"gosling",
	"eaglet",
	"duckling",
	"elver",
	"leveret",
	"gordon",
	"gopper",
	"gophie",
	"biker",
	"echidna",
	"quokka",
	"quola",
	"wallaby",
}

func randomName() string {
	return names[rand.Int()%len(names)]
}
