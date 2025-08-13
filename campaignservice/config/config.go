package config

import "specommerce/campaignservice/pkg/service_config"

type AppConfig struct {
	Server         service_config.RestServiceConfig `koanf:"server"`
	Env            string                           `koanf:"env"`
	Database       service_config.DbConfig          `koanf:"db"`
	Kafka          service_config.KafkaConfig       `koanf:"messagequeue"`
	OrderConsumer  service_config.KafkaConfig       `koanf:"orderEvents"`
	OrderSuccess   service_config.KafkaConfig       `koanf:"orderSuccess"`
	IphoneCampaign string                           `koanf:"iphoneCampaign"`
	Redis          service_config.RedisConfig       `koanf:"redis"`
}
