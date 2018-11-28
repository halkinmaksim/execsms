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
    
The application supports the following commands:

RESET LTE - reboot lte from cmd 

    systemctl restart lte
REBOOT GATEWAY  reboot system
    
    reboot
SET SERVER: address port - set server address and port

    sudo chmod +x execsms
    sudo crontab -e
    */1 * * * * /home/rxhf/execsms >> /tmp/sms.log 2>&1
    */1 * * * * /home/rxhf/execsms


