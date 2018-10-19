
package main

import (
	"fmt"
	"regexp"
	"os"
	"strings"
	"strconv"
)


func u_re_find_all(pattern string, target_string string) []map[string]string {

		re := regexp.MustCompile(pattern)

		all_indexes := re.FindAllSubmatchIndex([]byte(target_string), -1)

		all_groups := re.SubexpNames()

		res := make([]map[string]string, len(all_indexes))

		for i, submatch := range all_indexes {
			tmp_0 := make(map[string]string, 0)
			for j, group := range all_groups {
				if group == "" {
					tmp_0[fmt.Sprintf("%d", j)] = string(re.Expand([]byte(""), []byte(fmt.Sprintf("$%d", j)), []byte(target_string), submatch))
				} else {
					tmp_0[group] = string(re.Expand([]byte(""), []byte("$" + group), []byte(target_string), submatch))
					tmp_0[fmt.Sprintf("%d", j)] = string(re.Expand([]byte(""), []byte(fmt.Sprintf("$%d", j)), []byte(target_string), submatch))
				}
			}
			res[i] = tmp_0
		}

		return res

}

func u_get_full_path(path string, filename string) string {

	pattern := "^.*(\\|/)$"
	re := regexp.MustCompile(pattern)
	var full_path string = ""

	if path == "" {
		return strings.TrimRight(filename, "/\\")
	}

	if re.Match([]byte(path)) {
		path = strings.TrimRight(path, "/\\")
	}
	filename = strings.Trim(filename, "/\\")
	full_path = path + "/" + filename

	return full_path
}

func u_is_file_exist(path string, filename string) bool {

	var target_path string = u_get_full_path(path, filename)

	if _, err := os.Stat(target_path); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}

}

func u_str2int(value string) int {
	res, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return res
}

func u_str2int_2(value string, default_value int) int {
	res, err := strconv.Atoi(value)
	if err != nil {
		return default_value
	}
	return res
}

func u_str2bool(value string) bool {
	re_true := regexp.MustCompile(`(?i:true|yes|1|y|on|enable|\+|up)`)
	re_false := regexp.MustCompile(`(?i:false|no|0|on|off|disable|\-|down)`)
	if re_true.Match([]byte(value)) {
		return true
	}
	if re_false.Match([]byte(value)) {
		return false
	}
	return true
}


func u_str2bool_2(value string, default_value bool) bool {
	re_true := regexp.MustCompile(`(?i:true|yes|1|y|on|enable|\+|up)`)
	re_false := regexp.MustCompile(`(?i:false|no|0|on|off|disable|\-|down)`)
	if re_true.Match([]byte(value)) {
		return true
	}
	if re_false.Match([]byte(value)) {
		return false
	}
	return default_value
}

func u_words_length(list []string) []int {
	res := make([]int, 0)
	for _, word := range list {
		res = append(res, len(word))
	}
	return res
}

func u_align_words(list []string, stub string, target_len int) []string {
	lengths := u_words_length(list)

	if stub == "" {
		stub = " "
	}
	
	if target_len < 0 {
		target_len = ft_min_list_i(lengths)
		for i, v := range list {
			list[i] = v[0:target_len]
		}
	}

	if target_len == 0 {
		target_len = ft_max_list_i(lengths)
	}

	if target_len >= 0 {
		for i, v := range list {
			list[i] = v + strings.Repeat(stub, target_len - lengths[i])
		}
	}

	return list
}

func u_align_words_both_sides(list []string, stub string, target_len int) []string {
	lengths := u_words_length(list)

	if stub == "" {
		stub = " "
	}
	
	if target_len < 0 {
		target_len = ft_min_list_i(lengths)
		for i, v := range list {
			list[i] = v[0:target_len]
		}
	}

	if target_len == 0 {
		target_len = ft_max_list_i(lengths)
	}

	if target_len >= 0 {
		for i, v := range list {
			tmp1 := (target_len - lengths[i])/2
			list[i] = strings.Repeat(stub, tmp1) + v + strings.Repeat(stub, tmp1)
		}
	}

	return list
}

func u_align_words_2(value string, stub string, target_len int) string {
	
	var res string = ""
	list := strings.Fields(value)
	list = u_align_words(list, stub, target_len)

	res = strings.Join(list[:], "\t")
	return res
}

func u_align_words_3(value string, stub string, target_len, exception int, flag_both_side bool) string {

	var res string = ""
	var excp string = ""
	list := strings.Fields(value)
	excp = list[exception]
	if flag_both_side {
		list = u_align_words_both_sides(list[1:len(list)], stub, target_len)
	} else {
		list = u_align_words(list[1:len(list)], stub, target_len)
	}
	

	res = strings.Join(list[:], "\t")
	return excp + "\t" + res

}

func u_is_number(value string) bool {
	value = strings.TrimSpace(value)
	re := regexp.MustCompile(`^[\+\-]{0,1}[0-9]*[\.]{0,1}[0-9]*$`)
	if re.Match([]byte(value)) {
		return true
	} else {
		return false
	}
}

func u_is_integer(value string) bool {
	value = strings.TrimSpace(value)
	re := regexp.MustCompile(`^[\+\-]{0,1}[0-9]*[\.]{0}[0-9]*$`)
	if re.Match([]byte(value)) {
		return true
	} else {
		return false
	}
}
