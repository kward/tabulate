package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kward/tabulate/render"
	"github.com/kward/tabulate/table"
)

var (
	columns        int
	ifs, ofs       string
	renderer       string
	enableComments bool
	commentPrefix  string
	sectionReset   bool
)

func flagInit(rs []render.Renderer) {
	// Flag initialization.
	flag.Int("cols", 0, "Number of columns; 0=all.")

	flag.StringVar(&ifs, "I", " ", "Input field separator.")
	flag.StringVar(&ofs, "O", " ", "Output field separator.")
	flag.StringVar(&renderer, "r", "plain", "Output renderer.")

	flag.BoolVar(&enableComments, "enable_comments", true, "Enable comments.")
	flag.StringVar(&commentPrefix, "comment_prefix", "#", "Comment prefix.")

	flag.BoolVar(&sectionReset, "R", false, "Reset column widths after each section.")

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
	if columns < 0 {
		log.Fatalf("invalid number of columns: %v", columns)
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

	flagInit(render.Renderers)

	renderers := map[string]render.Renderer{}
	for _, r := range render.Renderers {
		renderers[r.Type()] = r
	}
	_, ok := renderers[renderer]
	if !ok {
		log.Fatalf("Invalid --render flag value %v.", renderer)
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
	n := columns
	if n == 0 {
		n = -1
	}
	tbl, err := table.Split(data, ifs, n,
		table.CommentPrefix(commentPrefix),
		table.EnableComments(enableComments),
		table.SectionReset(sectionReset),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Render file.
	r := renderers[renderer]
	switch r.(type) {
	case *render.PlainRenderer:
		r.(*render.PlainRenderer).SetOFS(ofs)
	}
	fmt.Print(r.Render(tbl))
}
