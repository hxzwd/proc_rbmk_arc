
package main


import (
	"fmt"
	"os"
	"bufio"
	"strings"
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


	filename := "/home/hjk/go_code/p_arc_2/dumps/arc.dump"
	var archive_data archiveData
	data := archive_data.load_text_data(filename)
	archive_data = archive_data.format_text_data(data)


}


func tools_get_cells_info() {



	var path string = "/home/hjk/mephim/arc_data/arc1_txt_part"
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

	arc_data.dump_text_data("/home/hjk/go_code/p_arc_2/dumps/arc.dump")

//	for key, value := range cells_data {
//		fmt.Printf("key: %s[%s]\n", key, cells_types[key])
//		value.print_info(0, "")
//	}

}
