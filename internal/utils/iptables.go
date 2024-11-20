package utils

import (
	"log"
	"os/exec"
)

func AddIptablesRule(ip string) {
	cmd := exec.Command("sudo", "iptables", "-t", "nat", "-A", "POSTROUTING", "-s", ip, "-j", "MASQUERADE")
	if err := cmd.Run(); err != nil {
		log.Printf("Error adding iptables rule: %v", ip)
	} else {
		log.Printf("Success adding iptables rule: %v", ip)
	}
}

func DeleteIptablesRule(ip string) {
	cmd := exec.Command("sudo", "iptables", "-t", "nat", "-D", "POSTROUTING", "-s", ip, "-j", "MASQUERADE")
	if err := cmd.Run(); err != nil {
		log.Printf("Error deleting iptables rule: %v", ip)
	} else {
		log.Printf("Success deleting iptables rule: %v", ip)
	}
}
