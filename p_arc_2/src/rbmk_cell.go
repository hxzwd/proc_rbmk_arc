
package main


import (
	"strings"
	"regexp"
	"os"
	"fmt"
	"strconv"
	"encoding/json"
)



type rbmk_cell struct {
	coord string	`json:"coord"`
	first int	`json:"first"`
	last int	`json:"last"`
}


func NewRbmkCell(/*cell *rbmk_cell*/) rbmk_cell {
/*
	res := new(cell)
	res.coord = ""
	res.first = 0
	res.last = 0
	return res
*/
	return rbmk_cell{ coord: "", first: 0, last: 0 }
}


func (cell rbmk_cell) empty() bool {
	if cell.coord == "" && cell.first == 0 && cell.last == 0 {
		return true
	} else {
		return false
	}
}

/*
func (cell *rbmk_cell) empty() bool {
	if cell.coord == "" && cell.first == 0 && cell.last == 0 {
		return true
	} else {
		return false
	}
}
*/

func (cell rbmk_cell) is_dkv(data []string) (bool, string) {

	match, err0 := regexp.MatchString("DKV[1-2]{1}", cell.get_type(data))

	if err0 != nil {
		fmt.Printf("(rbmk_cell) is_dkv() regexp.MatchString() error: %v\n", err0)
		os.Exit(1)
	}

	if match {
		return  true, strings.ToLower(cell.get_type(data))
	}

	return false, ""
}


func (cell rbmk_cell) to_str() string {
	return fmt.Sprintf("%s[%d, %d]", cell.coord, cell.first, cell.last)
}

/*
func (cell *rbmk_cell) to_str() string {
	return fmt.Sprintf("%s[%d, %d]", cell.coord, cell.first, cell.last)
}
*/

func (cell rbmk_cell) is_empty() bool {
	if cell.coord == "" {
		return true
	}
	return  cell.coord == "" && cell.first < 0 && cell.last < 0
}


func (cell rbmk_cell) get_type(data []string) string {
	if !strings.Contains(f_get_data_values_str(data[cell.first:cell.last + 1], "СУЗ", cell.last - cell.first), "нет") {
		return "SUZ"
	}
	if !strings.Contains(f_get_data_values_str(data[cell.first:cell.last + 1], "ДКЭВ1", cell.last - cell.first), "нет") {
		return "DKV1"
	}
	if !strings.Contains(f_get_data_values_str(data[cell.first:cell.last + 1], "ДКЭВ2", cell.last - cell.first), "нет") {
		return "DKV2"
	}
	return "UNKNOWN"
}

func (cell rbmk_cell) get_values(data []string, params string) []string {
	list_values := f_get_data_values_list(data[cell.first:cell.last + 1], params, cell.last - cell.first)
	return list_values
}

func (cell rbmk_cell) get_values2(data []string, params string) []string {
	list_values := f_get_data_values_list(data[cell.first:cell.last + 1], f_translate_param(params, nil), cell.last - cell.first)

	for _, v := range list_values {
		if strings.Contains(v, "нет") {
			return f_proc_none_data(list_values)
		}
	}

	match, err0 := regexp.MatchString("dkv[1-2]{1}-DKV[1-2]{1}", params + "-" + cell.get_type(data))

	if err0 != nil {
		fmt.Printf("(rbmk_cell) get_values2() regexp.MatchString() error: %v\n", err0)
		os.Exit(1)
	}

	if match {
		return  f_proc_dkv_data(list_values)
	}


	if params == "param" {
		return f_proc_param_data(list_values)
	}

	return list_values
}

func (cell rbmk_cell) print_cell_info(data []string) {
	cell_type := cell.get_type(data)
	fields_len := cell.get_fields_len(data)
	cell_params_new := cell.get_params(data, 1)
	fmt.Printf("CELL:\n%s [%s]\nBEGIN:\t%d\nEND:\t%d\n", cell.coord, cell_type, cell.first, cell.last)
	for index, value := range cell_params_new {
		fmt.Printf("%d:\t%s => %s\t%d\n", index, f_translate_param(value, nil), value, fields_len[value])
	}
}


func (cell rbmk_cell) get_all_data(data []string) map[string][]string {

	res := make(map[string][]string)
	params_list := cell.get_params(data, 1)
	for _, param := range params_list {
		res[param] = cell.get_values2(data, param)
	}

	return res

}

func (cell rbmk_cell) get_all_offset_data(data []string) map[string][]string {

	all_data := cell.get_all_data(data)
	flag, param := cell.is_dkv(data)
	if !flag {
		return nil
	}

	tmp_data := all_data[param]
	tmp_data_old := make([]string, len(tmp_data))
	copy(tmp_data_old, tmp_data)
	tmp_array := make([]float64, 4)
	var err0 error
	for index, value := range tmp_data {
		for i, v := range strings.Split(value, " ") {
			tmp_array[i], err0 = strconv.ParseFloat(v, 64)
			if err0 != nil {
				fmt.Printf("(rbmk_cell) get_all_offset_data() strconv.ParseFloat() error: %v\n", err0)
				os.Exit(1)
			}
		}
		tmp_data[index] = fmt.Sprintf("%.16f", ((tmp_array[0] + tmp_array[1]) - (tmp_array[2] + tmp_array[3]))/(tmp_array[0] + tmp_array[1] + tmp_array[2] + tmp_array[3]))
	}

	all_data[param + "_of"] = tmp_data
	all_data[param] = tmp_data_old

	return all_data

}

func (cell rbmk_cell) print_field_data(data []string, param string, begin int, end int) {
	data_size := cell.get_fields_len(data)[param]
	if end > data_size || end == 0 {
		end = data_size
	}
	if begin > data_size {
		begin = 0
	}
	if end < 0 {
		end = data_size + end + 1
	}
	if begin < 0 {
		begin = data_size + begin + 1
	}
	field_data := cell.get_values2(data, param)
	for i := begin; i < end; i++ {
		fmt.Printf("%d:  %s\n", i, field_data[i])
	}
}

func (cell rbmk_cell) get_fields_len(data []string) map[string]int {
	params := cell.get_params(data, 1)
	res := make(map[string]int)
	for _, param := range params {
		res[param] = len(cell.get_values2(data, param))
	}

	return res

}

func (cell rbmk_cell) print_fields_len(data []string) {
	len_table := cell.get_fields_len(data)
	for key, value := range len_table {
		fmt.Printf("%s:\t%d (length)\n", key, value)
	}
}

func (cell rbmk_cell) print_params(data []string, mode int) {
	params_list := cell.get_params(data, mode)
	for index, value := range params_list {
		fmt.Printf("%d:\t%s\n", index, value)
	}
}

func (cell rbmk_cell) get_params(data []string, mode int) []string {
	params_list := f_make_params_list(f_get_data_params(data[cell.first:cell.last + 1], cell.last - cell.first))
	if mode == 0 {
		return params_list
	}
	new_params_list := make([]string, len(params_list))
	if mode == 1 {
		for index, param := range params_list {
			new_params_list[index] = f_translate_param(param, nil)
		}
		return new_params_list
	}
	if mode == 2 {
		for index, param := range params_list {
			new_params_list[index] = fmt.Sprintf("%s => %s", param, f_translate_param(param, nil))
		}
		return new_params_list
	}
	if mode == 3 {
		for index, param := range params_list {
			new_params_list[index] = fmt.Sprintf("%s => %s", f_translate_param(param, nil), param)
		}
		return new_params_list
	}

	return nil
}

func (cell rbmk_cell) dump_data() string {
	var res string = ""
	res = fmt.Sprintf("{ \"coord\": \"%s\", \"first\": %d, \"last\": %d }", cell.coord, cell.first, cell.last)
	return res

}

func (cell *rbmk_cell) MarshalBinary() ([]byte, error) {
	return json.Marshal(cell)
}

func (cell *rbmk_cell) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &cell); err != nil {
		return err
	}
	return nil
}

/*
func (cell rbmk_cell) serialize() []byte {

	bin_data, err0 := &cell.MarshalBinary()
	if err0 != nil {
		fmt.Printf("Error rbmk_cell.serialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}
	return bin_data
}
*/

func (cell *rbmk_cell) serialize() []byte {

	bin_data, err0 := cell.MarshalBinary()
	if err0 != nil {
		fmt.Printf("Error *rbmk_cell.serialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}
	return bin_data
}

/*
func (cell rbmk_cell) deserialize(bin_data []byte) bool {
	err0 := &cell.UnmarshalBinary(bin_data)
	if err0 != nil {
		fmt.Printf("Error rbmk_cell.deserialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}
	return !cell.empty()
}
*/

func (cell *rbmk_cell) deserialize(bin_data []byte) bool {
	err0 := cell.UnmarshalBinary(bin_data)
	if err0 != nil {
		fmt.Printf("Error *rbmk_cell.deserialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}
	return !cell.empty()
}

/*

func (cell rbmk_cell) serialize() []byte {

	var bin_data bytes.Buffer
	encoder := gob.NewEncoder(&bin_data)
	err0 := encoder.Encode(cell)
	if err0 != nil {
		fmt.Printf("Error rbmk_cell.serialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}
	return bin_data.Bytes()
}


func (cell rbmk_cell) deserialize(bin_data []byte) bool {

	buf := bytes.NewBuffer(bin_data)
	decoder := gob.NewDecoder(buf)
	err0 := decoder.Decode(&cell)
	if err0 != nil {
		fmt.Printf("Error rbmk_cell.deserialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}

	return !cell.empty()
}

*/
