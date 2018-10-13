
package main

import (
	"fmt"
	"strings"
	"bytes"
	"text/tabwriter"
	"os"
	"encoding/json"
)



type rbmk_cell_data struct {
	cell rbmk_cell			`json:"ccell"`
	data map[string][]string	`json:"cdata"`
	cell_type string		`json:"ccell_type"`
}


func NewRbmkCellData() rbmk_cell_data {
	var Cell rbmk_cell = NewRbmkCell()
	return rbmk_cell_data{ cell: Cell, data: nil, cell_type: "" }
}


func (cell rbmk_cell_data) empty() bool {
	if cell.cell.empty() && cell.data == nil && cell.cell_type == "" {
		return true
	} else {
		return false
	}
	return false
}

/*
func (cell *rbmk_cell_data) empty() bool {
	if cell.cell.empty() && cell.data == nil && cell.cell_type == "" {
		return true
	} else {
		return false
	}
	return false
}
*/

func (cell rbmk_cell_data) get_info(info_mode int) map[string]string {

	info := make(map[string]string, 0)

	if info_mode <= 0 {
		info_mode = 3
	}

	if info_mode >= 1 {
		info["coord"] = cell.cell.coord
	}
	if info_mode >= 2 {
		info["type"] = cell.cell_type
	}
	if info_mode >= 3 {
		var buf bytes.Buffer
		var data_info string = ""
		tab_writer := new(tabwriter.Writer)
		tab_writer.Init(&buf, 0, 8, 2, '\t', 0 /*tabwriter.AlignRight*/)
		fmt.Fprintln(tab_writer, "    field:\tlength:\t")
		for key, value := range cell.data {
			fmt.Fprintln(tab_writer, fmt.Sprintf("    %s\t%d\t", key, len(value)))
		}
		tab_writer.Flush()
		data_info += buf.String()
		info["data"] = strings.Trim(data_info, "\n")
	}

	return info

}

func (cell rbmk_cell_data) print_info(info_mode int, header string) {

	if header != "none" {
		if header == "" {
			fmt.Printf("CELL INFO\n")
		} else {
			fmt.Printf("%s\n", header)
		}
	}

	info := cell.get_info(info_mode)
	for key, value := range info {
		if key != "data" {
			fmt.Printf("%s:\t%s\n", key, value)
	} else {
			fmt.Printf("%s:\n%s\n", key, value)
		}
	}

}

func (cell rbmk_cell_data) dump_data() string {

	var res string = ""
	res = fmt.Sprintf("{\n\t\"ccell\": %s,\n", cell.cell.dump_data())
	res += fmt.Sprintf("\t\"cdata\": %s,\n", dumpMap(cell.data))
	res += fmt.Sprintf("\t\"ccell_type\": \"%s\"\n}", cell.cell_type)
	return res
}

func (cell *rbmk_cell_data) MarshalBinary() ([]byte, error) {
	return json.Marshal(cell)
}

func (cell *rbmk_cell_data) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &cell); err != nil {
		return err
	}
	return nil
}


/*
func (cell *rbmk_cell_data) MarshalBinary() ([]byte, error) {
	return json.Marshal(cell)
}

func (cell *rbmk_cell_data) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &cell); err != nil {
		return err
	}
	return nil
}
*/

/*
func (cell rbmk_cell_data) serialize() []byte {

	bin_data, err0 := &cell.MarshalBinary()
	if err0 != nil {
		fmt.Printf("Error rbmk_cell_data.serialize() [cell: %s]:\t%v\n", cell.cell.to_str(), err0)
		os.Exit(1)
	}
	return bin_data
}
*/

func (cell *rbmk_cell_data) serialize() []byte {

	bin_data, err0 := cell.MarshalBinary()
	if err0 != nil {
		fmt.Printf("Error *rbmk_cell_data.serialize() [cell: %s]:\t%v\n", cell.cell.to_str(), err0)
		os.Exit(1)
	}
	return bin_data
}
/*
func (cell rbmk_cell_data) deserialize(bin_data []byte) bool {
	err0 := &cell.UnmarshalBinary(bin_data)
	if err0 != nil {
		fmt.Printf("Error rbmk_cell_data.deserialize() [cell: %s]:\t%v\n", cell.cell.to_str(), err0)
		os.Exit(1)
	}
	return !cell.empty()
}
*/

func (cell *rbmk_cell_data) deserialize(bin_data []byte) bool {
	err0 := cell.UnmarshalBinary(bin_data)
	if err0 != nil {
		fmt.Printf("Error *rbmk_cell_data.deserialize() [cell: %s]:\t%v\n", cell.cell.to_str(), err0)
		os.Exit(1)
	}
	return !cell.empty()
}


/*
func (cell rbmk_cell_data) serialize() []byte {

	var bin_data bytes.Buffer
	encoder := gob.NewEncoder(&bin_data)
	err0 := encoder.Encode(cell)
	if err0 != nil {
		fmt.Printf("Error rbmk_cell.serialize() [%s]:\t%v\n", cell.to_str(), err0)
		os.Exit(1)
	}
	return bin_data.Bytes()

}


func (cell rbmk_cell_data) deserialize(bin_data []byte) bool {

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
