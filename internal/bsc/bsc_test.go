package bsc

import (
	"fmt"
	"strings"
	"testing"
)

func Test_runTestExecution(t *testing.T) {
	tests := []struct {
		name      string
		command   string
		exitCode  int
		outputStr string
	}{
		// {
		// 	name:      "Successful python version check",
		// 	command:   "--version",
		// 	exitCode:  0,
		// 	outputStr: "Python 3",
		// },
		{
			name:      "Successful BSC default execution",
			command:   "main.py",
			exitCode:  0,
			outputStr: "Badminton Skill Calculator",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit_code, output_str := bscPythonTest(tt.command)
			fmt.Print(exit_code, output_str)
			if exit_code != tt.exitCode {
				t.Errorf("Test %d --- bscPythonTest() unexpected exit code = '%d', wanted exit code '%d'", i+1, exit_code, tt.exitCode)
			}
			if strings.Contains(output_str, tt.outputStr) != true {
				t.Errorf("Test %d --- bscPythonTest() unexpected string output = '%s', wanted string output '%s'", i+1, output_str, tt.outputStr)
			}
		})
	}
}

//TODO make tests for arguments compliation
