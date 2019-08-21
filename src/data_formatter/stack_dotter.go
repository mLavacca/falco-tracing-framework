package data_formatter

import (
	m "metrics"
	"strconv"
	"strings"
)

func CreateDotStacktrace(s map[string]m.Stacktrace) []byte {
	var dot string
	var subgraphRoot string
	var replacer = strings.NewReplacer(":", "_")

	dot += "digraph G {\n"
	dot += "\tcompound=true;\n"
	dot += "\tfontsize=22 label=\"Falco Stacktraces\" rankdir = \"LR\";\n"

	i := 0
	for ks, vs := range s {
		dot = dot + "\tsubgraph " + ks + " {\n"
		if i%2 == 0 {
			dot = dot + "\t\tnode [style=filled,color=lightgrey];\n"
		} else {
			dot = dot + "\t\tnode [style=filled,color=lightblue];\n"
		}

		dot = dot + "\t\t" + ks + " [label=\"" + strings.ToUpper(ks) + " - counter: " +
			strconv.FormatUint(vs.Counter, 10) + "\"," +
			" fontname=\"times-bold\",style=filled,color=white];\n"

		for kf, vf := range vs.Functions {

			caller := replacer.Replace(vf.Caller)
			called := replacer.Replace(kf)

			dot = dot + "\t\t" +
				called + strconv.Itoa(i) + " [label=\"" +
				called + " - latency: " + strconv.FormatUint(vf.Latency, 10) + "\"];\n"

			if vf.Caller != "root" {
				dot = dot + "\t\t" +
					caller + strconv.Itoa(i) +
					" -> " + called + strconv.Itoa(i) + ";\n"
			} else {
				dot = dot + "\t\t" + called + strconv.Itoa(i) +
					" -> " + ks +
					" [style=invis]\n"
				subgraphRoot = kf
			}
		}
		dot = dot + "\t}\n"
		i++
	}

	for j := 0; j < i; j++ {
		dot += "\troot -> " + subgraphRoot + strconv.Itoa(j) + ";\n"
	}

	dot += "\n"
	dot += "\troot[shape=Mdiamond];\n"

	dot += "}"

	return []byte(dot)
}
