package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	//	"strings"
)

func main() {
	var DATADIR string = "/openvpn_statistic/" // running chrooted

	var csv_file = DATADIR + os.Getenv("common_name") + ".csv"

	/*
		name, ok := os.LookupEnv("NAME")
		returns the value of the environment variable in its first parameter if
		set, otherwise the second parameter is false. Allows to distinguish
		unset from empty value. but not needed here I think ...
	*/

	// read selected parameters from environment
	var time_ascii = os.Getenv("time_ascii")
	var time_duration = os.Getenv("time_duration")
	var time_unix = os.Getenv("time_unix")
	var bytes_received = os.Getenv("bytes_received")
	var bytes_sent = os.Getenv("bytes_sent")
	var ifconfig_ipv6_local = os.Getenv("ifconfig_ipv6_local")
	var ifconfig_ipv6_netbits = os.Getenv("ifconfig_ipv6_netbits")
	var ifconfig_ipv6_remote = os.Getenv("ifconfig_ipv6_remote")
	var untrusted_ip = os.Getenv("untrusted_ip")
	var untrusted_ip6 = os.Getenv("untrusted_ip6")
	var untrusted_port = os.Getenv("untrusted_port")
	var ifconfig_pool_remote_ip = os.Getenv("ifconfig_pool_remote_ip")

	if err := ensureDir(DATADIR); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(0)
		/* openvpn requires exit status to be 0 for scripts, otherwise the user
		** is not allowed to login
		 */
	}

	/* openvpn args on OpenBSD doesn't honor something with args
	** at least on 6.5, so read 'state' from $0
	 */
	var status string
	switch os := filepath.Base(os.Args[0]); os {
	case "connect":
		status = "connect"
	case "disconnect":
		status = "disconnect"
	default:
		status = "connect_or_disconnect"
	}

	var head = "status;time_ascii;time_duration;time_unix;bytes_received;bytes_sent;ifconfig_ipv6_local;ifconfig_ipv6_netbits;ifconfig_ipv6_remote;untrusted_ip;untrusted_ip6;untrusted_port;ifconfig_pool_remote_ip\n"
	var body = status + ";" + time_ascii + ";" + time_duration + ";" + time_unix + ";" + bytes_received + ";" + bytes_sent + ";" + ifconfig_ipv6_local + ";" + ifconfig_ipv6_netbits + ";" + ifconfig_ipv6_remote + ";" + untrusted_ip + ";" + untrusted_ip6 + ";" + untrusted_port + ";" + ifconfig_pool_remote_ip + "\n"

	var f *os.File
	var err error

	if _, err = os.Stat(csv_file); err == nil {
		// file exists just print body
		f, err = os.OpenFile(csv_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(body)); err != nil {
			f.Close() // ignore error; Write error takes precedence
			log.Fatal(err)
		}
	} else {
		// file doesn't exist create header and body
		f, err = os.OpenFile(csv_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = f.Write([]byte(head)); err != nil {
			f.Close()
			log.Fatal(err)
		}
		if _, err = f.Write([]byte(body)); err != nil {
			f.Close()
			log.Fatal(err)
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	/*
		** not tested but maybe this could be used for looping through the env
		** simplifying the code (openvpn runs in a chroot)
			for _, s := range os.Environ() {
				kv := strings.SplitN(s, "=", 2) // unpacks "key=value"
				fmt.Printf("key:%q value:%q\n", kv[0], kv[1])
			}
	*/
	os.Exit(0)
}

func ensureDir(dirName string) error {

	err := os.MkdirAll(dirName, 0755)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

/*
** Modified BSD License
** ====================
**
** Copyright © 2020, Niels Thomas Haugård
** See ../LICENSE
 */
