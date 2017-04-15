package main

import (
	"flag"
	"fmt"
	"go/build"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/kisielk/gotool"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of %s:
	glyphcheck [flags]               # runs on package in current directory
	glyphcheck [flags] [packages]    # specified as Go import paths
	glyphcheck [flags] [directories] # where a '/...' suffix includes all sub-directories
	glyphcheck [flags] [files]       # all must belong to a single package
Flags:
`, os.Args[0])
		flag.PrintDefaults()
	}
	scanComments := flag.Bool("comments", false, "also scan comments")
	flag.Parse()

	var lines []string
	for _, file := range listFiles(flag.Args()) {
		lines = append(lines, scanFile(file, *scanComments)...)
	}
	if len(lines) > 0 {
		log.Fatalln(strings.Join(lines, "\n"))
	}
}

// listFiles parses args as a set of Go files or import paths and returns the
// set of all Go files it references.
func listFiles(args []string) (files []string) {
	if len(args) > 0 && strings.HasSuffix(args[0], ".go") {
		// assume ad-hoc package
		return args
	}
	for _, path := range gotool.ImportPaths(args) {
		buildPkg, err := build.Import(path, ".", 0)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range buildPkg.GoFiles {
			files = append(files, filepath.Join(buildPkg.Dir, file))
		}
	}
	return
}

// scanFile scans file and returns formatted error strings for each line in
// file that contains a homoglyph.
func scanFile(file string, scanComments bool) (lines []string) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	tokFile := fset.AddFile("", fset.Base(), len(src))
	var scanMode scanner.Mode
	if scanComments {
		scanMode = scanner.ScanComments
	}
	s.Init(tokFile, src, func(pos token.Position, msg string) {
		// exit if file cannot be scanned
		log.Fatalf("%v:%v: %v", pos.Filename, pos.Line, msg)
	}, scanMode)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			return
		}
		for _, c := range lit {
			if c >= utf8.RuneSelf && isHomoglyph(c) {
				lines = append(lines, fmt.Sprintf("%s:%v: %s (%#U)", file, fset.Position(pos), lit, c))
				break
			}
		}
	}
}

var homoglyphs map[rune]struct{}

func init() {
	homoglyphs = make(map[rune]struct{})
	for _, r := range `ᅟᅠ         　ㅤǃ！״″＂＄％＆＇﹝（﹞）⁎＊＋‚，‐𐆑－٠۔܁܂․‧。．｡⁄∕╱⫻⫽／ﾉΟοОоՕ𐒆ＯｏΟοОоՕ𐒆Ｏｏا１２３４５６𐒇７Ց８９։܃܄∶꞉：;；‹＜𐆐＝›＞？＠［＼］＾＿｀ÀÁÂÃÄÅàáâãäåɑΑαаᎪＡａßʙΒβВЬᏴᛒＢｂϲϹСсᏟⅭⅽ𐒨ＣｃĎďĐđԁժᎠⅮⅾＤｄÈÉÊËéêëĒēĔĕĖėĘĚěΕЕеᎬＥｅϜＦｆɡɢԌնᏀＧｇʜΗНһᎻＨｈɩΙІіاᎥᛁⅠⅰ𐒃ＩｉϳЈјյᎫＪｊΚκКᏦᛕKＫｋʟιاᏞⅬⅼＬｌΜϺМᎷᛖⅯⅿＭｍɴΝＮｎΟοОоՕ𐒆ＯｏΟοОоՕ𐒆ＯｏΡρРрᏢＰｐႭႳＱｑʀԻᏒᚱＲｒЅѕՏႽᏚ𐒖ＳｓΤτТᎢＴｔμυԱՍ⋃ＵｕνѴѵᏙⅤⅴＶｖѡᎳＷｗΧχХхⅩⅹＸｘʏΥγуҮＹｙΖᏃＺｚ｛ǀا｜｝⁓～ӧӒӦ` {
		homoglyphs[r] = struct{}{}
	}
}

func isHomoglyph(r rune) bool {
	_, ok := homoglyphs[r]
	return ok
}
