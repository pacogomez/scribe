package cli

import (
	"fmt"
	"github.com/adrg/xdg"
	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Globals struct {
	Version   VersionFlag  `name:"version" short:"v" help:"Print version information and quit"`
	Config    ConfigCmd    `cmd:"" help:"Display current config" default:"0"`
	Dashboard DashboardCmd `cmd:""`
}

type ConfigCmd struct {
}

func (cmd *ConfigCmd) Run(globals *Globals) error {
	//log.Println("Home config directory:", xdg.ConfigHome)
	appNane := "scribe"
	configFilePath, err := xdg.ConfigFile("" + appNane + "/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Save the config file at:", configFilePath)

	myConfig := Config{DataFile: "~/rnd/private-data/user.phr"}
	d, err := yaml.Marshal(&myConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = os.WriteFile(configFilePath, d, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	configFilePath, err = xdg.SearchConfigFile("" + appNane + "/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Config file was found at:", configFilePath)

	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Printf("file: %s\n", config.DataFile)
	return nil
}

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type Config struct {
	DataFile string `yaml:"data-file"`
}

type CLI struct {
	Globals
}

func Start() {

	cli := CLI{
		Globals: Globals{
			Version: VersionFlag("0.1.1"),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("scribe"),
		kong.Description("A scribe CLI"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.0.1",
		})
	err := ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)

}
