package config

import "specommerce/campaignservice/pkg/service_config"

type AppConfig struct {
	Server       service_config.RestServiceConfig `koanf:"server"`
	Env          string                           `koanf:"env"`
	Database     service_config.DbConfig          `koanf:"db"`
	Kafka        service_config.KafkaConfig       `koanf:"messagequeue"`
	OrderSuccess service_config.KafkaConfig       `koanf:"orderSuccess"`
}
