package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var nmapFileName = "nmap.txt"
var largeNmapFileName = "large-nmap.txt"
var udpNmapFileName = "udp-nmap.txt"

var interfacePtr *string
var udpPtr *bool
var ip string
var stdOutBuf bytes.Buffer
var mw = io.MultiWriter(os.Stdout, &stdOutBuf)

func main() {
	CheckRootUser()
	
	SetFlagUsage()
	ProcessFlags()

	interfacePtr = flag.String("i", "eth0", "network interface to use (default: eth0)")
	udpPtr = flag.Bool("u", false, "scan top 1000 udp ports after finishing tcp scan (default: false)")
	flag.Parse()

	ValidateArguments()
	ip = flag.Args()[0]
	ValidateIpAddress(ip)

	openPorts := GetMasscanOpenPorts()
	if openPorts != "" {
		PrintMasscanResult(openPorts)
		RunNmapOnPorts(openPorts, nmapFileName)
	}

	RunNmapOnAllPorts()

	if *udpPtr {
		RunUpdNmapScan()
	}
}

func IsValueFlag(flag byte) bool {
	switch flag {
	case
		'i':
		return true
	}
	return false
}

func SetFlagUsage() {
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("Usage of %s:\n", "ctfscan")
		PrintDashes()
		fmt.Printf("Ip address is mandatory!\n")
		fmt.Printf("Command needs to be run as root!\n")
		fmt.Printf("Example usage:\tctfscan 192.160.0.1 -i tun0 -u\n")
		PrintDashes()
		order := []string{"i", "u"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			fmt.Printf("-%s =>\t%s\n", flag.Name, flag.Usage)
		}
	}
}

func ProcessFlags() {
	var args []string
	var notargs []string

	for i := 0; i < len(os.Args); i++ {
		if i == 0 {
			notargs = append(notargs, os.Args[i])
		} else if os.Args[i][0] == '-' {
			notargs = append(notargs, os.Args[i])

			// non-bool flags need to have the value provided
			if IsValueFlag(os.Args[i][1]) {
				i++
				notargs = append(notargs, os.Args[i])
			}
		} else {
			args = append(args, os.Args[i])
		}
	}
	os.Args = append(notargs, args...)
}

func ValidateIpAddress(ip string) {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		flag.Usage()
		os.Exit(3)
	}

	for i := 0; i < len(parts); i++ {
		intVar, err := strconv.Atoi(parts[i])
		if err != nil || intVar < 0 || intVar > 255 {
			flag.Usage()
			os.Exit(3)
		}
	}
}

func ValidateArguments() {
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(3)
	}
}

func FilterMasscanResults(output string) string {
	msOutStr := strings.Split(output, "\n")
	var openPorts []string
	for i := 0; i < len(msOutStr); i++ {
		var line = msOutStr[i]
		if strings.Contains(line, "Discovered") {
			openPorts = append(openPorts, FetchPortFromMasscanLine(line))
		}
	}

	openPortsStr := strings.Join(openPorts, ",")
	return openPortsStr
}

func FetchPortFromMasscanLine(line string) string {
	index := strings.Index(line, "Discovered")
	subLine := line[index:len(line)]
	lineSlice := strings.Fields(subLine)
	portSlice := strings.Split(lineSlice[3], "/")
	port := portSlice[0]
	return port
}

func GetMasscanOpenPorts() string {
	fmt.Println("Starting quick masscan for open ports!")
	PrintDashes()

	msCmd := exec.Command("masscan", ip, "-p0-65535", "--rate", "1000", "-e", *interfacePtr)
	msCmd.Stdout = mw
	msCmd.Stderr = mw
	msErr := msCmd.Run()

	if msErr != nil {
		panic(msErr)
	}

	var openPorts = FilterMasscanResults(stdOutBuf.String())
	return openPorts
}

func RunNmapOnPorts(ports string, fileName string) {
	fmt.Printf("Starting nmap scan on %s ports!\n", ports)
	PrintDashes()

	nmapCmd := exec.Command("nmap", "-p", ports, ip, "-A", "-oN", fileName)
	nmapCmd.Stdout = mw
	nmapCmd.Stderr = mw
	nmapErr := nmapCmd.Run()
	if nmapErr != nil {
		panic(nmapErr)
	}

	PrintNmapFinished()
}

func RunNmapOnAllPorts() {
	RunNmapOnPorts("0-65535", largeNmapFileName)
}

func RunUpdNmapScan() {
	fmt.Println("Starting udp nmap scan on top 1000 ports!")
	PrintDashes()

	nmapCmd := exec.Command("nmap", "-sU", ip, "-A", "-oN", udpNmapFileName)
	nmapCmd.Stdout = mw
	nmapCmd.Stderr = mw
	nmapErr := nmapCmd.Run()
	if nmapErr != nil {
		panic(nmapErr)
	}

	PrintNmapFinished()
}

func PrintDashes() {
	fmt.Println("----------------------------------------------------------")
}

func PrintNmapFinished() {
	fmt.Println("=====================================================================")
	fmt.Println("=====================================================================")
	fmt.Println("=============== Nmap finished, check the output file! ===============")
	fmt.Println("=====================================================================")
	fmt.Println("=====================================================================")
}

func PrintMasscanResult(openPorts string) {
	PrintDashes()
	fmt.Printf("Masscan found %s open!\n", openPorts)
	PrintDashes()
}

func CheckRootUser() {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	uid, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		panic(err)
	}

	// root is 0
	if uid != 0 {
		fmt.Printf("Command needs to be run as root!\n")
		os.Exit(3)
	}
}
