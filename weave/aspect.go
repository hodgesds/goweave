package weave

import (
	"io/ioutil"
	"os"
	"strings"
)

// aspect contains advice, pointcuts and any imports needed
type Aspect struct {
	advize   Advice
	pointkut Pointcut
	importz  []string
}

// grab_aspects looks for an aspect file for each file
// this seems lame and contrary to what we want...
// I'd go as far as to say that we want this to be cross-pkg w/in root?
//
// maybe the rule should be - aspects are valid for anything in a
// project root?
func (w *Weave) loadAspects() {

	fz := w.findAspects()
	for i := 0; i < len(fz); i++ {

		buf, err := ioutil.ReadFile(fz[i])
		if err != nil {
			w.flog.Println(err)
		}
		s := string(buf)

		w.parseAspectFile(s)
	}

	if len(w.aspects) == 0 {
		w.flog.Println("no weaves")
		os.Exit(1)
	}

}

// parseImports returns an array of imports for the corresponding advice
func (w *Weave) parseImports(body string) []string {
	impbrace := strings.Split(body, "imports (")

	if len(impbrace) > 1 {
		end := strings.Split(impbrace[1], ")")[0]
		t := strings.TrimSpace(end)
		return strings.Split(t, "\n")
	} else {
		return []string{}
	}
}

// containsBefore returns true if the body has before advice
func (w *Weave) containsBefore(body string) bool {
	if strings.Contains(body, "before: {") {
		return true
	} else {
		return false
	}
}

// containsAfter returns true if the body has after advice
func (w *Weave) containsAfter(body string) bool {
	if strings.Contains(body, "after: {") {
		return true
	} else {
		return false
	}
}

// rightBraceCnt returns the number of right braces in a string
func (w *Weave) rightBraceCnt(body string) int {
	return strings.Count(body, "}")
}

// parseAdvice returns advice about this aspect
func (w *Weave) parseAdvice(body string) Advice {
	advize := strings.Split(body, "advice:")[1]

	a4t := ""
	b4t := ""
	ar4t := ""

	bbrace := strings.Split(advize, "before: {")
	if len(bbrace) > 1 {
		// fixme
		if w.containsAfter(bbrace[1]) {
			b4 := strings.Split(bbrace[1], "}")[0]
			b4t = strings.TrimSpace(b4)
			// ...
		} else {
			cnt := w.rightBraceCnt(bbrace[1])
			// have at most 3 right braces
			// 3 - 3 = 0
			// 4 - 3 = 1
			b4 := strings.SplitAfter(bbrace[1], "}")
			rb := ""
			if cnt == 3 {
				rb = b4[0]
				rb = rb[:len(rb)-1]
			} else {
				for i := 0; i < cnt-3; i++ {
					rb += strings.TrimSpace(b4[i])
				}
			}
			b4t = strings.TrimSpace(rb)
		}
	}

	abrace := strings.Split(advize, "after: {")
	if len(abrace) > 1 {
		a4 := strings.Split(abrace[1], "}")[0]
		a4t = strings.TrimSpace(a4)
	}

	arbrace := strings.Split(advize, "around: {")
	if len(arbrace) > 1 {
		ar4 := strings.Split(arbrace[1], "}")[0]
		ar4t = strings.TrimSpace(ar4)
	}

	return Advice{
		before: b4t,
		after:  a4t,
		around: ar4t,
	}

}

// parseAspectFile loads an individual file containing aspects
// there be tigers here
// subject to immediate and dramatic change
func (w *Weave) parseAspectFile(body string) {
	results := []Aspect{}

	aspects := strings.Split(body, "aspect {")

	for i := 0; i < len(aspects); i++ {

		aspect := aspects[i]
		azpect := Aspect{}

		if strings.TrimSpace(aspect) != "" {

			pk, err := w.parsePointCut(aspect)
			if err != nil {
				w.flog.Println(err.Error() + aspect)
				os.Exit(1)
			} else {
				azpect.pointkut = pk
			}

			azpect.importz = w.parseImports(aspect)

			azpect.advize = w.parseAdvice(aspect)

			results = append(results, azpect)
		}

	}

	w.aspects = results

}
