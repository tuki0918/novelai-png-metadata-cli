package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	pngstructure "github.com/dsoprea/go-png-image-structure/v2"
	"github.com/spf13/cobra"
)

// metadataCmd represents the metadata command
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "",
	Long:  ``,
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
			metadata := bytes.Split(c.Data, []byte("\x00"))
			title := string(metadata[0])
			data := string(metadata[1])
			set[title] = data
		}

		bytes, err := json.Marshal(set)
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
