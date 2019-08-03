package finder2d

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFinder2D_LoadMatrix(t *testing.T) {
	type fields struct {
		Content [][]int
		maxX    int
		maxY    int
	}
	type args struct {
		content []byte
		one     byte
		zero    byte
	}
	tests := []struct {
		name    string
		args    args
		want    fields
		wantErr bool
	}{
		{"empty matrix", args{[]byte{}, []byte(`1`)[0], []byte(`0`)[0]}, fields{}, false},
		{"2x2", args{[]byte("11\n00"), []byte(`1`)[0], []byte(`0`)[0]}, fields{Content: [][]int{{1, 1}, {0, 0}}, maxX: 2, maxY: 2}, false},
		{"3x3", args{[]byte("101\n010\n110"), []byte(`1`)[0], []byte(`0`)[0]}, fields{Content: [][]int{{1, 0, 1}, {0, 1, 0}, {1, 1, 0}}, maxX: 3, maxY: 3}, false},
		{"fillable", args{[]byte("101\n00\n1\n"), []byte(`1`)[0], []byte(`0`)[0]}, fields{Content: [][]int{{1, 0, 1}, {0, 0, 0}, {1, 0, 0}}, maxX: 3, maxY: 3}, false},
		{"empty lines", args{[]byte("11\n00\n\n\n"), []byte(`1`)[0], []byte(`0`)[0]}, fields{Content: [][]int{{1, 1}, {0, 0}, {0, 0}, {0, 0}}, maxX: 2, maxY: 4}, false},
		{"irregular", args{[]byte("101\n0000\n1\n"), []byte(`1`)[0], []byte(`0`)[0]}, fields{}, true},
		{"invalid value", args{[]byte("101\n0x0\n1yz\n"), []byte(`1`)[0], []byte(`0`)[0]}, fields{}, true},
		{"4x3", args{[]byte("1010\n0101\n1100\n"), []byte(`1`)[0], []byte(`0`)[0]}, fields{Content: [][]int{{1, 0, 1, 0}, {0, 1, 0, 1}, {1, 1, 0, 0}}, maxX: 4, maxY: 3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Matrix{}
			r := bytes.NewBuffer(tt.args.content)
			if err := m.Load(r, tt.args.one, tt.args.zero); (err != nil) != tt.wantErr {
				t.Errorf("Finder2D.Matrix.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotX, gotY := m.Size(); gotX != tt.want.maxX || gotY != tt.want.maxY {
				t.Errorf("Finder2D.Matrix.Size() = (%d,%d), want size (%d,%d)", gotX, gotY, tt.want.maxX, tt.want.maxY)
			}
			if !reflect.DeepEqual(m.Content, tt.want.Content) {
				t.Errorf("Finder2D.Matrix.Load() content = %v, want %v", m.Content, tt.want.Content)
			}
		})
	}
}

func TestMatrix_Sample(t *testing.T) {
	type args struct {
		x int
		y int
		w int
		h int
	}
	tests := []struct {
		name string
		args args
		want *Matrix
	}{
		{"empty", args{}, &Matrix{}},
		{"larger sample", args{15, 15, 10, 10}, nil},
		{"2x2", args{15, 15, 2, 2}, &Matrix{Content: [][]int{{1, 1}, {1, 1}}, maxX: 2, maxY: 2}},
		{"3x3", args{1, 2, 3, 3}, &Matrix{Content: [][]int{{1, 0, 1}, {1, 1, 1}, {0, 1, 1}}, maxX: 3, maxY: 3}},
		{"4x4", args{0, 16, 4, 4}, &Matrix{Content: [][]int{{1, 1, 0, 0}, {1, 0, 1, 0}, {1, 0, 0, 1}, {1, 0, 0, 0}}, maxX: 4, maxY: 4}},
	}
	m, err := LoadMatrix(bytes.NewBufferString(string(testMatrixData[0])), testMatrixOne, testMatrixZero)
	if err != nil {
		t.Errorf("Matrix.Sample() failed to load the matrix. %s", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Sample(tt.args.x, tt.args.y, tt.args.w, tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Matrix.Sample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Compare(t *testing.T) {
	tests := []struct {
		name    string
		m       []byte
		m1      []byte
		want    float64
		wantErr bool
	}{
		{"same", testMatrixData[0], testMatrixData[0], 100.0, false},
		{"diff", testMatrixData[0], testMatrixData[1], 94.0, false},
		{"diff size", testMatrixData[0], testMatrixData[2], 0, true},
		{"lot of spaces", testMatrixData[2], testMatrixData[3], 40.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := LoadMatrix(bytes.NewBufferString(string(tt.m)), testMatrixOne, testMatrixZero)
			if err != nil {
				t.Errorf("Matrix.Compare() failed to load the source matrix. %s", err)
			}
			m1, err := LoadMatrix(bytes.NewBufferString(string(tt.m1)), testMatrixOne, testMatrixZero)
			if err != nil {
				t.Errorf("Matrix.Compare() failed to load the target matrix. %s", err)
			}
			// t.Logf("Compare Matrix:\n%s == \n%s\n", m, m1)
			got, err := m.Compare(m1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Matrix.Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Matrix.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testMatrixOne = DefaultOne
var testMatrixZero = DefaultZero
var testMatrixData = [][]byte{
	[]byte(`++++++++++++++ +++++
+++++++++++++++++++ 
++ +++++++ +++++++++
++++++++++++++ + + +
+ +++++++++ +++++++ 
++++ +++++++++ +++++
+++++++++++++ ++++++
++ +++++++++++++ ++ 
+++++ +  ++++ ++++  
++++++++++++++++ +++
++++++ ++++ +++++++ 
++++++++++++++++++ +
+++++ ++++++++++  ++
++ + +++++ ++ ++ ++ 
++++++++++ +++++ +++
++++++++++++  +++++ 
++   +++++++ +++++++
+ +  ++++++ ++++ +++
+  + + ++++ ++ + +++
+   +  ++ ++++++++++`),

	[]byte(`++++++++++++++ +++++
+++++++++++++++++++ 
++ +++++++ +++++++++
++++++++++++++ + + +
+ +++++++++ +++++++ 
++++ +++++++++ +++++
++++++  +++++ ++++++
++ +++++++++++++ ++ 
+++++ +  ++++ ++++  
++++++++++++++++ +++
++++++ ++++ +++++++ 
++++++++++++++++++ +
+++++ ++++++++++  ++
++ + +++++ ++ ++ ++ 
++++++++++ +++++ +++
++++++++++++  +++++ 
+++++++++++ ++++++++
+ +++++++++ ++++ +++
+   +++ ++++ ++ + ++
+  ++ ++ ++++++++++ `),

	[]byte(`++ ++
+++++
+  ++
++++ 
+ +++`),

	[]byte(`    +
+    
+    
+    
+    `),
}
