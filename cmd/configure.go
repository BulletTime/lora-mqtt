// Copyright © 2018 Sven Agneessens <sven.agneessens@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/segmentio/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
)

type yamlConfig struct {
	InfluxDB influxdbConfig `yaml:"influxdb"`
	MQTT     mqttConfig     `yaml:"mqtt"`
}

type serverConfig struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type influxdbConfig struct {
	Server    serverConfig `yaml:"server"`
	Database  string       `yaml:"database"`
	Precision string       `yaml:"precision"`
}

type mqttConfig struct {
	Server   serverConfig `yaml:"server"`
	QoS      int          `yaml:"qos"`
	ClientID string       `yaml:"clientid"`
	Topic    string       `yaml:"topic"`
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure lora-mqtt",
	Long: `lora-mqtt configure creates a yaml configuration file for the mqtt tool.

Various different values for settings that are needed to use this tool are asked.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("configure called")

		influx := setupInflux()
		mqtt := setupMQTT()

		newConfig := &yamlConfig{
			InfluxDB: influx,
			MQTT:     mqtt,
		}

		output, err := yaml.Marshal(newConfig)
		if err != nil {
			log.WithError(err).Fatal("failed generating yaml config")
		}

		if len(viper.ConfigFileUsed()) == 0 {
			viper.SetConfigFile(cfgFile)
		}

		f, err := os.Create(viper.ConfigFileUsed())
		if err != nil {
			log.WithError(err).Fatal("failed creating log file")
		}

		defer f.Close()

		f.Write(output)
		log.WithField("path", viper.ConfigFileUsed()).Debug("new configuration file saved")
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printHeader(text string) {
	fmt.Printf("[===] %s [===]\n", text)
}

func printFooter() {
	fmt.Printf("\n")
}

func setupInflux() influxdbConfig {
	var config influxdbConfig
	var name = "InfluxDB"

	printHeader("Configure InfluxDB")
	defer printFooter()

	setupServer(&config.Server, name)

	config.Database = prompt.StringRequired("[%s] database (required)", name)
	config.Precision = prompt.String("[%s] precision (default `ms`)", name)

	return config
}

func setupMQTT() mqttConfig {
	var config mqttConfig
	var name = "MQTT"

	printHeader("Configure MQTT")
	defer printFooter()

	setupServer(&config.Server, name)

	config.QoS = prompt.Choose("[MQTT] Quality of Service (default `0`)", []string{"0", "1", "2"})
	config.ClientID = prompt.String("[%s] Client ID", name)
	config.Topic = prompt.StringRequired("[%s] Topic (eg. `+/devices/+/up`)", name)

	return config
}

func setupServer(config *serverConfig, name string) {
	for !isValidServer(config.Url) {
		config.Url = prompt.StringRequired("[%s] server in `scheme://host:port` format (required)", name)
	}

	config.Username = prompt.String("[%s] username", name)
	config.Password = prompt.PasswordMasked("[%s] password", name)
}

func isValidServer(server string) bool {
	u, err := url.Parse(server)
	if err != nil {
		log.WithError(err).Warn("invalid url")
		return false
	}

	if u.IsAbs() && len(u.Port()) > 0 {
		return true
	}

	return false
}
