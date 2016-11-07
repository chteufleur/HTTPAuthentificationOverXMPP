#!/bin/sh

#######################
name="httpAuth"
dirPath="/usr/local/sbin"
#######################
sh="/bin/bash"
pidPath="/var/run/http-auth/http-auth.pid"
#######################

do_start() {
	nohup $dirPath/$name &
	/bin/echo $! > $pidPath
}

do_stop() {
	if [ -e $pidPath ]
	then
		pid=$(/bin/cat $pidPath)
		kill $pid
		rm $pidPath
	fi
}

get_status() {
	if [ -e $pidPath ] ; then
		pid=$(/bin/cat $pidPath) ;
		/bin/echo "$name is running with PID $pid" ;
	else
		/bin/echo "$name is stopped" ;
	fi
}


case $1 in
	start)
		do_start
		;;
	stop)
		do_stop
		;;
	restart)
		do_stop
		do_start
		;;
	status)
		get_status
		;;
	*)
		/bin/echo "Usage: $0 {start|stop|restart|status}"
		exit 2
		;;
esac
