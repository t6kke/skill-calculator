package bsc

import (
	"log"
	"os/exec"
	"strings"
	"syscall"
)

const python = "python"
const python_app = "/opt/BSC/src/main.py"

type ExecutionArguments struct {
	Command            string
	DBName             string
	TournamentURLs     []string
	ExcelFile          string
	ExcelSheets        []string
	CategoryName       string
	CategoryDesc       string
	ReportName         string
	TournamentIDFilter string
	ListContent        bool
}

func (ea ExecutionArguments) BSCExecution() (int, string) {
	args := ea.compileArgs()
	parts := strings.Fields(args)
	cmd := exec.Command(python, parts...) //gosec:disable
	log.Printf("BSC execution arguments: %s", args) //gosec:disable

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

	log.Printf("BSC execution exit status: %d", exit_code) //gosec:disable
	log.Printf("BSC execution output: %s", string(output)) //gosec:disable

	return exit_code, string(output)
}

func (ea ExecutionArguments) compileArgs() string {
	//TODO not a good solution, this if/else is done to avoid creating empty catetgory values
	result_str := python_app
	if ea.ListContent {
		result_str = result_str + " " + ea.Command + " --db_name=" + ea.DBName + " --list --out=json"
	} else {
		result_str = result_str + " " + ea.Command + " --db_name=" + ea.DBName + " --r_name=" + ea.ReportName + " --file=" + ea.ExcelFile + " --c_name=" + ea.CategoryName + " --c_desc=" + ea.CategoryDesc + " --out=json"

		if ea.TournamentIDFilter != "" {
			result_str = result_str + " --r_tidf=" + ea.TournamentIDFilter
		}
		if len(ea.ExcelSheets) != 0 {
			for _, sheet := range ea.ExcelSheets {
				result_str = result_str + " --sheet=" + sheet
			}
		}
		if len(ea.TournamentURLs) != 0 {
			for _, url := range ea.TournamentURLs {
				result_str = result_str + " --url=" + url
			}
		}
	}
	return result_str
}

func bscPythonTest(command_args []string) (int, string) {
	cmd := exec.Command("python", command_args...) //gosec:disable
	cmd.Dir = "/home/runner/work/skill-calculator/skill-calculator/BSC/src"
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
