//go:build cit

package regression

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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

	d, _ := os.Open("/proc")
	defer d.Close()

	f1 := false
	f2 := false

	for {
		fmt.Println("Inside first loop")
		names, err := d.Readdirnames(0)
		if err == io.EOF {
			fmt.Printf("Error EOF is %v\n", err)
			break
		}
		for _, name := range names {
			if name[0] < '0' || name[0] > '9' {
				continue
			}
			intnum, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}
			fmt.Printf("The number is %v\n", intnum)

			x := "/proc/" + name + "/comm"
			dat, err := os.ReadFile(x)
			if err != nil {
				fmt.Printf("Read File Error is %v\n", err)
			}
			if f1 {
				fmt.Printf("The guest is present checking for %s and  %s\n", x, string(dat))
				n1 := strings.Trim(string(dat), "\n")
				if n1 == "sshd" {
					fmt.Printf("Inside SSH %s", string(dat))
					a1 := strings.Trim(string(name), "\n")
					x1 := "/proc/" + a1 + "/status"
					dat2, err := os.ReadFile(x1)
					if err != nil {
						fmt.Printf("Guest Error is %v\n", err)
					}
					fmt.Printf("The ssh is %s and %s\n", x1, string(dat2))
					f2 = true
					break
				}
			}
			fmt.Printf("The guest is not present checking for it %s and  %s\n", x, string(dat))
			n2 := strings.Trim(string(dat), "\n")
			if n2 == "google_guest_ag" {
				fmt.Printf("Inside Guest %s", string(dat))
				a1 := strings.Trim(string(name), "\n")
				x1 := "/proc/" + a1 + "/status"
				dat2, err := os.ReadFile(x1)
				if err != nil {
					fmt.Printf("Guest Error is %v\n", err)
				}
				fmt.Printf("The guest is %s and %s\n", x1, string(dat2))
				f1 = true
			}
		}
		if f1 && f2 {
			t2 := time.Now()
			fmt.Printf("Time is %s\n", t2.String())

			cmdx := exec.Command("cat", "/proc/uptime")
			outx, err := cmdx.Output()
			if err != nil {
				fmt.Printf("Uptime Error is %v\n", err)
			}
			fmt.Printf("The date is %s\n", string(outx))

			t3 := time.Now()
			fmt.Printf("Time is %s\n", t3.String())
			break
		}
	}
}
