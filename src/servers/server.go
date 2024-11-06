package servers

import "github.com/kylerequez/go-sample-dashboard/src/utils"

func Init() error {
	if err := utils.LoadEnvVariables(); err != nil {
		return err
	}

	hostname, err := utils.GetEnv("SERVER_HOSTNAME")
	if err != nil {
		return err
	}

	port, err := utils.GetEnv("SERVER_PORT")
	if err != nil {
		return err
	}

	app := NewAppServer(*hostname, *port)
	return app.Init()
}
