package config

import "specommerce/orderservice/pkg/service_config"

type AppConfig struct {
	Server                 service_config.RestServiceConfig `koanf:"server"`
	Env                    string                           `koanf:"env"`
	Database               service_config.DbConfig          `koanf:"db"`
	Kafka                  service_config.KafkaConfig       `koanf:"messagequeue"`
	ProcessPaymentRequest  service_config.KafkaConfig       `koanf:"processPaymentRequest"`
	ProcessPaymentResponse service_config.KafkaConfig       `koanf:"processPaymentResponse"`
	OrderEvents            service_config.KafkaConfig       `koanf:"orderEvents"`
}
