
package main

import (
	"strconv"
	"fmt"
	"os"
	"math"
	"text/tabwriter"
	"strings"
//	"reflect"
)



func f_split_two(target string, sep string) (string, string) {

	if !strings.Contains(target, sep) {
		return "", ""
	}

	res := strings.Split(target, sep)
	opt := strings.TrimSpace(res[0])
	arg := strings.TrimSpace(res[1])

	return opt, arg
}

func f_print_n_columns(options string, columns ...[]string) {

	col1 := columns[0]
	n := len(col1)
	col2 := make([]string, n)

	for j := 0; j < n; j++ {
		tmp_row := ""
		for i, v := range columns {
			if i == 0 {
				continue
			}
			tmp_row = tmp_row + v[j] + "\t"
		}
		col2[j] = strings.TrimSpace(tmp_row)
	}

	f_print_two_columns(col1, col2, options)

}


func f_print_two_columns(column1 []string, column2 []string, options string) {

	var rows_num int = 0
	var begin_pos int = 0
	var flag_numerate bool = false
	var columns_header string = ""
	var err0 error

	params := strings.Split(options, " ")
	for _, param := range params {
		opt, arg := f_split_two(param, ":")
		if opt == "rn" {
			rows_num, err0 = strconv.Atoi(arg)
			if err0 != nil {
				fmt.Printf("matt: f_print_two_columns() strconv.Atoi() error: %v\n", err0)
				os.Exit(1)
			}
			continue
		}
		if opt == "b" {
			begin_pos, err0 = strconv.Atoi(arg)
			if err0 != nil {
				fmt.Printf("matt: f_print_two_columns() strconv.Atoi() error: %v\n", err0)
				os.Exit(1)
			}
			continue
		}
		if opt == "n" {
			if arg == "true" {
				flag_numerate = true
			}
			continue
		}
		if opt == "h" {
			if arg != "" {
				columns_header = arg
			}
			if flag_numerate {
				columns_header = "index\t" + columns_header
			}
			continue
		}
	}

	if rows_num <= 0 {
		rows_num = len(column1)
	}



	tab_writer := new(tabwriter.Writer)
	tab_writer.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	if columns_header != "" {
		fmt.Fprintln(tab_writer, columns_header + "\t")
	}

	for i := begin_pos; i < rows_num + begin_pos; i++ {
		if flag_numerate {
			fmt.Fprintln(tab_writer, fmt.Sprintf("%d\t", i) + column1[i] + "\t" + column2[i] + "\t")
		} else {
			fmt.Fprintln(tab_writer, column1[i] + "\t" + column2[i] + "\t")
		}
	}

	tab_writer.Flush()
}

func f_str2float(seq []string) []float64 {

	res := make([]float64, 0)

	for _, v := range seq {
		tmp_float, err0 := strconv.ParseFloat(v, 64)
		if err0 != nil {
			fmt.Printf("matt: f_str2float() strconv.ParseFloat() error: %v\n", err0)
			os.Exit(1)
		}
		res = append(res, tmp_float)
	}

	return res

}

func f_float2str(fdata []float64) []string {

	res := make([]string, 0)

	for _, v := range fdata {
		res = append(res, fmt.Sprintf("%.16f", v))
	}

	return res
}

func m_mean(seq []string) float64 {

	fdata := f_str2float(seq)
	N := float64(len(fdata))
	var res float64 = 0.0

	for _, v := range fdata {
		res += v
	}

	return res/N
}

func m_disp(seq []string) float64 {

	fdata := f_str2float(seq)
	N := float64(len(fdata))
	x_mean := m_mean(seq)
	var res float64 = 0.0

	for _, v := range fdata {
		res += math.Pow((v - x_mean), 2.0)
	}

	return res/N
}


func m_disp_0(seq []string) float64 {

	fdata := f_str2float(seq)
	N := float64(len(fdata))
	x_mean := m_mean(seq)
	var res float64 = 0.0

	for _, v := range fdata {
		res += math.Pow((v - x_mean), 2.0)
	}

	return res/(N - 1.0)
}


func m_std(seq []string) float64 {

	fdata := f_str2float(seq)
	N := float64(len(fdata))
	x_mean := m_mean(seq)
	var res float64 = 0.0

	for _, v := range fdata {
		res += math.Pow((v - x_mean), 2.0)
	}

	return math.Pow(res/N, 0.5)
}


func m_std_0(seq []string) float64 {

	fdata := f_str2float(seq)
	N := float64(len(fdata))
	x_mean := m_mean(seq)
	var res float64 = 0.0

	for _, v := range fdata {
		res += math.Pow((v - x_mean), 2.0)
	}

	return math.Pow(res/(N - 1.0), 0.5)
}

func m_centered(seq []string) []string {

	fdata := f_str2float(seq)
	x_mean := m_mean(seq)

	for i, v := range fdata {
		fdata[i] = v - x_mean
	}

	return f_float2str(fdata)
}

func m_centered_normed(seq []string) []string {

	fdata := f_str2float(seq)
	x_mean := m_mean(seq)
	x_std := m_std(seq)

	for i, v := range fdata {
		fdata[i] = (v - x_mean)/x_std
	}

	return f_float2str(fdata)
}

func m_absi(n int) int {

	if n < 0 {
		return -n
	}

	return n
}

func m_ccf(x []string, y []string, k int) float64 {
	X := f_str2float(x)
	Y := f_str2float(y)
	var r float64 = 0
	n := len(x)
	MX := m_mean(x)
	MY := m_mean(y)
	std_X := m_std(x)*math.Pow(float64(n), 0.5)
	std_Y := m_std(y)*math.Pow(float64(n), 0.5)
	for i := 1; i <= n - m_absi(k); i++ {
		if k >= 0 {
			r += (X[i - 1] - MX)*(Y[i + k - 1] - MY)
		} else {
			r += (Y[i - 1] - MY)*(X[i - k - 1] - MX)
		}

	}

	r = r/(std_X*std_Y)

	return r

}
