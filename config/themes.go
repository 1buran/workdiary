package config

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

type themes struct {
	repo map[string]Theme
	tags map[string]string
}

func (t themes) List() (names []string) {
	for k, _ := range t.repo {
		names = append(names, k)
	}
	slices.Sort(names)
	return
}

func printThemePalette(output *termenv.Output, wday, doff string) {
	c1, _ := colorful.Hex(wday)
	c2, _ := colorful.Hex(doff)
	for i := range 10 {
		bg := output.Color(c1.BlendRgb(c2, float64(i)/float64(9)).Hex())
		fmt.Print(output.String("  ").Background(bg), " ")
	}
}

func (t themes) PrintList() {
	output := termenv.NewOutput(os.Stdout)
	fmt.Println("Color themes:")
	for _, k := range t.List() {
		// get theme colors and print sample pallete
		tt := t.Get(k)
		printThemePalette(output, tt.Color("workingDay"), tt.Color("dayOff"))

		// print theme name and aliases
		fmt.Print("  - ", k)
		var aliases []string
		for alias, name := range t.tags {
			if name == k {
				aliases = append(aliases, alias)
			}
		}
		if len(aliases) > 0 {
			fmt.Printf(" (aliases: %s)", strings.Join(aliases, ", "))
		}
		fmt.Println()
	}
}

func (t *themes) Add(name string, theme Theme, aliases ...string) {
	t.repo[name] = theme
	for _, alias := range aliases {
		t.tags[alias] = name
	}
}

func (t themes) Get(s string) *Theme {
	if theme, ok := t.repo[s]; ok {
		return &theme
	}

	if name, ok := t.tags[s]; ok {
		if theme, ok := t.repo[name]; ok {
			return &theme
		}
	}
	return nil
}

func NewThemes() themes {
	return themes{repo: make(map[string]Theme), tags: make(map[string]string)}
}

var Themes = NewThemes()

func init() {
	Themes.Add("Default", Theme{
		Colors: map[string]string{
			"dayOff":         "#a958ad",
			"workingDay":     "#0d420d",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "default", "def")

	Themes.Add("Dirt", Theme{
		Colors: map[string]string{
			"dayOff":         "#76767a",
			"workingDay":     "#49494d",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "dirt", "land")

	Themes.Add("StrangeThings", Theme{
		Colors: map[string]string{
			"dayOff":         "#f542ce",
			"workingDay":     "#6932cf",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "strange", "tng")

	Themes.Add("Dracula", Theme{
		Colors: map[string]string{
			"dayOff":         "#7d194e",
			"workingDay":     "#0d0d38",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "dracula", "drac", "dra")

	Themes.Add("Dune", Theme{
		Colors: map[string]string{
			"dayOff":         "#a87620",
			"workingDay":     "#1410e6",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "dune", "dun")

	Themes.Add("Minecraft", Theme{
		Colors: map[string]string{
			"dayOff":         "#cc0815",
			"workingDay":     "#09e61f",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "minecraft", "mine", "mcr")

	Themes.Add("Haki", Theme{
		Colors: map[string]string{
			"dayOff":         "#635e5b",
			"workingDay":     "#9ec29d",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "haki", "military")

	Themes.Add("Mars", Theme{
		Colors: map[string]string{
			"dayOff":         "#b34343",
			"workingDay":     "#de6b54",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "mars", "red")

	Themes.Add("Ocean", Theme{
		Colors: map[string]string{
			"dayOff":         "#0f04d6",
			"workingDay":     "#4f49c4",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "ocean", "blue")

	Themes.Add("Matrix", Theme{
		Colors: map[string]string{
			"dayOff":         "#249539",
			"workingDay":     "#243138",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "matrix", "neo")

	Themes.Add("Sand", Theme{
		Colors: map[string]string{
			"dayOff":         "#854c0c",
			"workingDay":     "#857d0c",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "sand", "sahara")

	Themes.Add("Olive", Theme{
		Colors: map[string]string{
			"dayOff":         "#b6da0e",
			"workingDay":     "#485804",
			"expectedAmount": "#ff9ff3",
			"infactAmount":   "#4cd137",
			"summary":        "#fd79a8",
		},
	}, "olive", "oil")

}
