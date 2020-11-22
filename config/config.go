package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// struct représentant la config
type ConfigJson struct {
	TraceActivate bool `json:"traceActivate"`
	DebugActivate bool `json:"debugActivate"`
	Processus []ProcessJson `json:"processus"`
}

// struct représentant un processus
type ProcessJson struct {
	Id int `json:"id"`
	Address string `json:"address"`
	Port int `json:"port"`

}

func GetConfig() *ConfigJson {
    // essaye de lire le fichier de config
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}

    // traduit le fichier de config en structure
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config ConfigJson
	json.Unmarshal(byteValue, &config)

    // S'il on est en mode "Debug", on affiche le contenu du fichier
	if config.DebugActivate {
		fmt.Println("Debug : " + strconv.FormatBool(config.DebugActivate))
		fmt.Println("Trace : " + strconv.FormatBool(config.TraceActivate))
        fmt.Println("")

		for i := 0; i < len(config.Processus); i++ {
			fmt.Println("Processus id : " + strconv.Itoa(config.Processus[i].Id))
			fmt.Println("Processus address : " + config.Processus[i].Address)
			fmt.Println("Processus port : " + strconv.Itoa(config.Processus[i].Port))
			fmt.Println("")
		}
	}

	return &config
}