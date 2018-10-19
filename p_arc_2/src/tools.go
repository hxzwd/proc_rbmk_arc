
package main


import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"text/tabwriter"
	"bytes"
)




func tools_read_cells_list_data(path string, filename string, cell_coords []string, info_mode_string string) (map[string]rbmk_cell_data, []string) {

	res := make(map[string]rbmk_cell_data, 0)

	const max_line_len = 1024 * 1024

	path = path + "/" + filename

	info_mode := ft_split_string_trim(info_mode_string, ";")

	if ft_in_list(info_mode, "files") {
		fmt.Printf("Target file: %s\n", path)
	}

	var counter int = 0
	data := make([]string, 1024)
	buf := make([]byte, max_line_len)
	hndl, err1 := os.Open(path)
	if err1 != nil {
		if ft_in_list(info_mode, "errors") {
			fmt.Printf("could not open file %s: %v\n", path, err1)
		}
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
 		if ft_in_list(info_mode, "errors") {
			fmt.Printf("scanner error: %v\n", err1)
		}
		os.Exit(1)
	}


	cells_list := f_get_cells_list(data, counter)
	target_cells := []rbmk_cell{}

	if len(cell_coords) == 0 {
		for _, v := range cells_list {
			cell_coords = append(cell_coords, v.coord)
		}
	}

	for _, v := range cell_coords {
		tmp_cell, status := f_find_cell_in_list(cells_list, v)
		if !status &&  ft_in_list(info_mode, "errors") {
			fmt.Printf("cell [coord: %s] not found\n", v)
			continue
		}
		target_cells = append(target_cells, tmp_cell)
	}

	if ft_in_list(info_mode, "cells-info") {
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

	return res, data

}




func tools_get_data_cells_list(path string, filenames []string, cells_coords []string, raw_data_flag bool, info_mode_string string)  (map[string]rbmk_cell_data, map[string]string, map[string]string) {

	cells_types := make(map[string]string, 0)
	raw_data := make(map[string]string, 0)
	big_data, big_raw_data :=  tools_read_cells_list_data(path, filenames[0], cells_coords, info_mode_string)
	for key, value := range big_data {
		cells_types[key] = value.cell.get_type(big_raw_data)
	}
	for i, _ := range filenames {
		if i == 0 {
			continue
		}
		tmp_big_data, _ /* tmp_big_raw_data */ :=  tools_read_cells_list_data(path, filenames[i], cells_coords, info_mode_string)
		for key, value := range tmp_big_data {
			tmp_map_data := sum_data_dumps(big_data[key].data, value.data)
			big_data[key] = rbmk_cell_data{ cell: value.cell, data: tmp_map_data, cell_type: cells_types[key] }
		}

	}

	if raw_data_flag {
		return big_data, cells_types, raw_data
	} else {
		return big_data, cells_types, nil
	}

}


func tools_get_cells_info_2() {


	filename := "/home/hxzwd/go_code/proc_rbmk_arc/p_arc_2/dumps/arc.dump"
	var archive_data archiveData
	data := archive_data.load_text_data(filename)
	archive_data = archive_data.format_text_data(data)


}


func tools_get_cells_info() {



	var path string = "/home/hxzwd/mephim/arc_data/arc1_txt_part"
	var info_mode_string string = "files; errors"
	filenames := []string{ "00.txt", "01.txt", "02.txt", "03.txt", "04.txt", "05.txt", "06.txt", "07.txt", "08.txt", "09.txt", "10.txt", "11.txt" }
	cells_coords := []string{ "44-23", "44-33", "30-51" }


	cells_data, cells_types, _ := tools_get_data_cells_list(path, filenames, cells_coords, false, info_mode_string)

	cells_types = cells_types

	Info := make(map[string]map[string]string, 0)
	Files := filenames
	Cells := cells_coords
	Data := cells_data
	for key, value := range cells_data {
		Info[key] = value.get_info(0)
	}

	var arc_data archiveData = NewEmptyArchiveData()
	arc_data = New(Data, Cells, Info, Files)

	arc_data.print_files("")
	arc_data.print_info("")

	arc_data.dump_text_data("/home/hxzwd/go_code/proc_rbmk_arc/p_arc_2/dumps/arc.dump")

//	for key, value := range cells_data {
//		fmt.Printf("key: %s[%s]\n", key, cells_types[key])
//		value.print_info(0, "")
//	}

}





func tools_get_archive_data(path string, filenames []string, info_mode_string string, cells_coords []string, raw_data_flag bool) (archiveData, map[string]string) {


	if info_mode_string == "default" {
		info_mode_string =  "files; errors"
	}


	cells_data, cells_types, cells_raw_data := tools_get_data_cells_list(path, filenames, cells_coords, raw_data_flag, info_mode_string)

	cells_types = cells_types

	Info := make(map[string]map[string]string, 0)
	Files := filenames
	Cells := cells_coords
	Data := cells_data
	for key, value := range cells_data {
		Info[key] = value.get_info(0)
	}

	var arc_data archiveData = NewEmptyArchiveData()
	arc_data = New(Data, Cells, Info, Files)

	if raw_data_flag {
		return arc_data, cells_raw_data
	} else {
		return arc_data, nil
	}

}

func tools_dump_archive_data(path string, filename string, data archiveData) bool {

	var target_path string
	if path != "" {
		target_path = path + "/" + filename
	} else {
			target_path = filename
	}

	status := data.dump_text_data(target_path)
	return status

}


func tools_read_archive_dump(path string, filename string) archiveData {

	var target_path string
	if path != "" {
		target_path = path + "/" + filename
	} else {
			target_path = filename
	}

	var arc_data archiveData
	raw_data := arc_data.load_text_data(target_path)
	arc_data = arc_data.format_text_data(raw_data)

	return arc_data

}



func tools_filter_archive_columns(arc_data archiveData, cells []string, fields map[string][]string, flag_not_none bool) (map[string][]string, []string, []string) {

	table_data := make(map[string][]string, 0)
	table_head := make([]string, 0)
	rows_index := make([]string, 0)

	if ft_in_map_s_ls(fields, "all") {
		tmp_0 := fields["all"]
		fields = make(map[string][]string, 0)
		for _, cell := range cells {
			fields[cell] = tmp_0
		}
	}

	for _, cell := range cells {
		for _, field := range fields[cell] {
			column_name := cell + ":" + field
			column_data := arc_data.at(cell, field)
			if flag_not_none {
				if ft_in_list(column_data, "none") {
					continue
				}
			}
			if column_data == nil {
				continue
			} else {
				table_data[column_name] = column_data
				table_head = append(table_head, column_name)
			}
		}
	}

	for _, column := range table_data {
		rows_index = f_int2str(ft_make_range(0, len(column), 1))
		break
	}

	return table_data, table_head, rows_index

}



func tools_get_table_view(arc_data archiveData, cells []string, fields map[string][]string, flag_not_none bool, options string) ([]string, string, [][]string) {

	var cols_num int
	var rows_num int
	var start_index int = 0
	var indexing string = "true" 
	var header string = ""
	table_data := make([]string, 0)


	table_columns, table_head, rows_index := tools_filter_archive_columns(arc_data, cells, fields, flag_not_none)

	cols_num = len(table_columns)
	rows_num = len(rows_index)

	columns := make([][]string, cols_num)

	for i, column_name := range table_head {
		columns[i] = table_columns[column_name]
		header = header + "\t" + column_name
	}
	header = strings.Trim(header, "\t\n ")

	if options == "" || options == "default" {
		options = fmt.Sprintf("rn = %d; b = %d; n = %s; h = %s;", rows_num, start_index, indexing, header)
	}

	table_data, header = ft_fprint_list_of_columns(options, columns, cols_num)

	return table_data, header, columns

}

func tools_write_data_csv(path string, filename string, csv_data []string, columns_head string, rows_num int) int {

	var target_path string
	var n_bytes int = 0

	if path != "" {
		target_path = path + "/" + filename
	} else {
			target_path = filename
	}

	if rows_num <= 0 {
		rows_num = len(csv_data)
	}

	var bytes_writed_num int = 0


	hndl, err0 := os.Create(target_path)

	if err0 != nil {
		fmt.Printf("Error: tools_write_data_csv() could not open output file %s: %v\n", target_path, err0)
		os.Exit(1)
	}
	csv_writer := bufio.NewWriter(hndl)

	if columns_head != "" {
		n_bytes, err0 = csv_writer.Write([]byte(columns_head))
		if err0 != nil {
			fmt.Printf("Error: tools_write_data_csv() Write() row index: [HEADER]\ttarget path: %s: %v\n", target_path, err0)
		os.Exit(1)
		}
		bytes_writed_num += n_bytes
	}

	for i := 0; i < rows_num; i++ {
		n_bytes, err0 = csv_writer.Write([]byte(csv_data[i]))
		if err0 != nil {
			fmt.Printf("Error: tools_write_data_csv() Write() row index: %d\ttarget path: %s: %v\n", i, target_path, err0)
		os.Exit(1)
		}
		bytes_writed_num += n_bytes
	}

	csv_writer.Flush()

	hndl.Sync()
	hndl.Close()

	fmt.Printf("csv: %d bytes writed in %s [path: %s]\n", bytes_writed_num, filename, path)

	return bytes_writed_num

}


func tools_dump_csv(path string, filename string, arc_data archiveData, cells []string, fields map[string][]string, flag_not_none bool, options string, csv_header string) int {

	var bytes_writed_num int

	table_data, table_head, _ := tools_get_table_view(arc_data, cells, fields, flag_not_none, options)

	if csv_header == "default" {
		csv_header = table_head
	}

	bytes_writed_num = tools_write_data_csv(path, filename, table_data, csv_header, 0)

	return bytes_writed_num

}

func tools_read_data_csv(path string, filename string, options string) []string {

	const max_line_len = 1024 * 1024
	pattern := "[\\s]*(?P<option>[^;^=^\\s]+)[\\s]*[=]{1}[\\s]*(?P<value>[^;]+)[\\s]*(?:;|\\z)"
	var start_row int = 0
	var num_of_rows int = -1 

	var target_path string = u_get_full_path(path, filename)
	var line_buffer string
	var line_counter int = 0
	var line_total_counter int = 0
	lines_buffer := make([]string, 0)
	var flag_stdout bool = true
	var flag_stdout_short = false
	var flag_stdout_columns = false
	var flag_stdout_header = false
	var flag_column_view = false
	var short_left_width = 0
	var short_right_width = 0
	columns_id := make([]int, 0) 

	buf := make([]byte, max_line_len)
	if options == "default" {
		flag_stdout = true
	} else {
		params := u_re_find_all(pattern, options)
		for _, param := range params {
			switch param["option"] {
				case "print":
					flag_stdout = u_str2bool(param["value"])
				case "start":
					start_row = u_str2int(param["value"])
				case "nlines":
					num_of_rows = u_str2int_2(param["value"], -1)
				case "short":
					tmp := f_str2int(strings.Split(param["value"], " "))
					flag_stdout_short = true
					short_left_width = tmp[0]
					short_right_width = tmp[1]
				case "columns":
					columns_id = f_str2int(strings.Split(param["value"], " "))
					flag_stdout_columns = true
					flag_stdout = true
				case "header":
					flag_stdout_header = u_str2bool(param["value"])
				case "column_view":
					flag_column_view = u_str2bool(param["value"])
			}
		}
	}

	if start_row < 0 {
		start_row = 0
	}
	if num_of_rows <= 0 {
		num_of_rows = -1
	}

	if !flag_stdout || short_left_width <= 0 || short_right_width <= 0 {
		flag_stdout_short = false
	}

	fn, err := os.Open(target_path)
	
	if err != nil {
		fmt.Printf("Error: tools_read_data_csv() could not open output file %s: %v\n", target_path, err)
		os.Exit(1)
	}
	defer fn.Close()
	
	csv_reader := bufio.NewScanner(fn)
	csv_reader.Buffer(buf, max_line_len)

	for csv_reader.Scan() {
		line_buffer = csv_reader.Text()
		if num_of_rows != -1 {
			if line_total_counter >= num_of_rows {
				break
			}
		} 
		if start_row <= line_counter {
			lines_buffer = append(lines_buffer, line_buffer)
			line_total_counter++
			line_counter++
		} else {
			if start_row != 0 && line_counter == 0 && flag_stdout_header && flag_stdout {
				line_buffer = u_align_words_3(line_buffer, "*", 16, 0, true)
				line_counter++
			} else {
				line_counter++
				continue
			}
		}
		if flag_stdout_columns {
			tmp1 := tools_extract_columns_from_string(line_buffer, columns_id)
			if flag_column_view {
				for _, v := range strings.Fields(tmp1) {
					fmt.Println(v)
				}
				continue
			}
			fmt.Printf("%s", tmp1)
			continue
		}
		if flag_stdout_short {
			tmp1 := len(line_buffer)
			fmt.Printf("%s\t. . .    %s\n", strings.TrimSpace(line_buffer[0:short_left_width]), strings.TrimSpace(line_buffer[tmp1 - short_right_width:tmp1]))
		}
		if flag_stdout && flag_stdout_short == false {
			if flag_column_view {
				for _, v := range strings.Fields(line_buffer) {
					fmt.Println(v)
				}
				continue
			}
			fmt.Println(line_buffer)
		}

	}

	err = csv_reader.Err()

	if err != nil {
		fmt.Printf("Error: tools_read_data_csv() scanner error: %v\n", err)
		os.Exit(1)
	}

	if !flag_stdout {
		return lines_buffer
	} else {
		return nil
	}

}

func tools_extract_columns_from_string(row string, cols_id []int) string {

	
	var res string = ""
	var buf bytes.Buffer
	tab_writer := new(tabwriter.Writer)
	tab_writer.Init(&buf, 0, 8, 2, '\t', 0)

	columns := strings.Fields(row)

	for index, value := range columns {
		if len(value) <= 3 && u_is_integer(value) && index > 0 {
			value = value + "." + strings.Repeat("0", 16)
		}
		columns[index] = strings.TrimSpace(value)
	}

	for _, column_id := range cols_id {

		res = res + fmt.Sprintf("%s\t", columns[column_id])
	}
	fmt.Fprintln(tab_writer, res)
	tab_writer.Flush()
	res = strings.TrimSpace(buf.String()) + "\n" 
	return res
}



