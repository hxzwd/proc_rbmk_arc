
package main

import (
	"strings"
	"regexp"
	"os"
	"fmt"
)

func f_test_0() (string) {
	return "func.go f_test_0 -- function\n"
}

func f_find_cell_in_list(cell_list []rbmk_cell, cell_coord string) (rbmk_cell, bool) {

	for _, cell := range cell_list {
		if cell.coord == cell_coord {
			return cell, true
		}
	}

	return rbmk_cell{ coord: "", first: -1, last: -1 }, false

}

func f_make_translation_table(flag int) map[string]string {

	tr_table := make(map[string]string)
	tr_table2 := make(map[string]string)

	tr_table["param"] = "Параметер"
	tr_table["load"] = "Загрузка, СКАЛА-Микро"
	tr_table["W"] = "Мощность"
	tr_table["G"] = "Расход"
	tr_table["P"] = "СУЗ"
	tr_table["w"] = "Энерговыработка ТВС реактора (МВт*сут)"
	tr_table["W_T"] = "Мощности ТК по ОНФР-2"
	tr_table["dkr1"] = "ДКЭР1"
	tr_table["dkr2"] = "ДКЭР2"
	tr_table["dkr1_k"] = "ДКЭР1 скоррекстированный"
	tr_table["dkr2_k"] = "ДКЭР2 скоррекстированный"
	tr_table["dkv1"] = "ДКЭВ1"
	tr_table["dkv2"] = "ДКЭВ2"
	tr_table["supp"] = "Запас до кризиса"

	if flag == 0 {
		return tr_table
	}


	for key, value := range tr_table {
		tr_table2[value] = key
	}

	if flag == 1 {
		return tr_table2
	}

	return nil

}

func f_translate_param(param string, tr_table map[string]string) string {

	if tr_table == nil {
		tr_table = f_make_translation_table(f_is_old_param(param))
	}

	return tr_table[param]

}

func f_is_old_param(param string) int {

	pattern, err0 := regexp.Compile("^[a-zA-Z].*")

	if err0 != nil {
		fmt.Printf("f_is_new_param() error: %v\n", err0)
		os.Exit(1)
	}

	if pattern.Match([]byte(param)) {
		return 0
	}

	return 1
}


func f_proc_param_data(param_data []string) []string {
	for i, v := range param_data {
		param_data[i] = strings.TrimSpace(strings.Replace(v, " ", "_", -1))
	}

	return param_data
}

func f_proc_none_data(none_data []string) []string {
	for i, v := range none_data {
		none_data[i] = strings.Trim(strings.Replace(v, "нет", "none", -1), "[] \t")
	}

	return none_data

}

func f_is_none_data(data []string) bool {
	for _, v := range data {
		if strings.Contains(v, "none") {
			return true
		}
	}
	return false
}

func f_proc_dkv_data(dkv_data []string) []string {
	res := []string{}
	var tmp_string string = ""
	for index, value := range dkv_data {
		if index % 4 == 0 && index != 0 || index == len(dkv_data) - 1 {
			res = append(res, strings.TrimSpace(tmp_string))
			tmp_string = ""
		}
		tmp_string = tmp_string + strings.Trim(value, "{}") + " "
	}

	return res
}


func f_get_cells_list(data []string, n int) []rbmk_cell {

	_, lines_per_cell, cells_pos := f_get_data_format(data, n)
	res := []rbmk_cell{}

	pattern, err0 := regexp.Compile("[0-9]{2}-[0-9]{2}")
	if err0 != nil {
		fmt.Printf("f_get_cells_list() error: %v\n", err0)
		os.Exit(1)
	}

	for _, v := range cells_pos {
		tmp_cell := rbmk_cell{ coord: pattern.FindAllString(data[v], -1)[0], first: v, last: v + lines_per_cell - 1 }
		res = append(res, tmp_cell)
	}

	return res
}

func f_get_data_values_list(data []string, params string, n int) []string {

	var data_string string = f_get_data_values_str(data, params, n)
	res := make([]string, strings.Count(data_string, ";"))
	var i int = 0
	for _, v := range strings.Split(data_string, ";") {
		if strings.TrimSpace(v) == "" {
			continue
		}
		res[i] = strings.Replace(strings.TrimSpace(v), ",", ".", -1)
		i++

	}
	return res
}

func f_get_data_values_str(data []string, param string, n int) string {

	var res string = ""

	for i := 0; i < n; i++ {
		if data[i] == "" || !strings.Contains(data[i], ";") {
			continue
		}
		if strings.Contains(data[i], param) {
			res = strings.Trim(strings.SplitN(data[i], ";", 2)[1], " \t")
			break
		}
	}

	return res

}

func f_make_params_list(params string) []string {

	var n int = strings.Count(params, "\n") + 1
	res := make([]string, n)
	copy(res, strings.Split(params, "\n"))
	return res

}

func f_get_data_params(data []string, n int) string {

	var res string = ""

	for i := 0; i < n; i++ {
		if data[i] == "" || !strings.Contains(data[i], ";") {
			continue
		}
		res += strings.Trim(strings.SplitN(data[i], ";", 2)[0], "\t ")
		res += "\n"
	}

	return strings.Trim(res, "\n")

}


func f_get_data_format(data []string, n int) (int, int, []int) {
	pos := [2]int{ -1, -1 }
	for i := 0; i < n; i++ {
		if pos[1] == -1 && data[i] == "" {
			pos[1] = i
			break
		}
		if pos[0] == -1 && !strings.Contains(data[i], ";") {
			pos[0] = i
			continue
		}
	}
	lines_per_cell := pos[1] + 1
	num_of_cells := (int)((n + 1)/lines_per_cell)
	res := make([]int, num_of_cells)
	res[0] = pos[0]
	for i := 1; i < num_of_cells; i++ {
		res[i] = res[i - 1] + lines_per_cell
	}

	return 	num_of_cells, lines_per_cell, res
}


