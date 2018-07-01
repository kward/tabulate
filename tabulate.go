package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kward/tabulate/tabulate"
)

var (
	columns = flag.Int("cols", 0, "Number of columns; 0=all.")

	ifs, ofs string
	render   string
	reset    bool

	comment  = flag.String("comment_prefix", "#", "Comment prefix.")
	comments = flag.Bool("comments", true, "Ignore comments.")
)

func flagInit(rs []tabulate.Renderer) {
	// Flag initialization.
	flag.StringVar(&ifs, "I", " ", "Input field separator.")
	flag.StringVar(&ofs, "O", " ", "Output field separator.")
	flag.StringVar(&render, "r", "plain", "Output renderer.")
	flag.BoolVar(&reset, "R", false, "Reset column width after each text block.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "Supported renderers:")
		for _, r := range rs {
			fmt.Fprintf(os.Stderr, "  %v\n", r.Type())
		}
	}

	flag.Parse()

	// Flag validation.
	if *columns < 0 {
		log.Fatalf("invalid number of columns: %v", *columns)
	}
}

func read(fh *os.File, data *[]string) error {
	s := bufio.NewScanner(fh)
	for s.Scan() {
		*data = append(*data, s.Text())
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("ERROR Reading file: %v", err)
	}
	return nil
}

func main() {
	var (
		err  error
		data []string
	)

	flagInit(tabulate.Renderers)

	renderers := map[string]tabulate.Renderer{}
	for _, r := range tabulate.Renderers {
		renderers[r.Type()] = r
	}

	// Open file.
	fh := os.Stdin
	if len(flag.Args()) > 0 {
		fh, err = os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close()
	}

	// Read file.
	err = read(fh, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Parse file.
	tbl, err := tabulate.NewTable()
	if err != nil {
		log.Fatal(err)
	}
	tbl.Split(data, ifs, *columns)

	// Render file.
	r, ok := renderers[render]
	if !ok {
		log.Fatalf("Invalid --render flag value %v.", r)
	}
	switch r.(type) {
	case *tabulate.PlainRenderer:
		r.(*tabulate.PlainRenderer).OFS = ofs
	}
	fmt.Print(r.Render(tbl))
}
