package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type Metadata struct {
	Title    string  `json:"title"`
	Software string  `json:"software"`
	Source   string  `json:"source"`
	Prompts  string  `json:"prompts"`
	Uc       string  `json:"uc"`
	Steps    int64   `json:"steps"`
	Strength float64 `json:"strength"`
	Seed     int64   `json:"seed"`
	Scale    int64   `json:"scale"`
	Sampler  string  `json:"sampler"`
	Noise    float64 `json:"noise"`
}

// metadataCmd represents the metadata command
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			panic(err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		pmp := pngstructure.NewPngMediaParser()
		intfc, err := pmp.ParseFile(filePath)
		if err != nil {
			panic(err)
		}

		cs := intfc.(*pngstructure.ChunkSlice)
		index := cs.Index()
		chunks, found := index["tEXt"]
		if found != true {
			panic("NovelAI metadata not found")
		}

		set := make(map[string]string, len(chunks))

		for _, c := range chunks {
			data := bytes.Split(c.Data, []byte("\x00"))
			title := string(data[0])
			value := string(data[1])
			set[title] = value
		}

		metadata := Metadata{
			Title:    set["Title"],
			Software: set["Software"],
			Source:   set["Source"],
			Prompts:  set["Description"],
			Uc:       gjson.Get(set["Comment"], "uc").String(),
			Steps:    gjson.Get(set["Comment"], "steps").Int(),
			Strength: gjson.Get(set["Comment"], "strength").Float(),
			Seed:     gjson.Get(set["Comment"], "seed").Int(),
			Scale:    gjson.Get(set["Comment"], "scale").Int(),
			Sampler:  gjson.Get(set["Comment"], "sampler").String(),
			Noise:    gjson.Get(set["Comment"], "noise").Float(),
		}

		bytes, err := json.Marshal(metadata)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metadataCmd.PersistentFlags().String("foo", "", "A help for foo")
	metadataCmd.PersistentFlags().String("file", "", "")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metadataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
