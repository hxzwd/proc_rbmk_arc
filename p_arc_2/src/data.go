
package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"bytes"
	"text/tabwriter"
	"encoding/json"
)


type ArchiveData interface {
	New() archiveData
	Load() bool
	Set() bool
	num_of_cells() int
	print_files() bool
	print_info() bool
}

type archiveData struct {
	Data map[string]rbmk_cell_data		`json:"data"`
	Cells []string				`json:"cells"`
	Info map[string]map[string]string	`json:"info"`
	Files []string				`json:"files"`
}

func (rbmk_data archiveData) dump_data() string {
	var res string = "{\n"
	res += fmt.Sprintf("\t\"data\": %s,\n", dumpMapCellData(rbmk_data.Data))
	res += fmt.Sprintf("\t\"cells\": %s,\n", dumpList(rbmk_data.Cells))
	res += fmt.Sprintf("\t\"info\": %s,\n", dumpMapInfo(rbmk_data.Info))
	res += fmt.Sprintf("\t\"files\": %s\n", dumpList(rbmk_data.Files))
	res += "}"
	return res
}

func (rbmk_data archiveData) format_text_data(arc_data map[string]interface{}) archiveData {
	archive_Cells := format_archive_cells(arc_data)
	archive_Info := format_archive_info(arc_data)
	archive_Files := format_archive_files(arc_data)
	archive_Data := format_archive_data(arc_data)
	return New(archive_Data, archive_Cells, archive_Info, archive_Files)
}

func NewEmptyArchiveData() archiveData {
	return archiveData{ Data: nil, Cells: nil, Info: nil, Files: nil }
}

func NewEmptyArchiveData2() archiveData {
	return archiveData{ Data: make(map[string]rbmk_cell_data, 0), Cells: make([]string, 0), Info: make(map[string]map[string]string, 0), Files: make([]string, 0) }
}

func New(data map[string]rbmk_cell_data, cells []string, info map[string]map[string]string, files []string) archiveData {
	return archiveData{ Data: data, Cells: cells, Info: info, Files: files }
}

func (rbmk_data archiveData) New(data map[string]rbmk_cell_data, cells []string, info map[string]map[string]string, files []string) archiveData {
	return archiveData{ Data: data, Cells: cells, Info: info, Files: files }
}


func (rbmk_data archiveData) Load(data map[string]rbmk_cell_data, cells []string, info map[string]map[string]string, files []string) bool {
	rbmk_data.Data = data
	rbmk_data.Cells = cells
	rbmk_data.Info = info
	rbmk_data.Files = files
	return true
}


func (rbmk_data archiveData) Set(data archiveData) bool {
	rbmk_data.Data = data.Data
	rbmk_data.Cells = data.Cells
	rbmk_data.Info = data.Info
	rbmk_data.Files = data.Files
	return true
}


func (rbmk_data archiveData) empty() bool {
	if rbmk_data.Data == nil && rbmk_data.Cells == nil && rbmk_data.Info == nil && rbmk_data.Files == nil {
		return true
	} else {
		return false
	}
}

/*
func (rbmk_data *archiveData) empty() bool {
	if rbmk_data.Data == nil && rbmk_data.Cells == nil && rbmk_data.Info == nil && rbmk_data.Files == nil {
		return true
	} else {
		return false
	}
}
*/

func (rbmk_data archiveData) num_of_cells() int {
	var n int = 0
	for _, v := range rbmk_data.Cells {
		if v != "" {
			n++
		}
	}
	return n
}


func (rbmk_data archiveData) print_files(header string) bool {

	if header != "none" {
		if header == "" {
			fmt.Printf("Archive data files:\n")
		} else {
			fmt.Printf("%s\n", header)
		}
	}

	var buf bytes.Buffer
	var message string
	tab_writer := new(tabwriter.Writer)
	tab_writer.Init(&buf, 0, 8, 2, '\t', 0)
	fmt.Fprintln(tab_writer, "index:\tfile name:\t")
	for i, v := range rbmk_data.Files {
		fmt.Fprintln(tab_writer, fmt.Sprintf("%d\t%s\t", i, v))
	}
	tab_writer.Flush()
	message = strings.Trim(buf.String(), "\n")
	fmt.Printf("%s\n", message)
	return true
}



func (rbmk_data archiveData) print_info(header string) bool {

	headers := make([]string, 0)
	var nc int = rbmk_data.num_of_cells()
	if header != "none" {
		if header == "" {
			headers_params := ft_sl2il(rbmk_data.Cells, nc)
			headers = ft_make_list_f("*** Cell[coord: %s] info: ***", headers_params, nc)
		} else {
			headers = ft_make_list(header, nc)
		}
	} else {
		headers = ft_make_list("", nc)
	}

	for i, v := range rbmk_data.Cells {

		fmt.Printf("%s\n", headers[i])
		info := rbmk_data.Info[v]
		for key, value := range info {
			if key != "data" {
				fmt.Printf("%s:\t%s\n", key, value)
			} else {
				fmt.Printf("%s:\n%s\n", key, value)
			}
		}

	}
	return true
}


func (rbmk_data *archiveData) MarshalBinary() ([]byte, error) {
	return json.Marshal(rbmk_data)
}

func (rbmk_data *archiveData) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &rbmk_data); err != nil {
		return err
	}
	return nil
}

/*
func (rbmk_data archiveData) serialize() []byte {

	bin_data, err0 := &rbmk_data.MarshalBinary()
	if err0 != nil {
		fmt.Printf("Error archiveData.serialize():\t%v\n", err0)
		os.Exit(1)
	}
	return bin_data
}
*/

func (rbmk_data *archiveData) serialize() []byte {

	bin_data, err0 := rbmk_data.MarshalBinary()
	if err0 != nil {
		fmt.Printf("Error *archiveData.serialize():\t%v\n", err0)
		os.Exit(1)
	}
	return bin_data
}


/*
func (rbmk_data archiveData) deserialize(bin_data []byte) bool {
	err0 := &rbmk_data.UnmarshalBinary(bin_data)
	if err0 != nil {
		fmt.Printf("Error archiveData.deserialize():\t%v\n", err0)
		os.Exit(1)
	}
	return !rbmk_data.empty()
}
*/

func (rbmk_data *archiveData) deserialize(bin_data []byte) bool {
	err0 := rbmk_data.UnmarshalBinary(bin_data)
	if err0 != nil {
		fmt.Printf("Error *archiveData.deserialize():\t%v\n", err0)
		os.Exit(1)
	}
	return !rbmk_data.empty()
}




func (rbmk_data archiveData) dump_text_data(filename string) bool {

	const max_line_len = 1024 * 1024
	var path string

	path = filename

	hndl, err1 := os.Create(path)
	if err1 != nil {
		fmt.Printf("Error: archiveData.dump_text_data():\t could not open file %s: %v\n", path, err1)
		os.Exit(1)
	}
	defer hndl.Close()

	var bytes_writed int
	byte_data := []byte(rbmk_data.dump_data())
	byte_data_size := len(byte_data)

	dump_writer := bufio.NewWriter(hndl)
	bytes_writed, err1 = dump_writer.Write(byte_data)

	if err1 != nil {
		fmt.Printf("Error: archiveData.dump_text_data(): Write() error: %v\n", err1)
		os.Exit(1)
	}

	err1 = dump_writer.Flush()

	if err1 != nil {
		fmt.Printf("Error: archiveData.dump_text_data(): Flush() error: %v\n", err1)
		os.Exit(1)
	}

	fmt.Printf("%d [total: %d] bytes writed in %s\n", bytes_writed, byte_data_size, path)

	if bytes_writed != byte_data_size {
		return false
	} else {
		return true
	}
}


func (rbmk_data archiveData) load_text_data(filename string) /* archiveData */ map[string]interface{} {

	const max_line_len = 1024 * 1024 * 8
	var path string

	path = filename
	var data string = ""
	var counter int = 0
	buf := make([]byte, max_line_len)

	hndl, err1 := os.Open(path)
	if err1 != nil {
		fmt.Printf("Error: archiveData.load_text_data():\t could not open file %s: %v\n", path, err1)
		os.Exit(1)
	}
	defer hndl.Close()

	scanner := bufio.NewScanner(hndl)
	scanner.Buffer(buf, max_line_len)
	for scanner.Scan() {
		data +=  scanner.Text()
		counter++
	}

	err1 = scanner.Err()

	if err1 != nil {
		fmt.Printf("Error: archiveData.load_text_data():\t scanner error: %v\n", err1)
		os.Exit(1)
	}


	bytes_readed := len([]byte(data))

	fmt.Printf("%d  bytes readed from %s [lines readed: %d]\n", bytes_readed, path, counter)


	/* archive_data := NewEmptyArchiveData2() */
	archive_data := make(map[string]interface{})
	err1 = json.Unmarshal([]byte(data), &archive_data)

	if err1 != nil {
		fmt.Printf("Error: archiveData.load_text_data():\t json.Unmarshal() error: %v\n", err1)
	}

	return archive_data
}


func (rbmk_data archiveData) dump(filename string) bool {

	const max_line_len = 1024 * 1024
	var path string

	path = filename

	hndl, err1 := os.Create(path)
	if err1 != nil {
		fmt.Printf("Error: archiveData.dump():\t could not open file %s: %v\n", path, err1)
		os.Exit(1)
	}
	defer hndl.Close()

	var bytes_writed int
	byte_data := rbmk_data.serialize()
	byte_data_size := len(byte_data)

	dump_writer := bufio.NewWriter(hndl)
	bytes_writed, err1 = dump_writer.Write(byte_data)

	if err1 != nil {
		fmt.Printf("Error: archiveData.dump(): Write() error: %v\n", err1)
		os.Exit(1)
	}

	err1 = dump_writer.Flush()

	if err1 != nil {
		fmt.Printf("Error: archiveData.dump(): Flush() error: %v\n", err1)
		os.Exit(1)
	}

	fmt.Printf("%d [total: %d] bytes writed in %s\n", bytes_writed, byte_data_size, path)

	if bytes_writed != byte_data_size {
		return false
	} else {
		return true
	}
}



func (rbmk_data archiveData) restore(filename string) bool {

	const max_line_len = 1024 * 1024
	var path string

	path = filename

	hndl, err1 := os.Create(path)
	if err1 != nil {
		fmt.Printf("Error: archiveData.dump():\t could not open file %s: %v\n", path, err1)
		os.Exit(1)
	}
	defer hndl.Close()

	var bytes_writed int
	byte_data := rbmk_data.serialize()
	byte_data_size := len(byte_data)

	dump_writer := bufio.NewWriter(hndl)
	bytes_writed, err1 = dump_writer.Write(byte_data)

	if err1 != nil {
		fmt.Printf("Error: archiveData.dump(): Write() error: %v\n", err1)
		os.Exit(1)
	}

	err1 = dump_writer.Flush()

	if err1 != nil {
		fmt.Printf("Error: archiveData.dump(): Flush() error: %v\n", err1)
		os.Exit(1)
	}

	fmt.Printf("%d [total: %d] bytes writed in %s\n", bytes_writed, byte_data_size, path)

	if bytes_writed != byte_data_size {
		return false
	} else {
		return true
	}
}



func archive_get_keys(arc map[string]interface{}) []string {
	res := make([]string, 0)
	for key, _ := range arc {
		res = append(res, key)
	}
	return res
}

func format_archive_info(arc map[string]interface{}) map[string]map[string]string {
	tmp0 := make(map[string]map[string][]string, 0)
	tmp0 = ft_i2msmss(arc["info"].(map[string]interface{}))
	res := make(map[string]map[string]string, 0)
	cells_list := make([]string, 0)
	for cell, _ := range tmp0 {
		cells_list = append(cells_list, cell)
		tmp_map := make(map[string]string, 0)
		tmp_map["coord"] = cell
		tmp_map["type"] = tmp0[cell]["type"][0]
		var info_data string = "    "
		tmp_data := tmp0[cell]["info_data"]
		var  j int = 0
		for i := 0; i < len(tmp_data); i++ {
			if tmp_data[i] == "" {
				continue
			}
			if j % 2 == 0  && j != 0 {
				info_data += "\n    "
			}
			info_data = info_data + tmp_data[i] + "\t"
			j++
		}
		var buf bytes.Buffer
		tab_writer := new(tabwriter.Writer)
		tab_writer.Init(&buf, 0, 8, 2, '\t', 0 /*tabwriter.AlignRight*/)
		fmt.Fprintln(tab_writer, info_data)
		tab_writer.Flush()

		tmp_map["data"] = strings.Trim(buf.String(), "\n")
		res[cell] = tmp_map
	}
	return res
}


func format_archive_cells(arc map[string]interface{}) []string {
 	res := ft_il2sl(arc["cells"].([]interface{}), 0)
	return res
}

func format_archive_files(arc map[string]interface{}) []string {
	res := ft_il2sl(arc["files"].([]interface{}), 0)
	return res
}


func format_archive_data(arc map[string]interface{}) map[string]rbmk_cell_data {

	tmp0 := arc["data"].(map[string]interface{})
	res := make(map[string]rbmk_cell_data, 0)

	for key, value := range tmp0 {
		tmp_cell := rbmk_cell{ coord: "", first: 0, last: 0 }
		for i, v := range (((value.(map[string]interface{}))["ccell"]).(map[string]interface{})) {
			if i == "coord" {
				tmp_cell.coord = v.(string)
				continue
			}
			if i == "first" {
				tmp_cell.first = int(v.(float64))
				continue
			}
			if i == "last" {
				tmp_cell.last = int(v.(float64))
				continue
			}
		}
		tmp_data := (value.(map[string]interface{}))["cdata"].(map[string]interface{})
		tmp_cell_type := (value.(map[string]interface{}))["ccell_type"].(string)
		res[key] = rbmk_cell_data{ cell: tmp_cell, data: ft_i2msls(tmp_data), cell_type: tmp_cell_type }

	}

	return res
}
