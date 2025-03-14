/*
Copyright The Finder2D Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package finder2d

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DefaultMinMatchPercentage default minimun match percentage. Any match
// percentage lower than this number is not considered
const DefaultMinMatchPercentage float64 = 50.0

// MinDelta minimum difference between coordinates of group of matches to be
// considered the same image
const MinDelta = 1

var (
	unoMatch  = "\033[47;1m \033[0m" // Bright Blue
	ceroMatch = "\033[45;1m \033[0m" // Bright Black
)

// Default values for a one and a zero in a matrix
var (
	DefaultOne  = []byte(`+`)[0]
	DefaultZero = []byte(` `)[0]
)

// Match represents the coordinate the target matrix was found in the source
// matrix and the percentage match
type Match struct {
	X, Y       int
	Percentage float64
}

// Finder2D is the struct used to find a 2D pattern into a 2D source matrix
type Finder2D struct {
	Target     *Matrix
	Source     *Matrix
	one, zero  byte
	Matches    []Match
	Percentage float64
	Delta      int
}

func (m *Match) String() string {
	return fmt.Sprintf("(%d,%d,%f)", m.X, m.Y, m.Percentage)
}

// New create an empty Finder 2D
func New(one, zero byte, percentage float64, delta int) *Finder2D {
	if percentage == 0 {
		percentage = DefaultMinMatchPercentage
	}
	if delta == 0 {
		delta = MinDelta
	}
	// if both are `0`
	if one+zero == 0 {
		one = DefaultOne
		zero = DefaultZero
	}
	return &Finder2D{
		one:        one,
		zero:       zero,
		Percentage: percentage,
		Delta:      delta,
	}
}

// String returns the found matches in JSON format
func (f *Finder2D) String() string {
	output, _ := json.Marshal(f.Matches)
	return string(output)
}

// Stringf return the matches in the requested format
func (f *Finder2D) Stringf(format string) string {
	switch format {
	case "", "text", "matrix":
		return f.Matrix()
	default: // json
		return f.String()
	}
}

// Values returns the values for the off and on bits
func (f *Finder2D) Values() (byte, byte) {
	return f.zero, f.one
}

// IsAMatchPoint returns true if the coordinate is a match coordinate
func (f *Finder2D) IsAMatchPoint(x, y int) bool {
	for _, m := range f.Matches {
		// TODO: return true if the coordinate is in the area of the match
		if m.X == x && m.Y == y {
			return true
		}
	}
	return false
}

// IsInMatchArea return true if the coordinate is in the area of the match
func (f *Finder2D) IsInMatchArea(x, y int) bool {
	for _, m := range f.Matches {
		if (m.X <= x && x <= m.X+f.Target.maxX) && (m.Y <= y && y <= m.Y+f.Target.maxY) {
			return true
		}
	}
	return false
}

// Matrix return the matches in the matrix
func (f *Finder2D) Matrix() string {
	var b bytes.Buffer
	for y := 0; y < f.Source.maxY; y++ {
		for x := 0; x < f.Source.maxX; x++ {
			o := uno
			z := cero
			if f.IsInMatchArea(x, y) {
				o = unoMatch
				z = ceroMatch
			}
			switch f.Source.Content[y][x] {
			case 0:
				b.WriteString(z)
			case 1:
				b.WriteString(o)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// LoadSource loads the source from a reader replacing the cell value given in
// `one` for `1` and `zero` for `0`
func (f *Finder2D) LoadSource(r io.Reader) error {
	m, err := LoadMatrix(r, f.one, f.zero)
	if err != nil {
		return err
	}

	f.Source = m
	return nil
}

// LoadTarget loads the target from a reader replacing the cell value given in
// `one` for `1` and `zero` for `0`
func (f *Finder2D) LoadTarget(r io.Reader) error {
	m, err := LoadMatrix(r, f.one, f.zero)
	if err != nil {
		return err
	}

	f.Target = m
	return nil
}

// SearchSimple find the occurences of the target in the source and the percentage
// match in the simplest way which is to iterate thru the entire matrix searching
// for the pattern, storing the match when the match percentage is higher than
// the required
func (f *Finder2D) SearchSimple() error {
	if f.Source == nil {
		return fmt.Errorf("not set source matrix")
	}
	if f.Target == nil {
		return fmt.Errorf("not set target matrix")
	}
	if f.Percentage == 0 {
		return fmt.Errorf("percentage cannot be 0%%")
	}
	if f.Delta == 0 {
		return fmt.Errorf("delta cannot be 0")
	}
	maxX, maxY := f.Source.Size()
	width, height := f.Target.Size()

	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			sample := f.Source.Sample(x, y, width, height)
			if sample == nil {
				break
			}
			p, err := sample.Compare(f.Target)
			if err != nil {
				return err
			}
			if p >= f.Percentage {
				f.Matches = append(f.Matches, Match{
					X:          x,
					Y:          y,
					Percentage: p,
				})
			}
		}
	}

	f.Matches = reduceMatches(f.Matches, f.Delta)

	return nil
}

func around(m Match, ms []Match, d int) bool {
	for _, m1 := range ms {
		dx := m1.X - m.X
		dy := m1.Y - m.Y
		if (dx >= -d && dx <= d) && (dy >= -d && dy <= d) {
			return true
		}
	}
	return false
}

func bestMatch(matches []Match) Match {
	var higherP float64
	var bestMatch Match

	for _, m := range matches {
		if m.Percentage >= higherP {
			bestMatch = m
			higherP = m.Percentage
		}
	}

	return bestMatch
}

func groupMatchesNear(m Match, initialUniv []Match, delta int) (group []Match, universe []Match) {
	univ := initialUniv
	var mov int
	group = []Match{m}

	for {
		newUniv := []Match{}
		for _, mi := range univ {
			if around(mi, group, delta) {
				group = append(group, mi)
				mov++
			} else {
				newUniv = append(newUniv, mi)
			}
		}
		univ = newUniv
		if mov == 0 {
			break
		}
		mov = 0
	}

	return group, univ
}

func reduceMatches(matches []Match, delta int) []Match {
	retMatches := []Match{}
	if len(matches) == 0 {
		return retMatches
	}

	for {
		var matchGroup []Match

		m := matches[0]
		matches = matches[1:]

		matchGroup, matches = groupMatchesNear(m, matches, delta)

		bestM := bestMatch(matchGroup)
		retMatches = append(retMatches, bestM)

		if len(matches) == 0 {
			break
		}
	}
	return retMatches
}
