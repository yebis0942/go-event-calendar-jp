package gojpcal

type Config struct {
	ConnpassGroups []string
}

func LoadConfig() *Config {
	return &Config{
		ConnpassGroups: []string{
			"asakusago",
			"ehimego",
			"fukuokago",
			"gocon",
			"golangtokyo",
			"go-online",
			"gopherdojo",
			"gophers-ex",
			"gospecreading",
			"gotalk",
			"kamakurago",
			"kanazawago",
			"kobego",
			"kyotogo",
			"nobishii-go",
			"okayamago",
			"sendaigo",
			"shibuya-go",
			"technical-book-reading-2",
			"tenntenn",
			"tinygo-keeb",
			"umedago",
			"womenwhogo-tokyo",
			"yokohama-go-reading",
		},
	}
}
