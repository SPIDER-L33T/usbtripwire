package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"bytes"
	"strings"
	"time"
	"os"
	"os/exec"
	"net/http"
	"github.com/dustin/go-humanize"
	"github.com/gotmc/libusb"
)

var conf = []string{}
var devlist string = ""
type WriteCounter struct {
	Total uint64
}

func main() {
    cmd := exec.Command("bash", "-c", "killall usbtripwire; rm /usr/local/bin/usbtripwire; rm /etc/usbtripwire.conf; rm /etc/systemd/system/usbtripwire.service")
    var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

    fmt.Println("Download Started")

	fileUrl := "https://github.com/SPIDER-L33T/usbtripwire/releases/download/v1.0/usbtripwire"
	err = DownloadFile("/usr/local/bin/usbtripwire", fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Download Finished")

	Setup()
	var keyinput string
	fmt.Printf("***** Remove your SECRET USB-device! And press 'ENTER' *****\n");
	fmt.Scanln(&keyinput)

	var clearl = []string(GetUsb())

	fmt.Printf("***** Now insert your SECRET USB-device! And press 'ENTER' *****\n");
	fmt.Scanln(&keyinput)
	var secl = []string(GetUsb())

	fmt.Printf("Searching secret device...\n");
	duration := time.Duration(5)*time.Second
	time.Sleep(duration)
	fmt.Printf("Done!\n");
	duration = time.Duration(1)*time.Second
	time.Sleep(duration)

	ConstSlice(secl, clearl);

	if(len(conf) == 0){
		fmt.Printf("==========================\n")
		fmt.Printf("Secret device not found!\n")
		fmt.Printf("==========================\n")
		return
	}

	input, err := ioutil.ReadFile("/etc/usbtripwire.conf")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "devlist") {
			lines[i] = `devlist="`+devlist+`"`
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("/etc/usbtripwire.conf", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("==========================\n")
    fmt.Println(len(conf), "secret USB-device(s) saved!")
	fmt.Printf("==========================\n")
    fmt.Printf("Install complete!\n")
    fmt.Printf("Please check /etc/usbtripwire.conf for tunning service.\n")
}
// ======================================================
func GetUsb() []string {
	ctx, err := libusb.NewContext()
        if err != nil {
                log.Fatal("Couldn't create USB context. Ending now.")
        }
        defer ctx.Close()

        devices, _ := ctx.GetDeviceList()
        if err != nil {
                log.Fatalf("Couldn't get devices")
        }

	snl := []string{}
	for _, device := range devices {
		usbDeviceDescriptor, err := device.GetDeviceDescriptor()
        if err != nil {
                        log.Printf("Error getting device descriptor: %s", err)
                        continue
                }
                handle, err := device.Open()
                if err != nil {
                        log.Printf("Error opening device: %s", err)
                        continue
                }
                defer handle.Close()
                serialNumber, err := handle.GetStringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
                if err == nil {
                        snl = append(snl, serialNumber)
                }
	}
	return snl
}

func Find(a []string, x string) bool {
	for _, n := range a {
                if x == n {
                        return false
                }
        }
	return true
}
func ConstSlice(c []string, s []string) bool {
	for _, n := range c {
		if(Find(s, n)){
			conf = append(conf, n)
			devlist = devlist + n + ";"
		}
    }
    return true
}
func Setup(){ 
	os.OpenFile("/etc/usbtripwire.conf", os.O_RDONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile("/etc/usbtripwire.conf", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
	etcl := `telegram_apikey=""
telegram_alarmtext="Something plug your device into PC"
telegram_users=""
devlist=""
cmd="date >> /root/log.txt"`
	if _, err = f.WriteString(etcl); err != nil {
		panic(err)
    }
	// ================================================
	os.OpenFile("/etc/systemd/system/usbtripwire.service", os.O_RDONLY|os.O_CREATE, 0600)
	f, err = os.OpenFile("/etc/systemd/system/usbtripwire.service", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()
	
	servlines := `[Unit]
Description=UsbTripWire (Small tripwire)
After=syslog.target
After=network.target

[Service]
RestartSec=2s
Type=simple
ExecStart=/usr/local/bin/usbtripwire
ExecReload=/bin/kill -HUP $MAINPID
Restart=always

[Install]
WantedBy=multi-user.target`
	if _, err = f.WriteString(servlines); err != nil {
		panic(err)
    }
    cmd := exec.Command("bash", "-c", "chmod +x /usr/local/bin/usbtripwire; systemctl daemon-reload")
    var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}
func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}
func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}
