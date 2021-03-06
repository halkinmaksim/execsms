package main

import (
	"fmt"

	"strings"
	"os"
	"log"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"os/exec"
	"path/filepath"
	"time"
)

/*
compiller for DragonBoard

set GOARCH=arm64
set GOOS=linux
go build -o execsms

compiller for Raspberry

set GOARCH=arm
set GOOS=linux
go build -o execsms


sudo wvdial -C /home/rxhf/risinghf/me909/dial.conf
enable text SMS mode
AT+CMGF=1		text mode on
AT+CSCS="GSM"
AT+CSCB=1

AT+CMGR=1
AT+CMGD=1

disable out
		^HCSQ: "LTE",35,29,21,26

		^RSSI: 13

		^HCSQ: "LTE",35,29,71,26

AT^CURC=0

Manufacturer: Huawei Technologies Co., Ltd.
Model: ME909s-120
Revision: 11.617.01.00.00
IMEI: 867377020177147
+GCAP: +CGSM,+DS,+ES



/opt/risinghf/pktfwd

*/



const CLR_0 = "\x1b[30;1m"

const CLR_R = "\x1b[31;1m"

const CLR_G = "\x1b[32;1m"

const CLR_Y = "\x1b[33;1m"

const CLR_B = "\x1b[34;1m"

const CLR_M = "\x1b[35;1m"

const CLR_C = "\x1b[36;1m"

const CLR_W = "\x1b[37;1m"

const CLR_N = "\x1b[0m"

var c_str_PASS = fmt.Sprintf("%sPASS%s",CLR_G,CLR_N)
var c_str_FAIL = fmt.Sprintf("%sFAIL%s",CLR_R,CLR_N)
var c_str_WARNING = fmt.Sprintf("%sWARNING%s",CLR_Y,CLR_N)


type prog_settings struct {
	FtpServer 			string
	FtpUser				string
	FtpPass				string
	WiFiSSID			string
	WiFiPASS			string
	ModemPort			string
	AddrConfig			string
}
func(pr_set *prog_settings)LoadFromFile(fileName string)bool{
	rawDataIn, err := ioutil.ReadFile(fileName)
	if err != nil {
		return false
		log.Fatal("Cannot load settings:", err)
	}
	err = json.Unmarshal(rawDataIn, &pr_set)
	if err != nil {
		return false
		log.Fatal("Invalid settings format:", err)
	}
	return true
}

/*
RESET LTE
REBOOT GATEWAY

RESET LTE - это sudo systemctl restart lte
REBOOT GATEWAY - это под суперюзером sudo reboot
SET SERVER: address port

sudo nano local_conf.json

{
    "gateway_conf": {
        "gateway_ID": "b827ebFFFE844c44",
        "server_address": "185.41.186.74",
        "serv_port_up": 1700,
        "serv_port_down": 1700
    }
}

*/
func main() {
	//fmt.Println("dd")
	fmt.Println("Service SMS command")
	var prog_sett prog_settings
	fn,_:= os.Executable()
	exPath := filepath.Dir(fn)
	exPath = exPath+"/programsettings.json"
	//fmt.Println(exPath)




	//create file for log
	f, err := os.OpenFile("/tmp/smsservice.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer to close
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)
	log.Println("\n\tStart Application")


	//	Load program settings
	if prog_sett.LoadFromFile(exPath)!=true{
		fmt.Println("Can't load programm settings from file programsettings.json")
	}

	LTEModem:=param_test_modem{"LTE",prog_sett.ModemPort,"","","","","","",""}
	var s []string



	/*
	s:="+CMGR: 1,,99\r\n sadsasd "
	fmt.Println(s)
	ars:=strings.Split(s,",")
	fmt.Println("ars",ars)
	fmt.Println("len(ars) = ",len(ars))
	//ars= append(ars[:1],ars[2:]...)
	ars = ars[:1+copy(ars[1:], ars[1+1:])]
	fmt.Println("ars",ars)
	fmt.Println("len(ars) = ",len(ars))



	fmt.Println(strings.HasPrefix(s," CMGR"))
	fmt.Println(strings.HasPrefix(s,"+CMGR"))
	fmt.Println(strings.HasSuffix(s," CMGR"))
	fmt.Println(strings.HasSuffix(s,"+CMGR"))

	fmt.Println(strings.TrimPrefix(s,"+CMGR"))
	fmt.Println(strings.TrimSuffix(s,"+CMGR"))
	fmt.Println(strings.TrimPrefix(s,"CMGR"))
	fmt.Println(strings.TrimSuffix(s,"CMGR"))

	str_buff := make([]byte, 0)
	str_buff = append(str_buff, '\r')
	str_buff = append(str_buff, '\n')

	fmt.Println(bytes.HasSuffix(str_buff,[]byte("\r\n")))


	*/
	//return
	const LINUX_REBOOT_MAGIC1 uintptr = 0xfee1dead
	const LINUX_REBOOT_MAGIC2 uintptr = 672274793
	const LINUX_REBOOT_CMD_RESTART uintptr = 0x1234567

	fmt.Println(LTEModem.AT_GetCSQ())
	fmt.Println(LTEModem.AT_ConfigSMS())

		for i:=1;i<11;i++{
			s,_=LTEModem.AT_ReadSMS(i)
			fmt.Println("Text SMS - ",s[6])
			if(strings.HasPrefix(s[6],"REBOOT GATEWAY")){
				fmt.Println("REBOOT GATEWAY")
				log_str:=fmt.Sprintf("\tReceived SMS %s, from num: %s",s[6],s[2])
				log.Println(log_str)
				fmt.Println(log_str)
				LTEModem.AT_DelSMS(i)
				cmdec := exec.Command("/sbin/reboot")
				cmdec.Run()

			}
			if(strings.HasPrefix(s[6],"RESET LTE")){
				fmt.Println("RESET LTE")
				log_str:=fmt.Sprintf("\tReceived SMS %s, from num: %s",s[6],s[2])
				log.Println(log_str)
				fmt.Println(log_str)
				cmdec := exec.Command("systemctl","restart","lte")
				cmdec.Run()
				time.Sleep(time.Second*2)
			}
			if(strings.HasPrefix(s[6],"SET SERVER:")){
				fmt.Println("SET SERVER:")
				log_str:=fmt.Sprintf("\tReceived SMS %s, from num: %s",s[6],s[2])
				log.Println(log_str)
				fmt.Println(log_str)

				serv_sett:=strings.Split(strings.TrimSpace(strings.TrimPrefix(s[6],"SET SERVER:"))," ")
				fmt.Println(serv_sett)
				serv_num,_:=strconv.Atoi(strings.TrimSpace(serv_sett[1]))
				//prog_sett.AddrConfig
				ChangeServer(prog_sett.AddrConfig,strings.TrimSpace(serv_sett[0]),serv_num)
				cmdec := exec.Command("systemctl","restart","pktfwd")
				cmdec.Run()
				time.Sleep(time.Second*2)
			}
			LTEModem.AT_DelSMS(i)
		}
		//time.Sleep(time.Second*20)


/*
	fmt.Println(LTEModem.AT_ReadSMS(2))
	fmt.Println("-------------------------")
	fmt.Println(LTEModem.AT_ReadSMS(3))
	fmt.Println("-------------------------")
	fmt.Println(LTEModem.AT_ReadSMS(4))
	fmt.Println("-------------------------")*/


	//waitsms()

}

func ChangeServer( configfile string,serv_addr string,serv_port int)  {
	config_json := map[string]interface{}{}
	rawDataIn, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Println("Cannot load settings:", err)
	}
	err = json.Unmarshal(rawDataIn, &config_json)
	if err != nil {
		log.Println("Invalid settings format:", err)
	}

	log.Println(config_json)

	config_json["gateway_conf"].(map[string]interface{})["server_address"]=serv_addr
	config_json["gateway_conf"].(map[string]interface{})["serv_port_up"]=serv_port

	rawDataOut, err := json.MarshalIndent(&config_json, "", "  ")
	if err != nil {
		log.Println("JSON marshaling failed:", err)
	}
	err = ioutil.WriteFile(configfile, rawDataOut, 0)
	if err != nil {
		log.Println("Cannot write updated settings file:", err)
	}
	log.Println(config_json)
}
