package main

import (
	"encoding/base64"
	"fmt"
	"github.com/klauspost/cpuid"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const licFile  = "license.key"

func createFile(key string) {
	f, err := os.Create(getPath() + licFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(key)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	infoMsg("File with license has been created, don't remove this file!", true)
}

func readFile() string {
	data, err := ioutil.ReadFile(getPath() + licFile)
	if err != nil {
		fmt.Println("File reading error", err)
		return ""
	}
	return string(data)
}

func getPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(dir,"\\") {
		dir = dir + "\\"
	} else if strings.Contains(dir, "/") {
		dir = dir + "/"
	}
	return dir
}

func getFileModifiedDate(filename string) string {
	// get last modified time
	file, err := os.Stat(filename)

	if err != nil {
		fmt.Println(err)
	}

	modifiedTime := file.ModTime().String()

	return modifiedTime
}

func getCpuInfo() string {
	cpuInfo := cpuid.CPU.BrandName + strconv.Itoa(cpuid.CPU.PhysicalCores) + strconv.Itoa(cpuid.CPU.ThreadsPerCore) + strconv.Itoa(cpuid.CPU.LogicalCores) + strconv.Itoa(cpuid.CPU.Family) + strconv.Itoa(cpuid.CPU.Model)
	return cpuInfo
}

func getUserInfo() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.Name + user.Username
}

func generateSerial(key string) string {
	if len(strings.TrimSpace(key)) > 0 {
		return getCpuInfo() + getUserInfo() + "____" + key
	}
	return getCpuInfo() + getUserInfo() + "____"
}

func salt() string {
	now := time.Now().String()
	return now
}

func encode(key string) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(key))
	return sEnc
}

func decode(key string) string {
	sDec, _ := base64.StdEncoding.DecodeString(key)
	return string(sDec)
}

func updateLicense(key string) {
	if decode(key) > salt() {
		createFile(encode(generateSerial(key)))
	}
	return
}

func install()  {
	key := ""
	_, err := os.Stat(getPath() + licFile)
	if err == nil {
		licenseKey := readFile()
		licenseKeyArr := strings.Split(licenseKey, "___")

		// check if license is valid
		if decode(strings.TrimSpace(licenseKeyArr[0])) == generateSerial("") {
			if decode(licenseKeyArr[1]) < salt() {
				newKey := ""
				infoMsg("You need to update your license, please insert new key: ", true)
				fmt.Scan(&newKey)
				if len(strings.TrimSpace(newKey)) == 0 {
					for i:=3;i>0;i-- {
						infoMsg("Remaining attempts: " + strconv.Itoa(i) + "You need to update your license, please insert new key: ", true)
						fmt.Scan(&newKey)
						if len(strings.TrimSpace(newKey)) > 0 {
							updateLicense(newKey)
						}
					}
				}
			}
		}

	} else {
		createFile(encode(generateSerial(key)))
	}


	infoMsg(" -> Please enter key license: ", true)
	fmt.Scanln(&key)
}
