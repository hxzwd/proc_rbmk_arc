
package main



import (
	"fmt"
	"bufio"
	"strings"
	"os"
)

func f_find_option(options []string, option string) (string, []string, bool) {


	for _, v := range options {
		opt, arg := f_get_option(v)
		if opt == "" {
			continue
		}
		if opt == option {
			return opt, arg, true
		}
	}

	return "", nil, false
}

func f_get_option(line string) (string, []string) {

	var option string
	if !strings.Contains(line, ":") {
		return "", nil
	}
	option = strings.TrimSpace(strings.Split(line, ":")[0])
	arg := strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), ",")
	for i, v := range arg {
		arg[i] = strings.TrimSpace(v)
	}

	return option, arg

}

func f_script_flags(script_data []string, flag_name string) bool {

	opt, arg, fl := f_find_option(script_data, flag_name)
	if fl && arg != nil {
		if arg[0] == "true" && opt == flag_name {
			return true
		}
	}

	return false

}

func f_script_inputs(script_data []string, option string) []string {

	opt, arg, fl := f_find_option(script_data, option)
	if fl && arg != nil {
		if opt == option {
			return arg
		}
	}

	return nil

}

func f_load_script(script_name string) []string {

	const max_line_len = 1024 * 1024

	var path string = script_name
	var counter int = 0
	data := make([]string, 1024)
	buf := make([]byte, max_line_len)
	hndl, err1 := os.Open(path)
	if err1 != nil {
		fmt.Printf("could not open script file %s: %v\n", path, err1)
		os.Exit(1)
	}
	defer hndl.Close()

	scanner := bufio.NewScanner(hndl)
	scanner.Buffer(buf, max_line_len)
	for scanner.Scan() {
		data[counter] = strings.Trim(scanner.Text(), " \t")
		counter++
	}

	err1 = scanner.Err()

	if err1 != nil {
		fmt.Printf("scanner error: %v\n", err1)
		os.Exit(1)
	}

	return data

}


