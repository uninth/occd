package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

func main() {
	/* The default logger in Go writes to stderr (2). redirect to file
	** openvpn requires exit status to be 0 for scripts, otherwise the user
	** is not allowed to login */
	var fErr *os.File
	var err error

	fErr, err = os.OpenFile("/tmp/openvpn_connect_err.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	syscall.Dup2(int(fErr.Fd()), 1) /* -- stdout */
	syscall.Dup2(int(fErr.Fd()), 2) /* -- stderr */

	var DATADIR string = "/openvpn_statistic/" // running chrooted

	var csv_file = DATADIR + os.Getenv("common_name") + ".csv"

	/*
		name, ok := os.LookupEnv("NAME")
		returns the value of the environment variable in its first parameter if
		set, otherwise the second parameter is false. Allows to distinguish
		unset from empty value. but not needed here I think ...
	*/

	// read selected parameters from environment
	var time_duration = os.Getenv("time_duration")
	var bytes_received = os.Getenv("bytes_received")
	var bytes_sent = os.Getenv("bytes_sent")
	var ifconfig_ipv6_local = os.Getenv("ifconfig_ipv6_local")
	var ifconfig_ipv6_netbits = os.Getenv("ifconfig_ipv6_netbits")
	var ifconfig_ipv6_remote = os.Getenv("ifconfig_ipv6_remote")
	var untrusted_ip = os.Getenv("untrusted_ip")
	var untrusted_ip6 = os.Getenv("untrusted_ip6")
	var untrusted_port = os.Getenv("untrusted_port")
	var ifconfig_pool_remote_ip = os.Getenv("ifconfig_pool_remote_ip")

	var currentTime = time.Now()
	var time_ascii = currentTime.Format("Mon Jan 2 15:04:05 MST 2006")
	var time_unix = strconv.FormatInt(currentTime.Unix(), 10)

	if err = ensureDir(DATADIR); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(0)
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

	if _, err = os.Stat(csv_file); err == nil {
		// file exists just print body
		f, err = os.OpenFile(csv_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("open csv_file failed: " + err.Error())
			os.Exit(0)
		}
		if _, err := f.Write([]byte(body)); err != nil {
			f.Close() // ignore error; Write error takes precedence
			fmt.Println("write csv_file failed: " + err.Error())
			os.Exit(0)
		}
	} else {
		// file doesn't exist create header and body
		f, err = os.OpenFile(csv_file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("create csv file failed: " + err.Error())
			os.Exit(0)
		}
		if _, err = f.Write([]byte(head)); err != nil {
			fmt.Println("write head to csv_file failed: " + err.Error())
			os.Exit(0)
			f.Close()
		}
		if _, err = f.Write([]byte(body)); err != nil {
			fmt.Println("write body to csv_file failed: " + err.Error())
			os.Exit(0)
			f.Close()
		}
	}

	if err := f.Close(); err != nil {
		fmt.Println("close csv_file failed: " + err.Error())
		os.Exit(0)
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
	/*
		{
			f := filepath.Base(os.Args[0])
			_, err := os.Lstat("bin/connect")
			if err != nil {
				err = os.Symlink(f, "bin/connect")
				if err != nil {
					log.Fatal(err)
				}

			}

			_, err = os.Lstat("bin/disconnect")
			if err != nil {
				err = os.Symlink(f, "bin/disconnect")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	*/

}

/*
 *  Modified BSD License
 *  ====================
 *
 *  Copyright © 2021, Niels Thomas Haugård, www.deic.dk, wwww.i2.dk.dk
 *  All rights reserved.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  1. Redistributions of source code must retain the above copyright
 *     notice, this list of conditions and the following disclaimer.
 *  2. Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *  3. Neither the name of the organisation  www.deic.dk, wwww.i2.dk.dk nor the
 *     names of its contributors may be used to endorse or promote products
 *     derived from this software without specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS” AND
 *  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 *  WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 *  DISCLAIMED. IN NO EVENT SHALL NIELS THOMAS HAUGÅRD BE LIABLE FOR ANY
 *  DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 *  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 *  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 *  ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 *  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 *  SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */
