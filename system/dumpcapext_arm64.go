// Copyright 2019-2020 Graham Clark. All rights reserved.  Use of this source
// code is governed by the MIT license that can be found in the LICENSE
// file.

package system

import (
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

//======================================================================

// DumpcapExt will run dumpcap first, but if it fails, run tshark. Intended as
// a special case to allow termshark -i <iface> to use dumpcap if possible,
// but if it fails (e.g. iface==randpkt), fall back to tshark. dumpcap is more
// efficient than tshark at just capturing, and will drop fewer packets, but
// tshark supports extcap interfaces.
func DumpcapExt(dumpcapBin string, tsharkBin string, args ...string) error {
	var err error

	dumpcapCmd := exec.Command(dumpcapBin, args...)
	log.Infof("Starting dumpcap command %v", dumpcapCmd)
	dumpcapCmd.Stdin = os.Stdin
	dumpcapCmd.Stdout = os.Stdout
	dumpcapCmd.Stderr = os.Stderr
	if dumpcapCmd.Run() != nil {
		var tshark string
		tshark, err = exec.LookPath(tsharkBin)
		if err == nil {
			log.Infof("Retrying with dumpcap command %v", append([]string{tshark}, args...))
			err = syscall.Exec(tshark, append([]string{tshark}, args...), os.Environ())
		}
	}

	return err
}
