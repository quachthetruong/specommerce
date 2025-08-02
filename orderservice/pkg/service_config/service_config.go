package service_config

type DbConfig struct {
	User            string `koanf:"user"`
	Password        string `koanf:"password"`
	DbName          string `koanf:"dbName"`
	Port            string `koanf:"port"`
	Host            string `koanf:"host"`
	EnableSsl       bool   `koanf:"enableSsl"`
	AutoMigrate     bool   `koanf:"autoMigrate"`
	EnableQueryHook bool   `koanf:"enableQueryHook"`
}

type KafkaConfig struct {
	Host            string `koanf:"host"`
	Retry           int    `koanf:"retry"`
	AutoCreateTopic bool   `koanf:"autoCreateTopic"`
	Topic           string `koanf:"topic"`
	ConsumerGroup   string `koanf:"consumerGroup"`
}

// GrpcServiceConfig defines the configuration for gRPC services
type GrpcServiceConfig struct {
	Endpoint      string            `koanf:"endpoint" yaml:"endpoint" required:"true"`
	Version       string            `koanf:"version" yaml:"version,omitempty" default:"v0.1.0"`
	Metadata      map[string]string `koanf:"metadata" yaml:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ServiceConfig string            `koanf:"serviceconfig" yaml:"serviceConfig,omitempty"`
	Name          string            `koanf:"name" yaml:"name,omitempty"`
}

type RestServiceConfig struct {
	Port int    `koanf:"port" yaml:"port" required:"true"`
	Name string `koanf:"name" yaml:"name" required:"true"`
}
