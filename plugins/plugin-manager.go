package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var PluginManager pluginManager = pluginManager{
	implementations: map[string]map[string]Plugin{},
	builders:        map[string]PluginBuilder{},
}

type Plugin interface{}

type pluginManager struct {
	implementations map[string]map[string]Plugin
	builders        map[string]PluginBuilder
}

func (pm *pluginManager) GetBuilder(interfaceName string) (PluginBuilder, error) {
	if pm.builders[interfaceName] == nil {
		return nil, fmt.Errorf("interface `%s` is not supported", interfaceName)
	}

	return pm.builders[interfaceName], nil
}

func (pm *pluginManager) LoadAll(pluginsDir string) error {
	pluginConfigs, err := os.ReadDir(pluginsDir)
	if err != nil {
		return err
	}

	for _, d := range pluginConfigs {
		if d.IsDir() {
			continue
		}
		config, err := loadConfig(filepath.Join(pluginsDir, d.Name()))
		if err != nil {
			continue
		}
		for _, ip := range config.IntegrationPoints {
			builder, err := pm.GetBuilder(ip.Interface)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			pluginClient := builder.Build(
				ip.GrpcConnection.Addr,
				ip.GrpcConnection.Port,
			)
			pm.RegisterImplementation(ip.Interface, config.Name, pluginClient)
		}
	}
	return nil
}

func (pm *pluginManager) RegisterInterface(name string, implementatons map[string]Plugin, pluginBuilder PluginBuilder) {
	pm.implementations[name] = implementatons
	pm.builders[name] = pluginBuilder
}

func (pm *pluginManager) RegisterImplementation(interfaceName string, implName string, plugin interface{}) {
	if pm.implementations[interfaceName] == nil {
		fmt.Printf("Extension Ponint %s is not supported by this version of the app.\n", interfaceName)
		return
	}
	pm.implementations[interfaceName][implName] = plugin
}

func loadConfig(configPath string) (*Config, error) {
	var config Config

	viperInstance := viper.NewWithOptions(viper.KeyDelimiter("::"))
	viperInstance.SetConfigFile(configPath)

	if err := viperInstance.ReadInConfig(); err != nil {
		fmt.Println("Can't read config: ", configPath)

		return nil, err
	}

	metaData := &mapstructure.Metadata{}
	if err := viperInstance.Unmarshal(&config, metadataConfig(metaData)); err != nil {
		return nil, err
	}

	if len(metaData.Keys) == 0 || len(metaData.Unused) > 0 {

		return nil, nil
	}

	return &config, nil
}

func metadataConfig(md *mapstructure.Metadata) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.Metadata = md
	}
}
