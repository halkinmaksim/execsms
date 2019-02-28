# execsms
exec commands from SMS

To compile the file under the Raspberry system, set the environment variables.

    set GOARCH = arm
    set GOOS = linux
and run
 
    go build -o execsms
    
To compile the file under the DragonBoard system, set the environment variables.

    set GOARCH = arm64
    set GOOS = linux
and run 

    go build -o execsms


copy findeth.sh to /home/rxhf 

copy execsms and programsettings.json to /home/rxhf

After then, use the command

    sudo chmod +x /home/rxhf/execsms
    
Add program to crontab - for this, use the command
    
    sudo crontab -e
and add task
    
    */1 * * * * /home/rxhf/execsms >> /tmp/sms.log 2>&1
    */1 * * * * /home/rxhf/execsms
    */1 * * * * /home/rxhf/findeth.sh



The application supports the following commands:

RESET LTE - reboot lte from cmd 

This analogue command

    systemctl restart lte
REBOOT GATEWAY  reboot system

This analogue command

    reboot
SET SERVER: address port - set server address and port

To view the event log, run the command

    tail -f /tmp/smsservice.log
 
 To view server settings, run the command
 
    sudo nano /opt/risinghf/pktfwd/local_conf.json


help for cmd on lte

    sudo journalctl -f -n 200 -u lte
    
    sudo journalctl -u pktfwd -f -n 250
    
    netstat -rn
    sudo route del default
    sudo route add default gw 10.64.64.64 dev ppp0
    
     sudo systemctl restart pktfwd
     sudo systemctl stop pktfwd

    sudo ./util_lbt_test -f 868.5 -r -80 -s 5000
     sudo cat /etc/resolv.conf
     
     #!/bin/bash
     route_found=$(/sbin/route -n | /bin/grep -c ^0.0.0.0)
     ppp_on=$(/sbin/ifconfig | /bin/grep -c ppp0)
     if [ $route_found -eq 0 ] && [ $ppp_on -eq 1 ]
       then /sbin/route add default ppp0
     fi