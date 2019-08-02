package finder2d

import (
	"reflect"
	"testing"
)

func Test_around(t *testing.T) {
	ms := []Match{
		Match{44, 20, 51.555556},
		Match{45, 21, 57.333333},
		Match{46, 22, 51.111111},
	}
	tests := []struct {
		name string
		m    Match
		ms   []Match
		d    int
		want bool
	}{
		{"empty list", Match{1, 2, 3.4}, []Match{}, 1, false},
		{"near X (1)", Match{43, 21, 0.0}, ms, 1, true},
		{"near X,Y (1)", Match{43, 19, 0.0}, ms, 1, true},
		{"near X (2)", Match{47, 23, 0.0}, ms, 1, true},
		{"near XY (2)", Match{47, 21, 0.0}, ms, 1, true},
		{"far X (1)", Match{42, 21, 0.0}, ms, 1, false},
		{"near X, far Y (1)", Match{43, 22, 0.0}, ms, 1, false},
		{"far X (2)", Match{48, 23, 0.0}, ms, 1, false},
		{"near X far Y (2)", Match{47, 20, 0.0}, ms, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := around(tt.m, tt.ms, tt.d); got != tt.want {
				t.Errorf("around() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bestMatch(t *testing.T) {
	tests := []struct {
		name    string
		matches []Match
		want    Match
	}{
		{"empty", []Match{}, Match{}},
		{"one", []Match{
			Match{74, 0, 52.888889},
			Match{75, 0, 54.222222},
			Match{76, 0, 57.777778},
			Match{77, 0, 56.444444},
			Match{78, 0, 55.555556},
			Match{79, 0, 72.000000},
			Match{80, 0, 99.111111},
			Match{81, 0, 71.555556},
			Match{82, 0, 55.555556},
			Match{83, 0, 58.222222},
			Match{84, 0, 60.000000},
			Match{85, 0, 53.777778},
			Match{76, 1, 50.222222},
			Match{77, 1, 51.111111},
		}, Match{80, 0, 99.111111}},
		{"multiple", []Match{
			Match{74, 0, 52.888889},
			Match{75, 0, 54.222222},
			Match{76, 0, 99.111111},
			Match{77, 0, 56.444444},
			Match{78, 0, 55.555556},
			Match{79, 0, 72.000000},
			Match{80, 0, 99.111111},
			Match{81, 0, 71.555556},
			Match{82, 0, 55.555556},
			Match{83, 0, 99.111111},
			Match{84, 0, 60.000000},
			Match{85, 0, 53.777778},
			Match{76, 1, 50.222222},
			Match{77, 1, 51.111111},
		}, Match{83, 0, 99.111111}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bestMatch(tt.matches); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bestMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
