package main

import (
	"time"
	"bytes"
	"log"
	"runtime"
	"os"
	"io"
	"os/exec"
	"bufio"
	"fmt"
	"github.com/gotmc/libusb"
	"github.com/Syfaro/telegram-bot-api"
	"strings"
	"strconv"
)

var alarmstate = 0
var prevnum = 0
var svls []string = []string(GetSerials())
var mainconf = "/etc/usbtripwire.conf"
var linlogfile string = "/var/log/usbtripewire.txt"
var tapi = ""
var tusr = ""
var tusers []string
var tmsg = ""
var clist []string
type Config map[string]string
// ===========================================
func main() {
	if _, err := os.Stat(mainconf); os.IsNotExist(err){
		log.Printf("File: %s not found!\n", mainconf)
		SaveToLog("Main config file not found!\n")
		os.Exit(3)
	}
	cfg, err := ReadConfig(mainconf)
	if err != nil {
        fmt.Println(err)
        os.Exit(3)
    }
    cmdline := cfg["cmd"]
	clist = strings.Split(cmdline, ";")
	tapi = cfg["telegram_apikey"]
	tmsg = cfg["telegram_alarmtext"]
	tusr = cfg["telegram_users"]
    tusers = strings.Split(tusr, ";")
   
	for {
		CheckUsb()
		time.Sleep(time.Millisecond * time.Duration(1000))
	}
}

func CheckUsb() bool {
	ctx, err := libusb.NewContext()
	if err != nil {
        log.Fatal("Couldn't create USB context. Ending now.")
    }
	defer ctx.Close()

    devices, _ := ctx.GetDeviceList()
    if err != nil {
        log.Fatalf("Couldn't get devices")
    }
    snList := []string{}
    for _, device := range devices {
        usbDeviceDescriptor, err := device.GetDeviceDescriptor()
        if err != nil {
            continue
        }
        handle, err := device.Open()
        if err != nil {
            continue
        }
        defer handle.Close()
        serialNumber, err := handle.GetStringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
        if err == nil {
            snList = append(snList, serialNumber)
        }
    }
	if(prevnum != len(snList)){
		if(len(svls) > 1){
			if(ConstSlice(snList, svls)){
				if(alarmstate == 0){					
					RunCmd()
					if(tapi != "" && tmsg != "" && tusr != ""){
					    SendToTelega(tmsg)
					}
				}
				alarmstate = 1
			}else{
				alarmstate = 0
			}
		}else{			
			RunCmd()
			if(tapi != "" && tmsg != "" && tusr != ""){
			    SendToTelega("Something plug USB-device into PC")
			}
		}
	}
	prevnum = len(snList)
	// ------------------------------------
	return true
}

func GetSerials() []string {
	cfg, err := ReadConfig(mainconf)
	if err != nil {
        fmt.Println(err)
        os.Exit(3)
    }
	devlist := cfg["devlist"]
	slist := strings.Split(devlist, ";")
	return slist
}

func ConstSlice(a []string, x []string) bool {
	for _, n := range a {
		for _, m := range x {
			if(n == m){
			    return true
			}
        }
    }
	return false
}

func Contains(a []string, x string) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}
func RunCmd() bool {
	for _, eachline := range clist {
		cmd := exec.Command("bash", "-c", eachline);
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", eachline);
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
	return true
}
func SaveToLog(msglog string) bool{
	os.OpenFile(linlogfile, os.O_RDONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile(linlogfile, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
	currentTime := time.Now()
	currentTime.Format("yyyy-MM-dd HH:mm:ss")
    if _, err = f.WriteString(currentTime.String() + " " + msglog); err != nil {
		panic(err)
    }
    return true
}
func ReadConfig(filename string) (Config, error) {
    config := Config{
        " ":     " ",
    }
    if len(filename) == 0 {
        return config, nil
    }
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := bufio.NewReader(file)

    for {
        line, err := reader.ReadString('\n')
        if equal := strings.Index(line, "="); equal >= 0 {
            if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
                value := ""
                if len(line) > equal {
                    value = strings.TrimSpace(line[equal+1:])
                    value = strings.Trim(value, `"`)
                }
                config[key] = value
            }
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
    }
    return config, nil
}
func SendToTelega(sendmsg string) {
    bot, err := tgbotapi.NewBotAPI(tapi)
    if err != nil {
		log.Panic(err)
	}
	for _, usr := range tusers {
	    if usrn, err := strconv.ParseInt(usr, 10, 64); err == nil {
	        //log.Println(usrn)
            msg := tgbotapi.NewMessage(usrn, sendmsg)
	        bot.Send(msg)
        }
	    
	}
}
