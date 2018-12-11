#!/bin/sh

# -q quiet
# -c nb of pings to perform

ping -q -c5 google.com > /dev/null

if [ $? -eq 0 ]
then
	echo "ok"
else
	/sbin/route del default
	/sbin/route add default ppp0
fi
	echo "fi"
