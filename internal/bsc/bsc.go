package bsc

import (
	"os/exec"
	"syscall"
	"strings"
)

type ExecutionArguments struct {
	DBName       string
	ExcelFile    string
	ExcelSheets  []string
	CategoryName string
	CategoryDesc string
}

const python = "python"
const python_app = "/opt/BSC/src/main.py"

func (ea ExecutionArguments) BSCExecution() (int, string) {
	args := ea.compileArgs()
	parts := strings.Fields(args)
	cmd := exec.Command(python, parts...)

	exit_code := 0 //TODO analyze if exit code output is really needed and just doing regular error output on failure is better
	output, err := cmd.CombinedOutput()
	if err != nil {
		exit_error, ok := err.(*exec.ExitError)
		if ok {
			status, ok := exit_error.Sys().(syscall.WaitStatus)
			if ok {
				exit_code = status.ExitStatus()
			}
		}
	}

	return exit_code, string(output)
}

func (ea ExecutionArguments) compileArgs() string {
	result_str := python_app + " insert --db_name=" + ea.DBName + " --file=" + ea.ExcelFile + " --c_name=" + ea.CategoryName + " --c_desc=" + ea.CategoryDesc + " --out=json"
	for _, sheet := range ea.ExcelSheets {
		result_str = result_str + " --sheet=" + sheet
	}
	return result_str
}

func bscPythonTest(command string) (int, string) {
	cmd := exec.Command("python", command)
	cmd.Dir = "/opt/BSC/src/"
	exit_code := 0 //TODO analyze if exit code output is really needed and just doing regular error output on failure is better
	output, err := cmd.CombinedOutput()
	if err != nil {
		exit_error, ok := err.(*exec.ExitError)
		if ok {
			status, ok := exit_error.Sys().(syscall.WaitStatus)
			if ok {
				exit_code = status.ExitStatus()
			}
		}
	}
	return exit_code, string(output)
}
