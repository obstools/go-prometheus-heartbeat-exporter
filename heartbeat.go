package heartbeat

// Creates new server with builded configuration based on the yaml-configuration.
// Returns new server pointer and error
func New(configPath string) (*Server, error) {
	var err error

	configuration, err := loadConfiguration(configPath)
	if err != nil {
		return nil, err
	}

	return newServer(configuration), err
}
