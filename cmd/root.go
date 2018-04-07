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
	"os"
	"path"
	"path/filepath"

	"os/signal"
	"syscall"

	"github.com/apex/log"
	cliHandler "github.com/apex/log/handlers/cli"
	textHandler "github.com/apex/log/handlers/logfmt"
	multiHandler "github.com/apex/log/handlers/multi"
	"github.com/bullettime/lora-mqtt/database"
	"github.com/bullettime/lora-mqtt/database/influxdb"
	"github.com/bullettime/lora-mqtt/input"
	"github.com/bullettime/lora-mqtt/parser"
	"github.com/bullettime/lora-mqtt/parser/ttnjson"
	"github.com/bullettime/lora-mqtt/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	logFile    *os.File
	verbose    bool
	debug      bool
	metricName string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lora-mqtt",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var logLevel = log.InfoLevel
		var logHandlers []log.Handler

		if verbose {
			logHandlers = append(logHandlers, cliHandler.Default)
		}

		if debug {
			logLevel = log.DebugLevel
		}

		absLogFileLocation, err := filepath.Abs("lora-mqtt.log")
		if err != nil {
			panic(err)
		}
		logFile, err = os.OpenFile(absLogFileLocation, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		if err == nil {
			logHandlers = append(logHandlers, textHandler.New(logFile))
		}

		log.SetHandler(multiHandler.New(logHandlers...))
		log.SetLevel(logLevel)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		checkConfig()
		start()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if logFile != nil {
			logFile.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lora-mqtt.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print everything to standard output")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug logs")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringVarP(&metricName, "metric-name", "m", ttnjson.LocationData, "define custom metric name")

	viper.SetDefault("influxdb.precision", "ms")
	viper.SetDefault("mqtt.clientid", fmt.Sprintf("lora-mqtt-%s", util.RandomString(4)))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		configName := ".lora-mqtt"

		// Search config in home directory with name ".lora-mqtt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(configName)

		// Add standard path to config file
		cfgFile = path.Join(home, string(append([]byte(configName), []byte(".yaml")...)))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func checkConfig() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("no config file found (run 'configure' first)")
	}
}

func start() {
	if len(metricName) == 0 {
		log.Fatal("you need to specify a valid metric name")
	}
	log.WithField("name", metricName).Debug("metric")

	p, err := ttnjson.New(metricName)
	if err != nil {
		log.WithError(err).Fatal("can't create ttnjson parser")
	}

	influxOptions := influxdb.InfluxOptions{
		Server:    viper.GetString("influxdb.server.url"),
		Username:  viper.GetString("influxdb.server.username"),
		Password:  viper.GetString("influxdb.server.password"),
		Database:  viper.GetString("influxdb.database"),
		Precision: viper.GetString("influxdb.precision"),
	}
	log.WithFields(log.Fields{
		"Server":    influxOptions.Server,
		"Username":  influxOptions.Username,
		"Database":  influxOptions.Database,
		"Precision": influxOptions.Precision,
	}).Debug("InfluxDB Options")
	db := influxdb.New(influxOptions)

	err = db.Connect()
	if err != nil {
		log.WithError(err).Fatal("can't connect to influxdb")
	}
	defer db.Close()
	log.WithFields(log.Fields{
		"server":   influxOptions.Server,
		"database": influxOptions.Database,
	}).Info("connected to influxdb")

	mqttOptions := input.MQTTOptions{
		Server:   viper.GetString("mqtt.server.url"),
		Username: viper.GetString("mqtt.server.username"),
		Password: viper.GetString("mqtt.server.password"),
		QoS:      viper.GetInt("mqtt.qos"),
		ClientID: viper.GetString("mqtt.clientid"),
		Debug:    viper.GetBool("mqtt.debug"),
	}
	log.WithFields(log.Fields{
		"Server":   mqttOptions.Server,
		"Username": mqttOptions.Username,
		"QoS":      mqttOptions.QoS,
		"ClientID": mqttOptions.ClientID,
		"Debug":    mqttOptions.Debug,
	}).Debug("MQTT Options")
	mqtt := input.New(mqttOptions)

	err = mqtt.Connect()
	if err != nil {
		log.WithError(err).Fatal("can't connect to mqtt")
	}
	defer mqtt.Close()
	log.WithField("server", mqttOptions.Server).Info("connected to mqtt")

	err = mqtt.Subscribe(viper.GetString("mqtt.topic"))
	if err != nil {
		log.WithError(err).Fatalf("can't subscribe to topic: %s", viper.GetString("mqtt.topic"))
	}
	log.WithField("topic", viper.GetString("mqtt.topic")).Info("mqtt subscribed to topic")

	go receiver(mqtt, p, db)

	waitForSignal()
}

func receiver(m *input.MQTT, p parser.Parser, db database.Database) {
	for {
		select {
		case <-m.Done:
			log.Debug("shutting down receiver")
			return
		case msg := <-m.Incoming:
			log.WithFields(log.Fields{
				"topic":   msg.Topic(),
				"payload": string(msg.Payload()),
			}).Debug("received message")
			metrics, err := p.Parse(msg.Payload())
			if err != nil {
				log.WithError(err).Warnf("could not parse payload: %s", string(msg.Payload()))
				continue
			}

			err = db.Write(metrics)
			if err != nil {
				log.WithError(err).Error("could not write metrics to database")
			}
		}
	}
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.WithField("signal", s).Warn("exiting")
}
