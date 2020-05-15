#!/usr/bin/env bash
#
# Assuming you have multible openvpn environments, and your environments are
# devined in /usr/local/etc/openvpn while the chroot is in /var/openvpn, you
# could use this script to ensure the logging is setup correctly
#

function contains() {
	word="$1"
	list="$2"
    if [[ ${list} =~ (^|[[:space:]])${word}($|[[:space:]]) ]]; then 
		return 0 
	else
		return 1
	fi
}

usage="$0 environment"

case $# in
	0)	echo $usage; exit
	;;
	*)	contains "$1" "$( ls /usr/local/etc/openvpn )"
		case $? in
			0)	:
			;;
			*)	echo unknown environment $1
				exit
			;;
		esac
	;;
esac

cd /var/openvpn

cd $1 || {
	echo no such environment or dir: $1
}

test -d bin || mkdir bin
(
	cd bin
	test -x occd		|| cp /var/openvpn/occd .
	test -x connect		|| ln occd connect
	test -x disconnect	|| ln occd disconnect
)

test -d openvpn_statistic || {
	mkdir openvpn_statistic
}

# allways do
chown -R _openvpn:_openvpn openvpn_statistic

test -d tmp || {
	mkdir tmp
}

chmod 1777 tmp
chown root:wheel tmp

echo bin, openvpn_statistic, tmp ok in $1

OPENPVN_ENV=$( basename $( pwd ) )

RCSTRING="--script-security 2 --client-connect /bin/connect --client-disconnect /bin/disconnect"

FOUND=$( grep $OPENPVN_ENV /etc/rc.conf.local|egrep '/bin/connect.* /bin/disconnect'|wc -l |tr -d ' ' )

case $FOUND in
	1)	echo "ok: client connect strings found in /etc/rc.conf.local for openvpn_${OPENPVN_ENV}"
	;;
	*)	echo "apply the string $RCSTRING to /etc/rc.conf.local for $OPENPVN_ENV"
	;;
esac

