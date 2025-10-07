package bsc

import (
	"os/exec"
	"syscall"
	"strings"
)

const python = "python"
const python_app = "/opt/BSC/src/main.py"


type ExecutionArguments struct {
	Command            string
	DBName             string
	ExcelFile          string
	ExcelSheets        []string
	CategoryName       string
	CategoryDesc       string
	ReportName         string
	TournamentIDFilter string
}

func (ean ExecutionArguments) BSCExecution() (int, string) {
	args := ean.compileArgsnew()
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

func (ean ExecutionArguments) compileArgsnew() string {
	result_str := python_app + " " + ean.Command + " --db_name=" + ean.DBName + " --r_name=" + ean.ReportName + " --file=" + ean.ExcelFile + " --c_name=" + ean.CategoryName + " --c_desc=" + ean.CategoryDesc + " --out=json"

	if ean.TournamentIDFilter != "" {
		result_str = result_str + " --r_tidf=" + ean.TournamentIDFilter
	}
	if len(ean.ExcelSheets) != 0 {
		for _, sheet := range ean.ExcelSheets {
			result_str = result_str + " --sheet=" + sheet
		}
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
