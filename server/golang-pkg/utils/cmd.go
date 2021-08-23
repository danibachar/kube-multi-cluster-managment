package utils

import (
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/labstack/gommon/log"
)

func runCommand(base, command string) (error, string) {
	log.Info("running command ", command)

	cmd := exec.Command(base, "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Info("StdoutPipe ", err)
		return err, ""
	}
	if err := cmd.Start(); err != nil {
		log.Info("Start ", err)
		return err, ""
	}
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Info("ReadAll ", err)
		return err, ""
	}
	if err := cmd.Wait(); err != nil {
		log.Info("Wait ", err)
		return err, ""
	}
	res := string(data)
	res = strings.TrimSuffix(res, "\n")
	log.Info("finished with res %s", res)
	return nil, res
}

func RunBash(command string) (error, string) {
	return runCommand("/bin/bash", command)
}

func RunShell(command string) (error, string) {
	return runCommand("sh", command)
}
