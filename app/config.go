package app

var config = configType{}

type configType struct {
	Name             string
	ShortDescription string
	RootDirectory    string
	GenesisFilePath  string
	DBFilePath       string
}

func GetConfig() configType {
	if !isInit {
		err := Init()
		if err != nil {
			panic(err)
		}
	}

	return config
}
