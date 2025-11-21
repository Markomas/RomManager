package config

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"gopkg.in/yaml.v3"
)

type Color sdl.Color

func (c *Color) UnmarshalYAML(value *yaml.Node) error {
	var hexColor string
	if err := value.Decode(&hexColor); err != nil {
		return err
	}

	if len(hexColor) > 0 && hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}

	var r, g, b, a uint8
	a = 255 // Default alpha

	var err error
	switch len(hexColor) {
	case 6: // #RRGGBB
		_, err = fmt.Sscanf(hexColor, "%02x%02x%02x", &r, &g, &b)
	case 8: // #RRGGBBAA
		_, err = fmt.Sscanf(hexColor, "%02x%02x%02x%02x", &r, &g, &b, &a)
	default:
		return fmt.Errorf("invalid hex color format: %s", value.Value)
	}

	if err != nil {
		return fmt.Errorf("error parsing hex color '%s': %w", value.Value, err)
	}

	*c = Color{R: r, G: g, B: b, A: a}
	return nil
}

type Config struct {
	Romm struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
	} `yaml:"romm"`
	Theme struct {
		FontPath                        string `yaml:"font_path"`
		TitleColor                      Color  `yaml:"title_color"`
		TitleBackgroundColor            Color  `yaml:"title_background_color"`
		TitleLineColor                  Color  `yaml:"title_line_color"`
		TitleFontSize                   int    `yaml:"title_font_size"`
		TitleLineWidth                  int    `yaml:"title_line_width"`
		BackgroundColor                 Color  `yaml:"background_color"`
		TextColor                       Color  `yaml:"text_color"`
		ListSelectedTextColor           Color  `yaml:"list_selected_text_color"`
		ListTextFontSize                int    `yaml:"list_text_font_size"`
		ListSelectedTextBackgroundColor Color  `yaml:"list_selected_text_background_color"`
		ListBackgroundColor             Color  `yaml:"list_background_color"`
		FooterTextFontSize              int    `yaml:"footer_text_font_size"`
		FooterTextColor                 Color  `yaml:"footer_text_color"`
		FooterLineColor                 Color  `yaml:"footer_line_color"`
		FooterBackgroundColor           Color  `yaml:"footer_background_color"`
		ImageBackgroundColor            Color  `yaml:"image_background_color"`
		TextBackgroundColor             Color  `yaml:"text_background_color"`
		TextPadding                     int    `yaml:"text_padding"`
	} `yaml:"theme"`
	System struct {
		MaxFPS       int    `yaml:"max_fps"`
		ShowFPS      bool   `yaml:"show_fps"`
		ImageTmpPath string `yaml:"image_tmp_path"`
	} `yaml:"system"`
}

func New(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
