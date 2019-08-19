package data_formatter

import (
	m "metrics"
	"strconv"
)

type folderStack struct {
	n     uint64
	funcs string
}

func CreateFoldedStacktrace(s map[string]m.Stacktrace) string {

	var fs []folderStack

	for _, stv := range s {
		var rootFunction string

		for fk, fv := range stv.Functions {
			if fv.Caller == "root" {
				rootFunction = fk
				break
			}
		}

		var stack []string

		recursiveFolding(stv, rootFunction, rootFunction, &stack)

	outer:
		for _, s := range stack {
			for i, f := range fs {
				if s == f.funcs {
					fs[i].n += stv.Counter
					continue outer
				}
			}
			fs = append(fs, folderStack{n: stv.Counter, funcs: s})
		}
	}

	folded := ""

	for _, f := range fs {
		folded = folded + f.funcs + " " + strconv.FormatUint(f.n, 10) + "\n"
	}

	return folded
}

func recursiveFolding(stv m.Stacktrace, rootFunction string, tmpFuncs string, funcsSlice *[]string) bool {

	var last = true

	for fk, fv := range stv.Functions {
		if fv.Caller == rootFunction {
			last = false

			tmpInnerFuncs := tmpFuncs + "; " + fk

			if recursiveFolding(stv, fk, tmpInnerFuncs, funcsSlice) == true {
				*funcsSlice = append(*funcsSlice, tmpInnerFuncs)
			}
		}
	}

	return last
}
