package externalcmds

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yulibaozi/kubectl-switch/server"
	"github.com/yulibaozi/kubectl-switch/server/fileutil"
)

func init() {
	server.RegisterExternalCmd((*Kubectl)(nil))
}

// Kubectl kubectl
type Kubectl struct{}

// Exec kubectl cmd exec
func (ctl *Kubectl) Exec(cmd *server.Command) error {
	kubePath := server.GetKubeConfigPath()
	exist, dir := fileutil.PathStatus(kubePath)
	if !exist || (exist && dir) {
		err := fileutil.Touch(kubePath)
		if err != nil {
			return fmt.Errorf("touch path failed. err:%v, path:%s", err, kubePath)
		}
	}
	clusterPath := server.GetConfigDir(cmd.SubCmd)
	if !server.CheckConfig(clusterPath) {
		return nil
	}
	configPath, err := server.GetConfigNameByClusterName(cmd.SubCmd)
	if err != nil {
		return err
	}
	equal, err := fileutil.FileMd5Equal(configPath, kubePath)
	if err != nil {
		return err
	}
	if equal {
		empty, err := fileutil.FileEmpty(configPath)
		if err != nil {
			return err
		}
		if empty {
			return fmt.Errorf("%s cluster config file(path:%s) is not allowed to be empty. ", cmd.SubCmd, configPath)
		}
	}
	if !equal {
		if err := fileutil.Copy(configPath, kubePath); err != nil {
			return fmt.Errorf("cp config file failed. err:%v, src:%s, des:%s", err, configPath, kubePath)
		}
		empty, err := fileutil.FileEmpty(kubePath)
		if err != nil {
			return err
		}
		if empty {
			return fmt.Errorf("config file is not allowed to be empty. path:%s", kubePath)
		}
	}
	ctl.execKubectl(cmd)
	return nil
}

func (ctl *Kubectl) execKubectl(cmd *server.Command) {
	if len(cmd.Flags) > 0 {
		cmd.Args = append(cmd.Args, cmd.Flags...)
	}
	var execCmd = exec.Command("kubectl", cmd.Args...)
	out, err := execCmd.CombinedOutput()
	failed := false
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(os.Stdout, text)
		if !strings.Contains(text, "WARNING") {
			failed = true
		}
	}

	if err != nil {
		fmt.Print("exec kubectl err: ", err)
		failed = true
	}
	if failed {
		os.Exit(1)
	}
}