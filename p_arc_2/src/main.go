
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"text/tabwriter"
)


//
//type rbmk_cell_data struct {
//	cell rbmk_cell
//	data map[string][]string
//}
//


func sum_data_dumps(data_dump ...map[string][]string) map[string][]string {

	res := make(map[string][]string)
	for index, value := range data_dump {
		if index == 0 {
			res = value
			continue
		}
		for k, v := range value {
			res[k] = append(res[k], v...)
		}
	}

	return res

}

func write_data_csv(filename string, csv_data []string, num_of_rows int) {


	tab_writer := new(tabwriter.Writer)


	var num_w_bytes int = 0


	f, err0 := os.Create(filename)

	if err0 != nil {
		fmt.Printf("write_data_csv() could not open output file %s: %v\n", filename, err0)
		os.Exit(1)
	}
	file_writer := bufio.NewWriter(f)
	tab_writer.Init(file_writer, 0, 8, 2, '\t', tabwriter.AlignRight)
	for i := 0; i < num_of_rows; i++ {
		fmt.Fprintln(tab_writer, strings.Trim(csv_data[i], "\n") + "\t")
	}
	fmt.Fprintln(tab_writer)
	file_writer.Flush()
	tab_writer.Flush()

	f.Sync()
	f.Close()

	fmt.Printf("%d bytes writed in %s\n", num_w_bytes, filename)

}


func read_data_main(path string, filename string, cell_coord string, flag_info bool) (map[string][]string, rbmk_cell, string, []string, int) {


	const max_line_len = 1024 * 1024

	path = path + "/" + filename
	fmt.Printf("Target file: %s\n", path)
	var counter int = 0
	data := make([]string, 1024)
	buf := make([]byte, max_line_len)
	hndl, err1 := os.Open(path)
	if err1 != nil {
		fmt.Printf("could not open file %s: %v\n", path, err1)
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


	cells_list := f_get_cells_list(data, counter)
	cell, status := f_find_cell_in_list(cells_list, cell_coord)

	if !status {
		fmt.Printf("cell [coord: %s] not found\n", cell_coord)
		return nil, cell, "", data, counter
	}

	if flag_info {
		cell.print_cell_info(data)
	}

	all_data := cell.get_all_data(data)
	all_data_offset := cell.get_all_offset_data(data)

	result_data := all_data

	f_is_dkv, _ := cell.is_dkv(data)
	if f_is_dkv {
		result_data = all_data_offset
	} else {
		result_data = all_data
	}

	if flag_info {
		fmt.Println(len(all_data))
	}

	return result_data, cell, cell.get_type(data), data, counter

}


func form_csv_table(data map[string][]string, cell_type string) ([]string, int) {

	table_len := len(data["param"]) + 2

	table_data := make([]string, table_len)
	var header string = "param\t"
	var table_row string = ""
	var row_counter int = 1
	for key, value := range data {
		if (key == strings.ToLower(cell_type) && strings.Contains(cell_type, "DKV")) || f_is_none_data(value) || key == "param" {
			continue
		}
		header = header +  key + "\t"
	}
	table_data[0] = strings.TrimSpace(header) + "\n"
	header_list := strings.Split(strings.Trim(header, "\n \t"), "\t")
	for i := 0; i < len(data["param"]); i++ {
		table_row = ""
		for _, key := range header_list {
			table_row = table_row + data[key][i] + "\t"
		}
		table_data[i + 1] = strings.TrimSpace(table_row) + "\n"
		row_counter++
	}
	return table_data, row_counter

}


func form_csv_table_without_param(data map[string][]string, cell_type string) ([]string, int) {

	fields := make([]string, 0)
	for key, _ := range data {
		fields = append(fields, key)
	}

	table_len := len(data[fields[0]]) + 2

	table_data := make([]string, table_len)
	var header string = ""
	var table_row string = ""
	var row_counter int = 1
	for key, value := range data {
		if (key == strings.ToLower(cell_type) && strings.Contains(cell_type, "DKV")) || f_is_none_data(value) || key == "param" {
			continue
		}
		header = header +  key + "\t"
	}
	table_data[0] = strings.TrimSpace(header) + "\n"
	header_list := strings.Split(strings.Trim(header, "\n \t"), "\t")
	for i := 0; i < len(data[fields[0]]); i++ {
		table_row = ""
		for _, key := range header_list {
			table_row = table_row + data[key][i] + "\t"
		}
		table_data[i + 1] = strings.TrimSpace(table_row) + "\n"
		row_counter++
	}
	return table_data, row_counter
}


func filter_table_columns(table_data map[string][]string, fields []string) map[string][]string {

	res := make(map[string][]string)
	var columns_str string = ""

	for _, value := range fields {
		columns_str += value + "\t"
	}
	columns_str = strings.TrimSpace(columns_str)

	for key, value := range table_data {
		if strings.Contains(columns_str, key) {
			res[key] = value
			continue
		}
	}

	return res


}

func main_old() {


	fmt.Println("Proccess rbmk-1000 archive data\n")

	var script_file string = "p_arc.script"
	script_data := f_load_script(script_file)



	var path string = "/home/hjk/mephim/arc_data/arc1_txt_part"
	var output_file string = "out2.txt"
	files := []string{ "00.txt", "01.txt", "02.txt", "03.txt", "04.txt", "05.txt", "06.txt", "07.txt", "08.txt", "09.txt", "10.txt", "11.txt" }
	var cell_coord string = "30-51"
	fields := []string{ "param", "P", "dkv1_of", "dkv2_of", "W", "G" }

	cell_coord_inputs := f_script_inputs(script_data, "cell")
	if cell_coord_inputs != nil {
		cell_coord = cell_coord_inputs[0]
		fmt.Printf("cell_coord is %s\n", cell_coord)
	}
	output_file_inputs := f_script_inputs(script_data, "output-file")
	if output_file_inputs != nil {
		output_file = output_file_inputs[0]
		fmt.Printf("output_file is %s\n", output_file)
	}
	fields_inputs := f_script_inputs(script_data, "fields")
	if fields_inputs != nil {
		fields = fields_inputs
		fmt.Printf("fields is ")
		fmt.Println(fields)
	}


/////////////////////////////////////////////////////////////////////

	cell_coords := []string{ "34-43", "34-33", "44-23", "44-33", "44-43", "34-23", "34-53", "44-53", "24-23", "24-33", "24-43", "24-53", "54-53", "54-43", "54-33", "54-23" }
	cell_types := make(map[string]string, 0)
	big_data, big_raw_data := read_data_main_for_cell_list(path, files[0], cell_coords, false)
	for key, value := range big_data {
		cell_types[key] = value.cell.get_type(big_raw_data)
	}
	for i, _ := range files {
		if i == 0 {
			continue
		}
		tmp_big_data, _ /* tmp_big_raw_data */ := read_data_main_for_cell_list(path, files[i], cell_coords, false)
		for key, value := range tmp_big_data {
			tmp_map_data := sum_data_dumps(big_data[key].data, value.data)
			big_data[key] = rbmk_cell_data{ cell: value.cell, data: tmp_map_data, cell_type: cell_types[key] }
		}

	}


//	for k, v := range big_data {
//		fmt.Printf("[%s]%d\n", k, len(v.data["P"]))
//		v.cell.print_field_data(big_raw_data, "P", 0, 3)
//	}

	fmt.Println(len(big_data))

	big_fields := []string{ "param", "P", "W", "G" }
	big_no_param_field := false
	var big_fields_str string = ""
	for _, v := range big_fields {
		big_fields_str = big_fields_str + v + "\t"
	}
	big_fields_str = strings.TrimSpace(big_fields_str)
	if !strings.Contains(big_fields_str, "param") {
		big_no_param_field = true
	}

	for key, value := range big_data {

		big_cell_type := cell_types[key]
		big_output_file := "data/" + strings.ToLower(big_cell_type) + "_" + key + ".txt"
		tmp_final_data := filter_table_columns(value.data, big_fields)
		tmp_table_data, tmp_row_counter := form_csv_table(tmp_final_data, big_cell_type)
		if big_no_param_field {
			tmp_table_data, tmp_row_counter = form_csv_table_without_param(tmp_final_data, big_cell_type)
		}

		write_data_csv(big_output_file, tmp_table_data, tmp_row_counter)
	}


	os.Exit(0)

/////////////////////////////////////////////////////////////////////

	data, cell, cell_type, raw_data, raw_data_counter := read_data_main(path, files[0], cell_coord, false)



/*
	for _, v := range script_data {
		opt, arg := f_get_option(v)
		if opt != "" {
			fmt.Println(opt, arg)
		}
	}

	os.Exit(0)
*/

	var close_prog bool = false
	var no_param_field bool = false
	var tmp_fields_str string = ""

	for _, v := range fields {
		tmp_fields_str = tmp_fields_str + v + "\t"
	}
	tmp_fields_str = strings.TrimSpace(tmp_fields_str)
	if !strings.Contains(tmp_fields_str, "param") {
		no_param_field = true
	}

	if f_script_flags(script_data, "cell-list") {
		close_prog = true
		cells_list := f_get_cells_list(raw_data, raw_data_counter)
		for _, v := range cells_list {
			fmt.Println(v.to_str(), v.get_type(raw_data))
		}
	}

	if f_script_flags(script_data, "cell-info") {
		close_prog = true
		cells_list := f_get_cells_list(raw_data, raw_data_counter)
		for _, v := range cell_coord_inputs {
			tmp_cell, f_cell_in := f_find_cell_in_list(cells_list, v)
			if f_cell_in {
				tmp_cell.print_cell_info(raw_data)
			}
		}
//		cell.print_cell_info(raw_data)
	}

	if close_prog {
		os.Exit(0)
	}




	if cell.is_empty() {
		fmt.Printf("Cell not found\n")
	}

	data1, _, _, _, _ := read_data_main(path, files[1], cell_coord, false)
	data2, _, _, _, _ := read_data_main(path, files[2], cell_coord, false)
	data3, _, _, _, _ := read_data_main(path, files[3], cell_coord, false)
	data4, _, _, _, _ := read_data_main(path, files[4], cell_coord, false)
	data5, _, _, _, _ := read_data_main(path, files[5], cell_coord, false)
	data6, _, _, _, _ := read_data_main(path, files[6], cell_coord, false)
	data7, _, _, _, _ := read_data_main(path, files[7], cell_coord, false)
	data8, _, _, _, _ := read_data_main(path, files[8], cell_coord, false)
	data9, _, _, _, _ := read_data_main(path, files[9], cell_coord, false)
	data10, _, _, _, _ := read_data_main(path, files[10], cell_coord, false)
	data11, _, _, _, _ := read_data_main(path, files[11], cell_coord, false)


	final_data := sum_data_dumps(data, data1, data2, data3, data4, data5, data6, data7, data8, data9, data10, data11)

///////////////////////////////////////////

	x := data["dkv1_of"]
	x = x
	y := data["W"]

	z := m_centered(y)

//	f_print_two_columns(y, z, "rn:15 b:0 n:true")
	f_print_n_columns("rn:15 b:0 n:true h:y\tz\tx", y, z, x)

	fmt.Println(m_mean(y))
	fmt.Println(m_ccf(y, x, 10))
	fmt.Println(m_ccf(y, y, 0), "\t", m_std(y), "\t", m_disp(y))

	os.Exit(0)

///////////////////////////////////////////

	final_data = filter_table_columns(final_data, fields)

	table_data, row_counter := form_csv_table(final_data, cell_type)


	if no_param_field {
		table_data, row_counter = form_csv_table_without_param(final_data, cell_type)
	}

	write_data_csv(output_file, table_data, row_counter)

}


func read_data_main_for_cell_list(path string, filename string, cell_coords []string, flag_info bool) (map[string]rbmk_cell_data, []string) {

	res := make(map[string]rbmk_cell_data, 0)

	const max_line_len = 1024 * 1024

	path = path + "/" + filename
	fmt.Printf("Target file: %s\n", path)
	var counter int = 0
	data := make([]string, 1024)
	buf := make([]byte, max_line_len)
	hndl, err1 := os.Open(path)
	if err1 != nil {
		fmt.Printf("could not open file %s: %v\n", path, err1)
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


	cells_list := f_get_cells_list(data, counter)
	target_cells := []rbmk_cell{}

	for _, v := range cell_coords {
		tmp_cell, status := f_find_cell_in_list(cells_list, v)
		if !status {
			fmt.Printf("cell [coord: %s] not found\n", v)
			continue
		}
		target_cells = append(target_cells, tmp_cell)
	}

	if flag_info {
		for _, v := range target_cells {
			v.print_cell_info(data)
		}
	}

	for _, cell := range target_cells {

		all_data := cell.get_all_data(data)
		all_data_offset := cell.get_all_offset_data(data)

		result_data := all_data

		f_is_dkv, _ := cell.is_dkv(data)
		if f_is_dkv {
			result_data = all_data_offset
		} else {
			result_data = all_data
		}
		res[cell.coord] = rbmk_cell_data{ cell: cell, data: result_data, cell_type: cell.get_type(data) }

	}

//	if flag_info {
//		fmt.Println(len(all_data))
//	}

	return res, data

}



func main() {

	tools_get_cells_info_2()

//	tools_get_cells_info()

}

