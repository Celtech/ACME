package configFactory

import (
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type Variable struct {
	name         string
	defaultValue string
	required     bool
}

func getValueFromVariable(variableName string) (string, error) {
	variableValue := os.Getenv(variableName)
	err := (error)(nil)

	if len(variableValue) == 0 {
		err = errors.New(fmt.Sprintf("variable %s not found", variableName))
	}

	return variableValue, err
}

func getVariableFromFile(variableName string) (string, error) {
	val, err := getValueFromVariable(fmt.Sprintf("%s_FILE", variableName))
	if err != nil {
		return "", err
	}

	dat, err := os.ReadFile(val)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func VariableToSetting(logger *zap.Logger, variable Variable) string {
	fileVal, err := getVariableFromFile(variable.name)
	if err == nil {
		return fileVal
	}

	varVal, err := getValueFromVariable(variable.name)
	if err == nil {
		return varVal
	}

	if variable.required {
		logger.Fatal(
			fmt.Sprintf(
				"could not find required environment variable %s or %s_FILE, please set one to start your server",
				variable.name,
				variable.name,
			),
		)
		os.Exit(1)
	}

	return variable.defaultValue
}
