package oneandone

import (
	"bytes"
	"fmt"
	"github.com/docker/machine/log"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"time"
	"strconv"
	"regexp"
)

// Function to perform busy-waiting for the selected TCP port to open on the first IP of the server.
//
// This functions cycles until the selected TCP port is open. Between each iteration a 5 sec sleep will be done.
func WaitForTcpPortToBeOpen(ip string, port int) {
	target := fmt.Sprintf("%v:%v", ip, port)
	log.Debugf("Wainting for port '%v' to open on IP '%v'.", port, ip)
	_, err := net.DialTimeout("tcp", target, 5*time.Second)
	for err != nil {
		log.Debugf("Port '%v' on IP '%v' still not open, wait 5 sec.", port, ip)
		time.Sleep(5 * time.Second)
		_, err = net.DialTimeout("tcp", target, 5*time.Second)
	}
}

// Function to get an gossh ssh client
//
// This function returns an instance of the gossh ssh client with given parameters
func getSSHClient(user string, ip string, port int, password string) (*gossh.Client, error) {
	sshConfig := &gossh.ClientConfig{
		User: user,
		Auth: []gossh.AuthMethod{gossh.Password(password)},
	}
	target := fmt.Sprintf("%v:%v", ip, port)
	client, err := gossh.Dial("tcp", target, sshConfig)
	if err != nil {
		return nil, err
	}
	return client, err
}

// Function to execute an ssh command
//
// This function executes an ssh command with the given gossh ssh client an the given command
func executeCmd(client *gossh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		log.Info(err)
		return "", err
	}
	var b bytes.Buffer
	session.Stdout = &b
	err = session.Run(cmd)
	if err != nil {
		return "", err
	}
	return b.String(), err
}

// Function to validate that the apt package manager is up to date
//
// This function validates that the apt package manager updates his cache within 30 seconds
func isAptUpToDate(client * gossh.Client) (bool) {
	lastRunS, lastRunErr := executeCmd(client, "stat -c %Y /var/cache/apt/")
	if lastRunErr != nil {
		log.Infof("Failed to fetch last package manager run: %v", lastRunErr)
	}
	currentTimeS, currentTimeErr := executeCmd(client, "date +%s")
	if currentTimeErr != nil {
		log.Infof("Failed to fetch last package manager run: %v", currentTimeErr)
	}

	var validTimestamp = regexp.MustCompile(`\d{10}`)
	lastRunS = validTimestamp.FindString(lastRunS)
	currentTimeS = validTimestamp.FindString(currentTimeS)

	lastRun, lastRunParseErr := strconv.Atoi(lastRunS)
	if lastRunParseErr != nil {
		log.Infof("Failed to parse last update run timestamp: %v", lastRunParseErr)
	}
	currentTime, currentTimeParseErr := strconv.Atoi(currentTimeS)
	if currentTimeParseErr != nil {
		log.Infof("Failed to parse current timestamp: %v", currentTimeParseErr)
	}

	if (currentTime - lastRun) > 30 {
		return true
	}
	return false
}
