
package main


import (
	"fmt"
	"strings"
	"os"
	"syscall"
	"unsafe"
	"bytes"
	"text/tabwriter"
	"strconv"
)


type winsize struct {
	Row uint16
	Col uint16
	Xpixel uint16
	Ypixel uint16
}

func ft_key_list_s_(m map[string][]string) []string {
	res := make([]string, 0)
	for key, _ := range m {
		res = append(res, key)
	}
	return res
}


func ft_in_map_s_ls(m map[string][]string, value string) bool {
	return ft_in_list(ft_key_list_s_(m), value)
}

func ft_key_list_s(m map[string]interface{}) []string {
	res := make([]string, 0)
	for key, _ := range m {
		res = append(res, key)
	}
	return res
}


func ft_key_list_sl(m map[string][]interface{}) []string {
	res := make([]string, 0)
	for key, _ := range m {
		res = append(res, key)
	}
	return res
}

func ft_in_map_s_i(m map[string]interface{}, value string) bool {
	return ft_in_list(ft_key_list_s(m), value)
}

func ft_in_map_s_li(m map[string][]interface{}, value string) bool {
	return ft_in_list(ft_key_list_sl(m), value)
}

func ft_i2msls(m map[string]interface{}) map[string][]string {
	res := make(map[string][]string, 0)
	for key, value := range m {
		res[key] = ft_il2sl(value.([]interface{}), len(value.([]interface{})))
	}
	return res
}


func ft_i2msmss(m map[string]interface{}) map[string]map[string][]string {

	res := make(map[string]map[string][]string, 0)
	for key, value := range m {
		tmp1 := value.(map[string]interface{})
		tmp0 := make(map[string][]string, 0)
		for k, v := range tmp1 {
			switch vv := v.(type) {
				case string:
					tmp0[k] = []string{ vv }
				default:
					tmp0[k] = ft_il2sl(v.([]interface{}), len(v.([]interface{})))

			}
		}
		res[key] = tmp0

	}
	return res

}

func ft_i2mss(m interface{}) map[string]string {

	res := make(map[string]string, 0)
	tmp1 := m.(map[string]interface{})
	for k, v := range tmp1 {
		res[k] = v.(string)
	}

	return res

}

func ft_il2sl(list []interface{}, n int) []string {
	if n <= 0 {
		n = len(list)
	}
	res := make([]string, n)
	for i, v := range list {
		res[i] = v.(string)
	}
	return res
}


func ft_sl2il(list []string, n int) []interface{} {
	if n <= 0 {
		n = len(list)
	}
	res := make([]interface{}, n)
	for i, v := range list {
		res[i] = v
	}
	return res
}

func ft_max_i(values ...int) (int, int) {
	var res int = values[0]
	var i int = 0
	n := len(values)
	for i = 1; i < n; i++ {
		if values[i] > res {
			res = values[i]
		}
	}
	return res, i
}


func ft_max_list_i(values []int) int {
	var res int = values[0]
	n := len(values)
	for i := 1; i < n; i++ {
		if values[i] > res {
			res = values[i]
		}
	}
	return res
}


func ft_min_i(values ...int) int {
	var res int = values[0]
	n := len(values)
	for i := 1; i < n; i++ {
		if values[i] < res {
			res = values[i]
		}
	}
	return res
}


func ft_min_list_i(values []int) int {
	var res int = values[0]
	n := len(values)
	for i := 1; i < n; i++ {
		if values[i] < res {
			res = values[i]
		}
	}
	return res
}


func ft_list_add(sep1 string, sep2 string, list ...[]string) []string {

	tmp := make(map[int][]string, 0)
	lens := make([]int, 0)
	res := make([]string, 0)
	for i, v := range list {
		tmp[i] = v
		lens = append(lens, len(v))
	}
	max := ft_max_list_i(lens)
	for i := 0; i < max; i++ {
		res = append(res, sep1)
		for _, v := range list {
			if len(v) <= i {
				res[i] += sep2
			} else {
				res[i] += v[i]
			}
		}
		res[i] = strings.Trim(res[i], sep1)
	}
	return res
}

func ft_make_list_f(target string, args []interface{}, n int) []string {

	res := make([]string, 0)
	for i := 0; i < n; i++ {
		res = append(res, fmt.Sprintf(target, args[i]))
	}
	return res

}

func ft_make_list(value string, n int) []string {

	res := make([]string, 0)
	for i := 0; i < n; i++ {
		res = append(res, value)
	}
	return res
}

func ft_in_list(list []string, value string) bool {

	if ft_index_list(list, value) >= 0 {
		return true
	} else {
		return false
	}

}

func ft_make_range(begin int, end int, step int) []int {
	res := make([]int, 0)
	if step == 0 {
		return nil
	}
	if begin == end {
		for i := 0; i < begin; i++ {
			res = append(res, step)
		}
		return res
	}
	for i := begin; i < end; i += step {
		res = append(res, i)
	}
	return res
}

func ft_reversed(list []interface{}, n int) []interface{} {

	if n == 0 {
		n = len(list)
	}
	if n < 0 {
		return nil
	}
	res := make([]interface{}, 0)
//	nc := copy(res, list)
//	if nc != n {
//		fmt.Printf("Number of copied elements [%d] not equal size of source list [%d]\n", nc, n)
//	}
	for i := n - 1; i >= 0; i-- {
		res = append(res, list[i])
	}
	return res
}

func ft_index_list_all(list []string, value string) []int {

	res := make([]int, 0)

	for i, v := range list {
		if v == value {
			res = append(res, i)
		}
	}

	if len(res) == 0 {
		return nil
	} else {
		return res
	}

}

func ft_index_list(list []string, value string) int {

	for i, v := range list {
		if value == v {
			return i
		}
	}
	return -1
}

func ft_split_string_trim(target string, sep string) []string {

	if !strings.Contains(target, sep) {
		return []string{ strings.TrimSpace(target) }
	}
	res := strings.Split(target, sep)
	for i, v := range res {
		res[i] = strings.TrimSpace(v)
	}
	return res

}



func ft_get_window_width() uint {
	win_size := &winsize{}

	status, _, err0 := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(win_size)))

	if int(status) == -1 {
		fmt.Printf("Error: ft_get_window_width()\terrno:%d\terror:%s\n", err0, err0.Error())
		os.Exit(1)
	}

	return uint(win_size.Col)
}


func dumpList(list []string) string {
	var res string = "["
	var n int = len(list)
	for i := 0; i < n - 1; i++ {
		res += fmt.Sprintf("\"%s\", ", list[i])
	}
	res += fmt.Sprintf("\"%s\"]", list[n - 1])
	return res
}

func dumpMap(m map[string][]string) string {
	var res string = "{\n"
	for k, v := range m {
		res += fmt.Sprintf("\t\"%s\": %s,\n", k, dumpList(v))
	}
	res = strings.Trim(res, ",\n")
	res += "\n}"
	return res
}


func dumpMapStr(m map[string]string) string {
	var res string = "{\n"
	for k, v := range m {
		if k == "data" {
			tmp := strings.Replace(v, "\n", "", -1)
			tmp1 := strings.Split(tmp, "\t")
			tmp2 := make([]string, 0)
			for _, vv := range tmp1 {
				if vv != "" {
					tmp2 = append(tmp2, strings.Trim(vv, "\t\n "))
				}
			}

			res += fmt.Sprintf("\t\"info_%s\": %s,\n", k, dumpList(tmp2))
		} else {
			res += fmt.Sprintf("\t\"%s\": \"%s\",\n", k, v)
		}
	}
	res = strings.Trim(res, ",\n")
	res += "\n}"
	return res
}

/*
func dumpMapStr(m map[string]string) string {
	var res string = "{\n"
	for k, v := range m {
		res += fmt.Sprintf("\t\"%s\": \"%s\",\n", k, v)

	}
	res = strings.Trim(res, ",\n")
	res += "\n}"
	return res
}
*/

func dumpMapCellData(m map[string]rbmk_cell_data) string {
	var res string = "{\n"
	for k, v := range m {
		res += fmt.Sprintf("\t\"%s\": %s,\n", k, v.dump_data())
	}
	res = strings.Trim(res, ",\n")
	res += "\n}"
	return res

}

func dumpMapInfo(m map[string]map[string]string) string {
	var res string = "{\n"
	for k, v := range m {
		res += fmt.Sprintf("\t\"%s\": %s,\n", k, dumpMapStr(v))
	}
	res = strings.Trim(res, ",\n")
	res += "\n}"
	return res
}




////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////


func ft_fprint_list_of_columns(options string, columns [][]string, cols_num int) ([]string, string) {

	if cols_num <= 0 {
		cols_num = len(columns)
	}

	col1 := columns[0]
	n := len(col1)
	col2 := make([]string, n)

	for j := 0; j < n; j++ {
		tmp_row := ""
		for i := 1; i < cols_num; i++ {
			tmp_row = tmp_row + columns[i][j] + "\t"
		}
		col2[j] = strings.TrimSpace(tmp_row)
	}

	rows_data, headers := ft_fprint_two_columns(col1, col2, options)

	return rows_data, headers

}


func ft_fprint_two_columns(column1 []string, column2 []string, options string) ([]string, string) {

	var rows_num int = 0
	var begin_pos int = 0
	var flag_numerate bool = false
	var columns_header string = ""
	var err0 error = nil

	var buf bytes.Buffer
	var result_header string = ""
	result_rows := make([]string, 0)


	pattern := "[\\s]*(?P<option>[^;^=^\\s]+)[\\s]*[=]{1}[\\s]*(?P<value>[^;]+)[\\s]*(?:;|\\z)"

	params := u_re_find_all(pattern, options)


	for _, param := range params {
		switch param["option"] {
			case "rn":
				rows_num, err0 = strconv.Atoi(param["value"])
			case "n":
				if param["value"] == "true" {
					flag_numerate = true
				}
			case "b":
				begin_pos, err0 = strconv.Atoi(param["value"])
			case "h":
				if param["value"] != "" {
					columns_header = param["value"]
				}
				if flag_numerate {
					columns_header = "index\t" + columns_header
				}
		}
		if err0 != nil {
				fmt.Printf("matt: f_print_two_columns() strconv.Atoi() error: %v\n", err0)
				os.Exit(1)
		}
	}


	if rows_num <= 0 {
		rows_num = len(column1)
	}



	tab_writer := new(tabwriter.Writer)
	tab_writer.Init(&buf, 0, 8, 2, '\t', 0 /* tabwriter.AlignRight */)

	if columns_header != "" {
		fmt.Fprintln(tab_writer, columns_header + "\t")
		tab_writer.Flush()
		result_header = buf.String()
		buf.Reset()
	}

	for i := begin_pos; i < rows_num + begin_pos; i++ {
		if flag_numerate {
			fmt.Fprintln(tab_writer, fmt.Sprintf("%d\t", i) + column1[i] + "\t" + column2[i] + "\t")
		} else {
			fmt.Fprintln(tab_writer, column1[i] + "\t" + column2[i] + "\t")
		}
		tab_writer.Flush()
		result_rows = append(result_rows, buf.String())
		buf.Reset()
	}

	return result_rows, result_header

}