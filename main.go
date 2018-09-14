package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	var filename = flag.String("input", "", "input csv filename")
	var columnLengthP = flag.Int("colLen", 0, "csv column length")
	var outfilename = flag.String("output", "", "output csv filename")
	flag.Parse()
	columnLength := *columnLengthP
	log.Printf("colLen=%d\n", columnLength)
	fp, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("input file open error")
		return
	}
	defer fp.Close()

	op, err := os.Create(*outfilename)
	if err != nil {
		log.Fatalf("output file open error")
		return
	}
	defer op.Close()

	reader := bufio.NewReader(fp)
	writer := bufio.NewWriter(op)
	out_words := []string{}
	out_lines := []string{}
	for {
		//行を読み込む
		line, _, err := reader.ReadLine()
		log.Printf("read line -> %s\n", string(line))
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("read line unexpected error : %v\n", err)
			panic(err)
		}
		line_str := string(line)
		words := strings.Split(line_str, ",")
		if len(words) > columnLength {
			log.Fatalf("words length must be less than columnLength %s\n", string(line))
			return
		}

		for idx, w := range words {
			if len(out_words) > 0 && idx == 0 { //すでに前のループで単語を登録していたら
				out_words[len(out_words)-1] = out_words[len(out_words)-1] + "\n" + w
				continue
			}
			out_words = append(out_words, w)
		}
		fmt.Printf("len(out_words)=%d\n", len(out_words))
		if len(out_words) == columnLength { //同じ長さだったら、"で囲む
			for idx, w := range out_words {
				out_words[idx] = "\"" + w + "\""
			}
			out_line := strings.Join(out_words, ",")
			fmt.Printf("output -> %s\n", out_line)
			out_lines = append(out_lines, out_line)
			out_words = []string{}
		}
	}

	for _, l := range out_lines {
		writer.WriteString(l + "\n")
	}
	writer.Flush()
}
