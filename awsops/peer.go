package awsops

import (
	"aws-api/setting"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func GetPeerService(publicip, peertype string) string {

	command := "bash"

	args := []string{setting.AppSetting.ScriptsPath + "/awsops-get-peer-service.sh", publicip, peertype}

	cmd := exec.Command(command, args...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	log.Println(stdBuffer.String())

	return stdBuffer.String()
}

func SetPeerService(publicip, runcmd string) string {

	command := "bash"

	args := []string{setting.AppSetting.ScriptsPath + "/awsops-set-peer-service.sh", publicip, runcmd}

	cmd := exec.Command(command, args...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	// Execute the command
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}

	log.Println(stdBuffer.String())

	return stdBuffer.String()
}
