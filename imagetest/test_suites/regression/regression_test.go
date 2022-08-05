//go:build cit

package regression

import (
	"fmt"
	"github.com/GoogleCloudPlatform/guest-test-infra/imagetest/utils"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	PATTERN_KERNEL_TIME      = " in (.*?)s"
	PATTERN_USER_TIME        = " \\+ (.*?)s"
	PATTERN_SYSTEMD_CP_LINE1 = "\\+(.*?)s"
	PATTERN_SYSTEMD_CP_LINE2 = "@(.*?)s"
	NETWORK_TARGET           = "network.target"
	LOCALFS_TARGET           = "local-fs.target"
	SSH_SERVICE              = "ssh.service"
	SSHD_SERVICE             = "sshd.service"
	GUESTSCRIPTS_SERVICE     = "google-startup-scripts.service"
)

// TestKernelStart test
func TestKernelStart(t *testing.T) {
	cmd1 := exec.Command("cat", "/proc/uptime")
	out1, err := cmd1.Output()
	if err != nil {
		t.Fatalf("cat command failed %v", err)
	}
	cmd2 := exec.Command("date", "+%s.%N")
	out2, err := cmd2.Output()
	if err != nil {
		t.Fatalf("date command failed %v", err)
	}
	uptime_epoc := strconv.ParseFloat(string(out2), 64) - strconv.ParseFloat(strings.Split(string(out1), " ")[0], 64)

}

func TestKernelFinish(t *testing.T) {
	cmd := exec.Command("systemd-analyze")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("systemd-analyze command failed %v", err)
	}
	kernel_time, _ := ParseKernelUserTimes(string(out))
	if kernel_time == 0.0 {
		t.Fatalf("Error parsing kernel time")
	}
}

func TestNetworkReady(t *testing.T) {
	cmd := exec.Command("systemd-analyze", "critical-chain", NETWORK_TARGET)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("systemd-analyze command failed %v", err)
	}
	network_time, err := ParseSystemDCriticalChainOutput(string(out))
	if err != nil {
		t.Fatalf("systemd-analyze command failed %v", err)
	}
}

func ParseKernelUserTimes(str string) (float64, float64, error) {
	kernel_time := 0.0
	user_time := 0.0
	kernelRegex, _ = regexp.Complile(PATTERN_KERNEL_TIME)
	userRegex, _ = regexp.Complile(PATTERN_USER_TIME)
	matchedKernel := kernelRegex.MatchString(str)
	matchedUser := userRegex.MatchString(str)
	if matchedKernel {
		kernel_time, err := ParseSeconds(kernelRegex.FindString(str)[4:])
		if err != nil {
			return nil, nil, err
		}
	}
	if matchedUser {
		user_time, err := ParseSeconds(userRegex.FindString(str)[3:])
		if err != nil {
			return nil, nil, err
		}
	}
	return kernel_time, user_time, nil
}

func ParseSeconds(str string) (float64, error) {
	formattedStr := strings.Trim(str, " ")
	secs := 0.0
	for _, element := range strings.Split(formattedStr, " ") {
		if strings.HasSuffix(element, "min") {
			value, _ := strconv.ParseFloat(element[0:len(element)-3], 64)
			secs += (value * 60.00)
		} else if strings.HasSuffix(element, "ms") {
			value, _ := strconv.ParseFloat(element[0:len(element)-2], 64)
			secs += (value / 1000)
		} else if strings.HasSuffix(element, "s") {
			value, _ := strconv.ParseFloat(element[0:len(element)-1], 64)
			secs += value
		} else {
			return nil, fmt.Errorf("error parsing seconds in %s", str)
		}
	}
	return secs, nil
}

func ParseSystemDCriticalChainOutput(str string) (float64, error) {
	lines := strings.Split(str, "\n")
	if len(lines) < 5 {
		return nil, fmt.Errorf("Invalid format. Critical chain command: %s", str)
	}

	systemdLine1Regex, _ := regexp.Complile(PATTERN_SYSTEMD_CP_LINE1)
	systemdLine2Regex, _ := regexp.Complile(PATTERN_SYSTEMD_CP_LINE2)

	firstLine := lines[3]
	secondLine := lines[4]
	totalSecs := 0.0
	foundValue := false

	if strings.Contains(firstLine, "+") {
		boolmatch := systemdLine1Regex.MatchString(firstLine)
		if boolmatch {
			match := systemdLine1Regex.FindString(firstLine)[1:]
			secs, err := ParseSeconds(match)
			if err != nil {
				return nil, err
			}
			totalSecs += secs
			foundValue = true
		}
	} else {
		secondLine = firstLine
	}
	boolmatch := systemdLine2Regex.MatchString(secondLine)
	if boolmatch {
		match := systemdLine2Regex.FindString(secondLine)[1:]
		secs, err := ParseSeconds(match)
		if err != nil {
			return nil, err
		}
		totalSecs += secs
		foundValue = true
	}
	if !foundValue {
		return nil, fmt.Errorf("Invalid format. Critical chain command: %s", str)
	}
	return totalSecs, nil
}
