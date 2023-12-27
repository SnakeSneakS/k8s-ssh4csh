package config

import (
	env "github.com/Netflix/go-env"
)

type Ssh4CurrentCsh_ServerEnv struct {
	CONFIG struct {
		//TARGET_CONTAINER string `env:"TARGET_CONTAINER,required=true"`
		SSH_PUB_KEY string `env:"SSH_PUB_KEY,required=true"`
		SERVER_PORT string `env:"PORT,default=22,required=false"`
		SHELL       string `env:"SHELL,default=sh,required=false"`
	}

	TARGET_CONTAINER string `env:"TARGET_CONTAINER,required=true"`
}

func LoadSsh4CshServerEnv() (Ssh4CurrentCsh_ServerEnv, error) {
	var sshSidecarEnv Ssh4CurrentCsh_ServerEnv
	_, err := env.UnmarshalFromEnviron(&sshSidecarEnv)
	if err != nil {
		return sshSidecarEnv, err
	}
	return sshSidecarEnv, nil
}
