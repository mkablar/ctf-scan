
<h1 style="text-align: center;">CTF Scan</h1>

A fast, but thorough scanner for all your CTF or work needs!

## About the project

This scanner does its job in a couple of steps:
1. Quick masscan to get the open ports
2. Nmap scan on the discovered ports (saved to nmap.txt) - so you can start hacking in a few minutes
3. A full nmap scan on all the ports (saved to large-nmap.txt) - to be sure we did not miss anything
4. A udp nmap scan on the top 1000 ports (saved to udp-nmap.txt) - if the -u flag is specified

## Getting Started

Instructions on how to setup the project.

### Prerequisites

To build and run this program, you need to have golang installed on the machine. You can check if it is present with
```sh
go version
```
with the response showing you the version similar to
```sh
go version go1.18.3 linux/amd64
```

### Installation

1. Clone the repo
```sh
git clone https://github.com/mkablar/ctf-scan.git
```
2. Cd into the ctf-scan directory
```sh
cd ctf-scan
```
3. Change the desired file names, located at the top of the file, with the editor of your choice
```go
var nmapFileName = "nmap.txt"
var largeNmapFileName = "large-nmap.txt"
var udpNmapFileName = "udp-nbmap.txt"
```
4. Build and place the file in /usr/bin
```sh
sudo go build -o /usr/bin ctfscan.go
```

## Usage

1. The program needs to be run as root
2. Ip address is required
3. -i flag is to specify a network interface, default eth0 if not provided
4. -u flag is to specify that you want to run a udp scan, after the tcp scan finishes

### Example
```sh
sudo ctfscan 10.10.10.3 -u -i tun0
```
This will run a scan on the 10.10.10.3 ip, using the network interface tun0, and will do a udp scan as well.