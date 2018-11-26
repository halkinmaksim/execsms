package main

import (
	"os/exec"
	"fmt"
	"time"

	//"github.com/warthog618/modem/at"
	"strings"
	"github.com/tarm/serial"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/gpio"
	"net"
	"os"
	"bytes"
	"strconv"
)

/*
	20181029
		Проверить ф-цию reboot

*/


type GSM_Cmd struct {
	Cmd 				string
	CmdResponseOnOk		string
	TimeoutMs			int
	DelayMs				int
	Skip				uint8

}


type param_test_modem struct {
	Dev_name 					string
	Port 						string
	GPRS_test					string
	APN							string
	Password					string
	Username					string
	Phone						string
	ATcommand					string
	PhoneSMSandCall				string
	//	tst_res						result_test_i2c
}

type result_test_modem struct {
	find_dev					bool
	IMEI						string	/* Read IME number*/
	init_modem					bool 	/* Initialize the modem*/
	sim_card					bool	/* Check is Sim connected*/
	rssi						string	/* Check rssi */
	pdr_name					string	/* read cellular provider*/
	ip							string	/* IP*/
	ping_test					bool	/* Ping test*/
	dwnload_test				bool	/* Download packet with checksum*/
	test_start_str				string
	test_stop_str				string
}



func (tst_res result_test_modem) LogResult(param_test param_test_modem)  {
	if tst_res.find_dev == true{
		//Error.Printf("Modem %s not detect \n",param_test.Dev_name)
	} else{
		//Info.Printf("Modem %s detected\n",param_test.Dev_name)
		//Info.Printf("Modem IMEI number %s\n",tst_res.IMEI)
		if tst_res.init_modem{
			//	Info.Println("Modem init - OK")
		}else {
			//	Error.Println("Modem init - Failed")
		}
		if tst_res.sim_card{
			//	Info.Println("Modem Sim Card - OK")
		}else {
			//	Error.Println("Modem Sim Card - Failed")
		}
		if tst_res.ping_test{
			//	Info.Println("Modem ping test - OK")
		}else {
			//	Error.Println("Modem ping test - Failed")
		}
		if tst_res.dwnload_test{
			//	Info.Println("Modem download test - OK")
		}else {
			//	Error.Println("Modem download test - Failed")
		}
		//Info.Printf("Modem RSSI %s\n",tst_res.rssi)
		//Info.Printf("Modem provider %s\n",tst_res.pdr_name)
	}
}
func (tst_res result_test_modem) PrintResult(param_test param_test_modem,attempt int)  {
	str_array,_:=tst_res.ResultToStr(param_test,attempt)
	for _,str:=range str_array{
		fmt.Println(str)
	}

}
func (tst_res result_test_modem) ResultToStr(param_test param_test_modem, attempt int) ([]string,bool){

	var status bool = true
	res_test_str := []string{}
	var fail_string string = c_str_FAIL
	if attempt >0 {
		fail_string = c_str_WARNING
	}
	res_test_str = append(res_test_str,tst_res.test_start_str)
	if tst_res.find_dev != true{
		res_test_str = append(res_test_str,fmt.Sprintf("    Device %s not detect - %s",param_test.Dev_name,c_str_FAIL))
		//fmt.Printf("    Device %s not detect \n",param_test.Dev_name)
	} else {
		res_test_str = append(res_test_str, fmt.Sprintf("    Device %s detected - %s", param_test.Dev_name,c_str_PASS))
		//fmt.Printf("    Device %s detected\n",param_test.Dev_name)

		res_test_str = append(res_test_str, fmt.Sprintf("    Modem IMEI number %s", tst_res.IMEI))

		if tst_res.sim_card {
			res_test_str = append(res_test_str, "    Modem Sim Card - "+c_str_PASS)
		} else {
			res_test_str = append(res_test_str, "    Modem Sim Card - "+fail_string)
			status = false
		}
		res_test_str = append(res_test_str, fmt.Sprintf("    Modem RSSI %s", tst_res.rssi))
		res_test_str = append(res_test_str, fmt.Sprintf("    Modem provider %s", tst_res.pdr_name))

		if param_test.GPRS_test == "true"{
			if tst_res.init_modem {
				res_test_str = append(res_test_str, fmt.Sprintf("    Modem init - %s",c_str_PASS))
			} else {
				res_test_str = append(res_test_str,fmt.Sprintf("    Modem init - %s",fail_string))
				status = false

			}
			if tst_res.ping_test {
				res_test_str = append(res_test_str, fmt.Sprintf("    Modem ping test - %s",c_str_PASS))
			} else {
				res_test_str = append(res_test_str, fmt.Sprintf("    Modem ping test - %s",fail_string))
				status = false
			}
			if tst_res.dwnload_test {
				res_test_str = append(res_test_str,fmt.Sprintf("    Modem download test - %s",c_str_PASS))
			} else {
				res_test_str = append(res_test_str, fmt.Sprintf("    Modem download test - %s",fail_string))
				status = false
			}
		}



	}
	if param_test.GPRS_test == "true"{
		if tst_res.find_dev&&tst_res.dwnload_test&&tst_res.init_modem&&tst_res.ping_test&&tst_res.sim_card{
			res_test_str = append(res_test_str,tst_res.test_stop_str+" - "+c_str_PASS)
			//fmt.Println(tst_res.test_stop_str+" - PASS")
		}else {
			if tst_res.find_dev == false {
				fail_string = c_str_FAIL
			}
			res_test_str = append(res_test_str,tst_res.test_stop_str+" - "+fail_string)
			//fmt.Println(tst_res.test_stop_str+" - FAIL")
			status = false
		}
	}else {
		if tst_res.find_dev&&tst_res.sim_card{
			res_test_str = append(res_test_str,tst_res.test_stop_str+" - "+c_str_PASS)
			//fmt.Println(tst_res.test_stop_str+" - PASS")
		}else {
			if tst_res.find_dev == false {
				fail_string = c_str_FAIL
			}
			res_test_str = append(res_test_str,tst_res.test_stop_str+" - "+fail_string)
			//fmt.Println(tst_res.test_stop_str+" - FAIL")
			status = false
		}
	}

	//res_test_str = append(res_test_str,tst_res.test_stop_str)
	res_test_str = append(res_test_str," ")
	//fmt.Println("")
	return res_test_str,status
}

/*
echo 12 > /sys/class/gpio/export
echo 36 > /sys/class/gpio/export

AT&F
AT+CPIN?
AT+COPS?
at+creg?
echo out > /sys/class/gpio/gpio12/direction
echo out > /sys/class/gpio/gpio36/direction
root@linaro-alip:/home/linaro# echo 0 > /sys/class/gpio/gpio12/value
root@linaro-alip:/home/linaro# echo 1 > /sys/class/gpio/gpio36/value
root@linaro-alip:/home/linaro# echo 1 > /sys/class/gpio/gpio12/value
root@linaro-alip:/home/linaro# echo 0 > /sys/class/gpio/gpio12/value
root@linaro-alip:/home/linaro# echo 1 > /sys/class/gpio/gpio12/value
root@linaro-alip:/home/linaro# echo 1 > /sys/class/gpio/gpio36/value
root@linaro-alip:/home/linaro# echo 0 > /sys/class/gpio/gpio36/value
root@linaro-alip:/home/linaro# echo 0 > /sys/class/gpio/gpio12/value
root@linaro-alip:/home/linaro# echo 1 > /sys/class/gpio/gpio12/value
root@linaro-alip:/home/linaro# echo 1 > /sys/class/gpio/gpio36/value

sudo systemctl stop ModemManager


AT+QCFG="nwscanseq"
AT+QCFG="nwscanmode"
sudo minicom -b 9600 -D /dev/tty96B0
sudo minicom -b 115200 -D /dev/tty96B0
sudo minicom -b 9600 -D /dev/tty96B1
sudo minicom -b 115200 -D /dev/tty96B1

*/


func (device_param_tests param_test_modem)RebootModem_fpga(){
	cmdec := exec.Command("bash","-c","echo 36 > /sys/class/gpio/export")
	cmdec.Run()
	//fmt.Println(cmdec.Run())
	cmdec = exec.Command("bash","-c","echo out > /sys/class/gpio/gpio36/direction")
	cmdec.Run()

	cmdec = exec.Command("bash","-c","echo 12 > /sys/class/gpio/export")
	cmdec.Run()
	//fmt.Println(cmdec.Run())
	cmdec = exec.Command("bash","-c","echo out > /sys/class/gpio/gpio12/direction")
	cmdec.Run()
	fmt.Println("Reboot Modem FPGA")
	cs_r := gpioreg.ByName("GPIO12")
	cs := gpioreg.ByName("GPIO36")

	c := &serial.Config{Name: device_param_tests.Port, Baud: 115200,ReadTimeout:time.Millisecond*100}
	s, err := serial.OpenPort(c)
	if err != nil {
		return
	}

	cs.Out(gpio.Low)
	cs_r.Out(gpio.Low)

	time.Sleep(time.Millisecond*100)
	cs.Out(gpio.High)
	cs_r.Out(gpio.High)
	time.Sleep(time.Millisecond*800)
	cs.Out(gpio.Low)
	cs_r.Out(gpio.Low)

	ret_str := readResponseResetOrPwrDown(s)
	//fmt.Println("Reboot ret = ",ret_str)
	if(ret_str == "RDY"){
		s.Close()
		//fmt.Println("Reboot ret = RDY exit")
		return
	}else{
		time.Sleep(time.Millisecond*2000)
		cs.Out(gpio.High)
		cs_r.Out(gpio.High)
		time.Sleep(time.Millisecond*800)
		cs.Out(gpio.Low)
		cs_r.Out(gpio.Low)
		ret_str = readResponseResetOrPwrDown(s)
		//fmt.Println("Reboot ret = ",ret_str)
		s.Close()
		return
	}

	// Проверяем RDY или POWERED DOWN


}

func (device_param_tests param_test_modem)RebootModem_1(){
	cmdec := exec.Command("bash","-c","echo 36 > /sys/class/gpio/export")
	cmdec.Run()
	//fmt.Println(cmdec.Run())
	cmdec = exec.Command("bash","-c","echo out > /sys/class/gpio/gpio36/direction")
	cmdec.Run()

	cmdec = exec.Command("bash","-c","echo 12 > /sys/class/gpio/export")
	cmdec.Run()
	//fmt.Println(cmdec.Run())
	cmdec = exec.Command("bash","-c","echo out > /sys/class/gpio/gpio12/direction")
	cmdec.Run()
	fmt.Println("Reboot Modem")
	cs_r := gpioreg.ByName("GPIO12")
	cs := gpioreg.ByName("GPIO36")
	cs.Out(gpio.Low)
	cs_r.Out(gpio.Low)
	time.Sleep(time.Millisecond*800)
	cs.Out(gpio.High)


	time.Sleep(time.Millisecond*900)
	cs.Out(gpio.Low)
	time.Sleep(time.Millisecond*700)
	cs.Out(gpio.High)
	time.Sleep(time.Millisecond*900)
	i:=0
	for(i<40){
		fmt.Print(".")
		time.Sleep(time.Millisecond*1000)
		i = i+1
	}

	cs.Out(gpio.Low)
	time.Sleep(time.Millisecond*300)
	cs.Out(gpio.High)
	cs_r.Out(gpio.High)
	i=0
	for(i<25){
		fmt.Print(".")
		time.Sleep(time.Millisecond*1000)
		i = i+1
	}
}

func (device_param_tests param_test_modem)RebootModem(){
	cmdec := exec.Command("bash","-c","echo 36 > /sys/class/gpio/export")
	cmdec.Run()
	//fmt.Println(cmdec.Run())
	cmdec = exec.Command("bash","-c","echo out > /sys/class/gpio/gpio36/direction")
	cmdec.Run()
	fmt.Println("Reboot Modem")
	cs := gpioreg.ByName("GPIO36")
	cs.Out(gpio.Low)
	time.Sleep(time.Millisecond*800)
	cs.Out(gpio.High)


	time.Sleep(time.Millisecond*900)
	cs.Out(gpio.Low)
	time.Sleep(time.Millisecond*700)
	cs.Out(gpio.High)
	time.Sleep(time.Millisecond*900)
	i:=0
	for(i<40){
		fmt.Print(".")
		time.Sleep(time.Millisecond*1000)
		i = i+1
	}

	cs.Out(gpio.Low)
	time.Sleep(time.Millisecond*300)
	cs.Out(gpio.High)
	i=0
	for(i<25){
		fmt.Print(".")
		time.Sleep(time.Millisecond*1000)
		i = i+1
	}
}


func (device_param_tests param_test_modem)AT_SendReq(atcmd string,timeout int, f ...func())(string,error){
	time.Sleep(time.Millisecond*500)
	CMD_REQ := GSM_Cmd{Cmd:fmt.Sprintf("AT%s\r\n",atcmd),CmdResponseOnOk:"OK",TimeoutMs:timeout,DelayMs:50,Skip:0}
	c := &serial.Config{Name: device_param_tests.Port, Baud: 115200,ReadTimeout:time.Millisecond*100}
	s, err := serial.OpenPort(c)
	if err != nil {
		return "", fmt.Errorf("Can't open port")
	}
	respstr := ReqAT(s,CMD_REQ)
	fmt.Println("AT_SendReq respstr",respstr)
	arr_resp := RespToInfo(respstr)
	for _,el:=range arr_resp{
		if(el == "ERROR"){
			s.Close()
			return "",fmt.Errorf("ERROR ",atcmd)
		}
	}
	if(strings.Contains(atcmd,"CMGS")){
		s.Write([]byte{0x1A})
	}
	s.Close()
	//fmt.Println(arr_resp)
	if(len(arr_resp)>0){
		return arr_resp[0],nil
	}
	return "",nil
}
//	Return slice answer, end define as OK or ERROR
func (device_param_tests param_test_modem)AT_ReqAns(atcmd string,timeout int, f ...func())([]string,error){
	time.Sleep(time.Millisecond*500)
	CMD_REQ := GSM_Cmd{Cmd:fmt.Sprintf("AT%s\r\n",atcmd),CmdResponseOnOk:"OK",TimeoutMs:timeout,DelayMs:50,Skip:0}
	c := &serial.Config{Name: device_param_tests.Port, Baud: 115200,ReadTimeout:time.Millisecond*100}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, fmt.Errorf("Can't open port")
	}
	sendCMDToSerial(s,CMD_REQ)
	buff :=ReadResponceLines(s)
	var ret []string
	for _,el:=range buff{
		if(strings.HasSuffix(strings.TrimSpace(el),strings.TrimSpace(CMD_REQ.Cmd))==true){
			continue
		}
		if(strings.HasPrefix(strings.TrimSpace(el),strings.TrimSpace(CMD_REQ.Cmd))==true){
			continue
		}
		if((el == "")||(el==" ")){
			fmt.Println("(el == \"\")||(el==\" \")",el)
			continue
		}
		ret=append(ret,el)
	}
	return ret,nil
}


// Идет первым, если ошибка - перезагружаем
func (device_param_tests param_test_modem)AT_TestModemAndSim()(string,error){
	cmd_nm:= exec.Command("systemctl","stop","ModemManager")
	cmd_nm.Start()
	i:=0
	var err error
	var str_buff string
	for(i<3){
		str_buff,err= device_param_tests.AT_SendReq("+CPIN?",200)
		if(err!=nil){
			i=i+1
		}else{
			break
		}
	}
	if(err!=nil){
		return "", err
	}
	if(strings.Contains(str_buff,"READY")){
		return "Sim Card READY", nil
	}else{
		return "Sim Card Not READY", nil
	}
}
func (device_param_tests param_test_modem)AT_GetProviderName()(string,error){
	i:=0
	var err error
	var str_buff string
	for(i<3){
		str_buff,err= device_param_tests.AT_SendReq("+COPS?",200)
		if(err!=nil){
			i=i+1
		}else{
			break
		}
	}
	if(err != nil){
		return str_buff, err
	}
	str_arr := strings.Split(str_buff,",")
	if(len(str_arr)>2){
		return str_arr[2],nil
	}else {
		return "",fmt.Errorf("Error Operator Register")
	}
	return "",nil
}
func (device_param_tests param_test_modem)AT_GetIMEI()(string,error){
	i:=0
	var err error
	var str_buff string
	for(i<3){
		str_buff,err= device_param_tests.AT_SendReq("+CIMI",300)
		if(err!=nil){
			i=i+1
		}else{
			break
		}
	}
	return str_buff,err
}
func (device_param_tests param_test_modem)AT_GetCSQ()(string,error){
	i:=0
	var err error
	var str_buff string
	for(i<3){
		str_buff,err= device_param_tests.AT_SendReq("+CSQ",500)
		if(err!=nil){
			i=i+1
		}else{
			break
		}
	}
	return str_buff,err
	//return device_param_tests.AT_SendReq("+CSQ",200)
}

//Read and return the message at index idx or None if no message is found at that index.
//
func (device_param_tests param_test_modem)AT_ReadSMS(num  int)([]string,error){

	//device_param_tests.AT_ConfigSMS()
	//AT+CMGR=5
	cmd:=fmt.Sprintf("+CMGR=%d",num)
	//device_param_tests.AT_ReqAns(cmd,500)
	str,_:=device_param_tests.AT_ReqAns(cmd,500)
	var ret []string
	ret = make([]string,7)
	ret[0] = strconv.Itoa(num)

	for i,el:=range str{
		if(strings.HasPrefix(el,"+CMGR:")){
			el=strings.TrimSpace(strings.TrimPrefix(el,"+CMGR:"))
			buff:=strings.Split(el,",")
			if(len(buff)>3){
				//this is text mode(or GSM)
				ret[1] = buff[0]
				ret[2] = buff[1]
				ret[3] = buff[3]
				ret[4] = buff[2]
				ret[5] = buff[2]
				if((i+1)<len(str)){
					ret[6] = str[i+1]
				}
			}
		}
		//fmt.Println("el - ",el)
		//fmt.Println("i - ",i)
	}
	return ret,nil
	//return device_param_tests.AT_SendReq("+CSQ",200)
}

func (device_param_tests param_test_modem)AT_DelSMS(num  int)([]string,error){
	//AT+CMGD=
	cmd:=fmt.Sprintf("+CMGD=%d",num)
	//device_param_tests.AT_ReqAns(cmd,500)
	str,_:=device_param_tests.AT_ReqAns(cmd,500)
	return str,nil
}

func (device_param_tests param_test_modem)AT_ConfigSMS()([]string,error){

	//AT+CMGR=5
	//AT+CMGF=1
	//AT+CSCS="GSM"
	//AT+CSCB=1
	cmd:="+CMGF=1"
	//device_param_tests.AT_ReqAns(cmd,500)
	str,_:=device_param_tests.AT_ReqAns(cmd,500)
	//fmt.Println("11111111111111")
	//fmt.Println(str)
	cmd="+CSCS=\"GSM\""
	str,_=device_param_tests.AT_ReqAns(cmd,500)
	//fmt.Println("11111111111111")
	//fmt.Println(str)
	cmd="+CSCB=1"
	str,_=device_param_tests.AT_ReqAns(cmd,500)
	//fmt.Println("11111111111111")
	//fmt.Println(str)

	return str,nil
	//return device_param_tests.AT_SendReq("+CSQ",200)
}


func (device_param_tests param_test_modem)AT_ATH()(string,error){

	var err error
	var str_buff string

	str_buff,err= device_param_tests.AT_SendReq("H",500)
	return str_buff,err
	//return device_param_tests.AT_SendReq("+CSQ",200)
}
func (device_param_tests param_test_modem)AT_SetAPN()(string,error){
	str_buff := fmt.Sprintf("+QICSGP=1,1,\"%s\"",device_param_tests.APN)
	return device_param_tests.AT_SendReq(str_buff,5000)
}
func (device_param_tests param_test_modem)AT_DisconnectGPRS()(string,error){
	return device_param_tests.AT_SendReq("+QIDEACT=1",10000)
}
func (device_param_tests param_test_modem)AT_ConnectGPRS()(string,error){

	_,reterr:= device_param_tests.AT_DisconnectGPRS()
	if(reterr != nil){
		return "AT_DisconnectGPRS", reterr
	}
	time.Sleep(time.Millisecond*500)
	_,reterr= device_param_tests.AT_SetAPN()
	if(reterr != nil){
		return "AT_SetAPN", reterr
	}
	time.Sleep(time.Millisecond*500)

	//retstr,reterr= device_param_tests.AT_GetIP()
	return device_param_tests.AT_SendReq("+QIACT=1",5000)
}
func (device_param_tests param_test_modem)AT_GetIP()(string,error){
	retstr,reterr := device_param_tests.AT_SendReq("+QIACT?",5000)
	if(reterr != nil){
		return "AT_GetIP", reterr
	}
	arrretstr := strings.Split(retstr,",")
	if(len(arrretstr)<4){
		return "", fmt.Errorf("No IP")
	}else{
		return arrretstr[3], nil
	}
}

func (device_param_tests param_test_modem)AT_CMGF()(string,error){

	var err error
	var str_buff string

	str_buff,err= device_param_tests.AT_SendReq("+CMGF=1",200)
	return str_buff,err
	//return device_param_tests.AT_SendReq("+CSQ",200)
}

func (device_param_tests param_test_modem)AT_ConnectGPRSAndPing()(string,error){
	i:=0
	var err error
	var str_buff string
	for(i<3){
		str_buff,err= device_param_tests.AT_SendReq("+QPING=1,\"google.com\",1,1",4000)
		if(err!=nil){
			fmt.Println("ping test ",i)
			i=i+1
		}else{
			str_buff,err=device_param_tests.AT_SendReq("",300)
			return str_buff,err
		}
	}

	return device_param_tests.AT_SendReq("+QPING=1,\"google.com\"",4000)

}

func (device_param_tests param_test_modem)AT_GetError()(string,error){
/*
	m, err := serial.New(device_param_tests.Port, 115200)
	if err != nil {
		fmt.Println("Device not found")
		fmt.Println(err)
		//tst_res.find_dev = false
	}else{
		defer m.Close()
		var mio io.ReadWriter = m
		a := at.New(mio)

		ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Millisecond)
		info, err := a.Command(ctx, "+QIGETERROR?")
		cancel()
		if err != nil {
			fmt.Println("Modem not answer +QIGETERROR?", err)
			//tst_res.find_dev = false
			//tst_res.IMEI = "Not read IMEI"
		} else {
			for _, l := range info {
				return strings.TrimLeft(l,"+QIGETERROR: "),nil
				fmt.Println("info",l)
			}
		}



	}*/
	return "",nil
}

func (device_param_tests param_test_modem)AT_TestFull()(result_test_modem){
	var tst_res result_test_modem
	tst_res.find_dev = false
	tst_res.init_modem = false
	tst_res.dwnload_test = false


	if device_param_tests.Dev_name != ""{
		//	tst_res.test_start_str = time.Now().Format("01/02/2006 15:04:05")+" "+device_param_tests.Dev_name+" Test Start"
		//Если команда не проходит - выходим

		str_buff,err:=device_param_tests.AT_TestModemAndSim()
		if err != nil {
			fmt.Println("Device not found")
			fmt.Println(err)
			tst_res.find_dev = false
			return tst_res
		}
		tst_res.find_dev = true
		tst_res.sim_card = true

		str_buff,err=device_param_tests.AT_GetIMEI()
		if err != nil {
			fmt.Println("Modem not answer +CIMI", err)
			tst_res.IMEI = "Not read IMEI"
		} else {
			tst_res.IMEI = str_buff
		}

		str_buff,err=device_param_tests.AT_GetCSQ()
		if err != nil {
			fmt.Println("+CSQ",err)
			tst_res.rssi = "No Signal"
		} else {
			tst_res.rssi = str_buff
		}
		str_buff,err=device_param_tests.AT_GetProviderName()
		if err != nil {
			//fmt.Println("Error AT_GetProviderName():",err)
			tst_res.pdr_name = "Not detect"
		} else {
			tst_res.pdr_name = str_buff
		}
		if((device_param_tests.PhoneSMSandCall!="")){
			fmt.Println("send sms")
			device_param_tests.AT_SendSMS("Hello",100)
			fmt.Println("Call")
			device_param_tests.AT_Call(200)
		}

		//		GPRS Test
		if(device_param_tests.GPRS_test == "true"){
			str_buff,err=device_param_tests.AT_ConnectGPRS()
			if(err!=nil){
				//fmt.Println("Error AT_ConnectGPRS(): ",err)
				tst_res.init_modem = false
				tst_res.dwnload_test = false
				return tst_res
			}else{
				tst_res.init_modem = true
				tst_res.dwnload_test = true
			}

			//Получим IP
			if(tst_res.init_modem){
				str_buff,err=device_param_tests.AT_GetIP()
				tst_res.ip = str_buff
			}
			if(tst_res.ip == ""){
				tst_res.init_modem = false
				tst_res.dwnload_test = false
				return tst_res
			}

			str_buff,err=device_param_tests.AT_ConnectGPRSAndPing()

			if(err != nil){
				//fmt.Println("Error AT_ConnectGPRSAndPing():",err)
				tst_res.init_modem = false
				tst_res.dwnload_test = false
				return tst_res
			}
			tst_res.ping_test = true
		//	fmt.Println("AT_ConnectGPRSAndPing() Ok",str_buff)

			device_param_tests.AT_DisconnectGPRS()
		}else{
			tst_res.init_modem = true
			tst_res.dwnload_test = true
			tst_res.ping_test = true
		}


	}

	return tst_res
}

func (device_param_tests param_test_modem)AT_LoopTestFull()(result_test_modem){

	var tst_res result_test_modem
	tst_res.test_start_str = time.Now().Format("01/02/2006 15:04:05")+" "+device_param_tests.Dev_name+" Test Start"
	//device_param_tests.RebootModem_fpga()
	//if()
	//device_param_tests.RebootModem_fpga()
	//time.Sleep(time.Second*5)

	//device_param_tests.RebootModem()
	tst_res = device_param_tests.AT_TestFull()

	if(tst_res.init_modem != true){
		if(device_param_tests.Port=="/dev/tty96B0"){
			device_param_tests.RebootModem_fpga()
		}else{
			device_param_tests.RebootModem()
		}

		tst_res = device_param_tests.AT_TestFull()
	}
	tst_res.test_stop_str = time.Now().Format("01/02/2006 15:04:05")+" "+device_param_tests.Dev_name+" Test Finish"
	return tst_res
}


func TestModemDevice(modem_dev param_test_modem) result_test_modem{

	cmd_nm:= exec.Command("systemctl","stop","ModemManager")
	cmd_nm.Start()



	//tst_res:=TestFullATModemDevice(modem_dev)
	fmt.Println("Please wait")
	var tst_res result_test_modem


	//tst_res=modem_dev.AT_LoopTestFull()

	if(modem_dev.ATcommand == "true"){
		//Тестируем только при помощи АТ команд
		tst_res=modem_dev.AT_LoopTestFull()
	} else{

		//tst_res=TestATModemDevice(modem_dev)

		if  tst_res.find_dev && modem_dev.GPRS_test == "true"{
			//ReconfigWvdial(modem_dev)
			cmd := exec.Command("wvdial")

			//stdout, err := cmd.StdoutPipe()
			//stdout, err := cmd.StdoutPipe()

			if err := cmd.Start(); err != nil {
				fmt.Println(err)
				tst_res.init_modem=false
				tst_res.ping_test = false
			}else{
				fmt.Println("Please Wait")
				fmt.Printf("Connect to GPRS")
				for i:=0; i<40;i++{
					fmt.Printf(".")
					time.Sleep(time.Second)
				}
				iface,err:=net.Interfaces()
				if err!=nil{
					tst_res.init_modem=false
					tst_res.ping_test = false
				}else {
					for _,i := range iface{
						if strings.Contains(i.Name,"ppp"){
							fmt.Printf("\n")
							fmt.Println("Connected Ok")
							tst_res.init_modem = true
							tst_res.ping_test = true
							tst_res.dwnload_test = true
						}
					}
				}
				fmt.Println("Disconnect, please wait")
				cmd.Process.Signal(os.Interrupt)
				time.Sleep(10*time.Second)

				if err := cmd.Process.Kill(); err != nil {

					fmt.Println("failed to kill process: ", err)
				}
			}
		}
	}



	/*
	// Запускаем wvdial и ждем секунд 10
	//для начала перезаписываем файл wvdial.conf

	if  tst_res.find_dev && modem_dev.GPRS_test == "true"{
		ReconfigWvdial(modem_dev)

		cmd := exec.Command("wvdial")

		//stdout, err := cmd.StdoutPipe()
		//stdout, err := cmd.StdoutPipe()

		if err := cmd.Start(); err != nil {
			fmt.Println(err)
			tst_res.init_modem=false
			tst_res.ping_test = false
		}else{
			fmt.Println("Please Wait")
			fmt.Printf("Connect to GPRS")
			for i:=0; i<40;i++{
				fmt.Printf(".")
				time.Sleep(time.Second)
			}
			iface,err:=net.Interfaces()
			if err!=nil{
				tst_res.init_modem=false
				tst_res.ping_test = false
			}else {
				for _,i := range iface{
					if strings.Contains(i.Name,"ppp"){
						fmt.Printf("\n")
						fmt.Println("Connected Ok")
						tst_res.init_modem = true
						tst_res.ping_test = true
						tst_res.dwnload_test = true
					}
				}
			}
			fmt.Println("Disconnect, please wait")
			cmd.Process.Signal(os.Interrupt)
			time.Sleep(10*time.Second)

			if err := cmd.Process.Kill(); err != nil {

				fmt.Println("failed to kill process: ", err)
			}

			//fmt.Println("cmd := exec.Command Kill")
		}

	}
	*/

	tst_res.test_stop_str = time.Now().Format("01/02/2006 15:04:05")+" "+modem_dev.Dev_name+" Test Finish"
	tst_res.PrintResult(modem_dev,0)
	tst_res.LogResult(modem_dev)
	return tst_res
}

/*
	Для команды ping сразу выдает ok
*/
func readResponseResetOrPwrDown(s *serial.Port)string{

	f_end:= false
	f_OkError:= false
	count:=100
	buff_byte := make([]byte, 1)
	str_buff := make([]byte, 0)
	//fmt.Println("READ START --------------------------------")
	for !f_end{
		_,err:=s.Read(buff_byte)
		if(err!=nil){
			count -=1
			if(count<0){
				return "ERROR"
			}
		}else{
			//fmt.Print(string(buff_byte[0]))
			str_buff = append(str_buff, buff_byte[0])
			//str_buff = bytes.Trim(str_buff,cmd.Cmd)
			if(strings.Contains(string(str_buff),"RDY")){
				f_OkError = true
				return "RDY"
			}
			if((strings.Contains(string(str_buff),"POWERED DOWN"))){
				f_OkError = true
				return "POWERED DOWN"
			}
			if((buff_byte[0] == '\n')&&(f_OkError)){
				f_end = true
			}
		}
	}
	//fmt.Print("\n")
	//fmt.Println("READ STOP ---------------------------------")
	//str_buff = strings.TrimLeft(string(str_buff),cmd.Cmd)
	return string(str_buff)
}

/*
	Для команды ping сразу выдает ok
*/
func readResponse(s *serial.Port,cmd GSM_Cmd)string{

	f_end:= false
	f_OkError:= false
	count:=cmd.TimeoutMs/100
	buff_byte := make([]byte, 1)
	str_buff := make([]byte, 0)
	//fmt.Println("READ START --------------------------------")
	for !f_end{
		_,err:=s.Read(buff_byte)
		if(err!=nil){
			count -=1
			if(count<0){
				f_end = true
			}
		}else{
			//fmt.Print(string(buff_byte[0]))
			str_buff = append(str_buff, buff_byte[0])
			//str_buff = bytes.Trim(str_buff,cmd.Cmd)
			if(strings.Contains(string(str_buff),"OK")||(strings.Contains(string(str_buff),"ERROR"))||
				(strings.Contains(string(str_buff),"RING"))||(strings.Contains(string(str_buff),"BUSY"))){
				f_OkError = true
			}
			if((buff_byte[0] == '\n')&&(f_OkError)){
				f_end = true
			}
		}
	}
	return strings.TrimLeft(string(str_buff),cmd.Cmd+": ")
}




func ReqAT(s *serial.Port,cmd GSM_Cmd)string{
	sendCMDToSerial(s,cmd)
	str := readResponse(s,cmd)
	return str
}

func RespToInfo(str string) []string  {
	var retstr []string

	arr_resp:=strings.Split(str,"\r\n")
	//fmt.Println("RespToInfo arr_resp",arr_resp)
	//fmt.Println("RespToInfo arr_resp len", len(arr_resp))
	for _,el:=range arr_resp {
		if((strings.TrimSpace(el) !="")){//&&(strings.TrimSpace(el) !="OK")){
			retstr = append(retstr,strings.TrimSpace(el))
		}
	}
	fmt.Println("RespToInfo retstr",retstr)
	return retstr
}

func ReadResponceLines(s *serial.Port) ([]string) {

	f_end:= false
	//f_OkError:= false
	count:=2		//20msec
	buff_byte := make([]byte, 1)
	str_buff := make([]byte, 0)
	var temp_str string
	var retbuf []string
	//fmt.Println("READ START --------------------------------")
	for !f_end{
		_,err:=s.Read(buff_byte)
		if(err!=nil){
			count -=1
			if(count<0){
				f_end = true	// if not responce
			}
		}else{
			//fmt.Print(string(buff_byte[0]))
			str_buff = append(str_buff, buff_byte[0])
			if(bytes.HasSuffix(str_buff,[]byte("\r\n"))){
				temp_str = string(bytes.TrimSuffix(str_buff,[]byte("\r\n")))
				if(temp_str!=""){
					retbuf = append(retbuf,strings.TrimSpace(temp_str))
					if(strings.HasPrefix(temp_str,"OK")){
						f_end = true
					}
					if(strings.HasPrefix(temp_str,"ERROR")){
						f_end = true
					}
				}
				str_buff = make([]byte, 0)
			}
		}
	}
	return retbuf
}



/*
	Отправка запросса
*/

func sendCMDToSerial(s *serial.Port,cmd GSM_Cmd){
	s.Write([]byte(cmd.Cmd))
}


func (device_param_tests param_test_modem)AT_SendSMS(atcmd string,timeout int, f ...func())(string,error){


	var err error
	var str_buff string

	device_param_tests.AT_CMGF()

	str_buff=fmt.Sprintf("+CMGS=\"%s\"\r\n Test\r\n",device_param_tests.PhoneSMSandCall)
	str_buff = str_buff+string(0x1a)
	fmt.Println("str_buff = ",str_buff)
	str_buff,err= device_param_tests.AT_SendReq(str_buff,500)

	device_param_tests.AT_ATH()

	//fmt.Println(str_buff)
	return str_buff,err


	/*
	c := &serial.Config{Name: device_param_tests.Port, Baud: 115200,ReadTimeout:time.Millisecond*500}
	s, err := serial.OpenPort(c)
	if err != nil {
		//log.Fatal("Error from open port ",device_param_tests.Port)
		//log.Fatal(err)
		return "", fmt.Errorf("Can't open port")
	}
	var mio io.ReadWriter = s
	g := gsm.New(mio)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()
	if err = g.Init(ctx); err != nil {
		//log.Fatal(err)
		fmt.Println("Error g.Init = ",err)
	}

	mr, err := g.SendSMS(ctx, device_param_tests.PhoneSMSandCall, "Test LTE")
	// !!! check CPIN?? on failure to determine root cause??  If ERROR 302
	log.Printf("%v %v\n", mr, err)
	return mr,err
	*/
}

func (device_param_tests param_test_modem)AT_Call(timeout int)(string,error){

	var err error
	var str_buff string
	str_buff,err= device_param_tests.AT_SendReq("D"+device_param_tests.PhoneSMSandCall+";",200)


	if(err != nil){
		return str_buff, err
	}
	//если все нормально, ждем ответ

	time.Sleep(time.Millisecond*500)
	c := &serial.Config{Name: device_param_tests.Port, Baud: 115200,ReadTimeout:time.Millisecond*100}
	s, err := serial.OpenPort(c)
	if err != nil {
		//log.Fatal("Error from open port ",device_param_tests.Port)
		//log.Fatal(err)
		return "", fmt.Errorf("Can't open port")
	}

	CMD_REQ := GSM_Cmd{Cmd:"",CmdResponseOnOk:"OK",TimeoutMs:20000,DelayMs:50,Skip:0}
	respstr := readResponse(s,CMD_REQ)

	arr_resp := RespToInfo(respstr)
	for _,el:=range arr_resp{
		if(el == "ERROR"){
			s.Close()
			return "",fmt.Errorf("ERROR Call")
		}
	}
	s.Close()
	fmt.Println(arr_resp)

	device_param_tests.AT_ATH()


	return "",nil


	fmt.Println(str_buff)

	return "",nil
}


