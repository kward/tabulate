package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kward/tabulate/tabulate"
)

const (
	defaultIFS    = " "
	defaultOFS    = " "
	defaultRender = "plain"
)

var (
	columns = flag.Int("cols", 0, "Number of columns; 0=all.")

	ifs, ofs string
	render   string

	comment  = flag.String("comment_prefix", "#", "Comment prefix.")
	comments = flag.Bool("comments", true, "Ignore comments.")
)

func flagInit(renderers []tabulate.Renderer) {
	// Flag initialization.
	flag.StringVar(&ifs, "I", defaultIFS, "Input field separator.")
	flag.StringVar(&ofs, "O", defaultOFS, "Output field separator.")
	flag.StringVar(&render, "r", defaultRender, "Output renderer.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "Supported renderers:")
		for _, r := range renderers {
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

	rmap := map[string]tabulate.Renderer{}
	for _, r := range tabulate.Renderers {
		rmap[r.Type()] = r
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
	t := tabulate.NewTable(tabulate.NewTableConfig())
	t.Split(data, ifs, *columns)

	// Render file.
	renderer, ok := rmap[render]
	if !ok {
		log.Fatalf("Invalid --render flag value %v.", render)
	}
	switch renderer.(type) {
	case *tabulate.PlainRenderer:
		renderer.(*tabulate.PlainRenderer).OFS = ofs
	}
	fmt.Print(renderer.Render(&t))
}
