package main

import (
	"bytes"
	"time"
  "fmt"
	"io"
	"log"
	"os"
	"os/exec"
	//"runtime"
	"sync"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
  "flag"
)

type Cfg struct {
  StartCMD  string `yaml:"start"`
  EndCMD  string `yaml:"end"`
  Interval string `yaml:"interval"`
  Shell string `yaml:"shell"`
}


func main() {

  cfgPath := flag.String("c", "./cfg.yml", "configuration file path")
  flag.Parse()

  cfg := getYmlFile(*cfgPath)
  log.Println("start by:", cfg.StartCMD)
  log.Println("end by:", cfg.EndCMD)
  log.Println("interval:", cfg.Interval)

  interval1, err:=  time.ParseDuration(cfg.Interval)
  if err != nil {
      log.Fatalf("Wrong format of cfg.yml's interval", err)
  }
	//timer1 := time.NewTimer(interval1)

  startCMD := cfg.StartCMD
  killCMD := cfg.EndCMD

	log.Println("[START]")
  go RunCMD(cfg.Shell, startCMD)

  for {
    time.Sleep(interval1)
    log.Println("[RESTART]")

    RunCMD(cfg.Shell, killCMD)
    time.Sleep(20* time.Millisecond)
    go RunCMD(cfg.Shell, startCMD) 
  }
	
	fmt.Println("END: ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("ok")
}

func getYmlFile(fname string) Cfg {
  filename, _ := filepath.Abs(fname)
  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
      panic(err)
  }
  var cfg Cfg
  err = yaml.Unmarshal(yamlFile, &cfg)
  return cfg
}

func RunCMD(shell string, bashcmd string) {
	// copy from https://github.com/kjk/go-cookbook/blob/master/advanced-exec/03-live-progress-and-capture-v3.go
	cmd := exec.Command(shell, "-c", bashcmd)
	// On windows you can just use PowerShell so that OS specification is unnecessary here
	/*if runtime.GOOS == "windows" {
		cmd = exec.Command(shell, "/c", bashcmd)
	}*/

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	log.Println("[RUN CMD]", bashcmd)
	if err != nil {
		time.Sleep(1000* time.Millisecond)
		go RunCMD(shell, bashcmd)
		log.Printf("cmd.Start() failed with '%s', try to restart\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		time.Sleep(1000* time.Millisecond)
		go RunCMD(shell, bashcmd)
		log.Printf("cmd.Start() failed with '%s', try to restart\n", err)
		return
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	//outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	//fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
