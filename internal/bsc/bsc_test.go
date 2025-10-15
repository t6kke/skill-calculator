package bsc

import (
	"fmt"
	"strings"
	"testing"
)

func Test_runTestExecution(t *testing.T) {
	tests := []struct {
		name         string
		command_args []string
		exitCode     int
		outputStr    string
	}{
		{
			name:         "Successful python version check",
			command_args: []string{"--version"},
			exitCode:     0,
			outputStr:    "Python 3",
		},
		{
			name:         "Successful BSC default execution",
			command_args: []string{"main.py", "version"},
			exitCode:     0,
			outputStr:    "Badminton Skill Calculator",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exit_code, output_str := bscPythonTest(tt.command_args)
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
