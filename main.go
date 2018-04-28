package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/smtp"
	"strings"
)

type Info struct {
	Nickname string `json:nickname`
	Subject  string `json:subject`
}

type Email struct {
	To       []string `json:"to"`
	User     string   `json:"user"`
	Password string   `json:"password"`
	Host     string   `json:"host"`
	Port     string
	Info     `json:"info"`
}

func main() {
	ip, err := getIP()
	if err != nil {
		log.Fatalln(err)
	}

	config, err := getConfig()
	if err != nil {
		log.Fatalln(err)
	}
	var e Email
	json.Unmarshal(config, &e)

	err = sendMail(ip, e)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("send email success")
}

func getConfig() ([]byte, error) {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	return b, err
}

// 获取IP地址
func getIP() ([]string, error) {
	var ip []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = append(ip, ipnet.IP.String())
			}
		}
	}
	if ip == nil {
		return nil, errors.New("获取IP失败")
	}
	return ip, nil
}

func sendMail(ip []string, email Email) error {
	auth := smtp.PlainAuth("", email.User, email.Password, email.Host)

	smtpAddr := email.Host + ":" + email.Port
	contentType := "Content-Type: text/plain; charset=UTF-8"
	body := strings.Join(ip, ",")

	msg := []byte("To:" + strings.Join(email.To, ",") + "\r\nFrom:" + email.Nickname + "<" + email.User + ">\r\nSubject: " + email.Subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail(smtpAddr, auth, email.User, email.To, msg)
	if err != nil {
		return err
	}
	return nil
}
