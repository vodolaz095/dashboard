package config

type Broadcaster struct {
	// ConnectionName defines redis/mqtt3/mqtt5 connection name we broadcast model.Updates into
	ConnectionName string `yaml:"connection_name"`
	// SensorsToListen defines array of sensors names, from which we do publish updates, if left
	// blank - all sensors' readings will be broadcasted
	SensorsToListen []string `yaml:"sensors_to_listen"`
	// Subject defines routing parameters for publishing updates, currently it is tempalte with sensor name
	// like `vodolaz095/dashboard/{{sensorName}}`
	Subject string `yaml:"subject"`
	// ValueOnly defines way we send model.Update - raw value in decimal float encoding, or do we send JSON message
	ValueOnly bool `yaml:"value_only"`
}
