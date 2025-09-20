package bsc

import (
	"os/exec"
	"syscall"
)

type ExecutionArguments struct {
	python       string
	pythonApp    string
	DBName       string
	ExcelFile    string
	ExcelSheet   string //TODO should be sheets of multiple strings
	CategoryName string
	CategoryDesc string
}

const python = "python"
const python_app = "/opt/BSC/src/main.py"

func (ea ExecutionArguments) BSCExecution() (int, string) {
	ea.python = python
	ea.pythonApp = python_app

	cmd := exec.Command(ea.python, ea.pythonApp, "--db_name="+ea.DBName, "--file="+ea.ExcelFile, "--sheet="+ea.ExcelSheet, "--c_name="+ea.CategoryName, "--c_desc="+ea.CategoryDesc)

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
