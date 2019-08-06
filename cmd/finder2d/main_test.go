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

package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func testSetup() {
	os.Setenv(envPrefix+"_NAME", "Johandry")
	os.Setenv(envPrefix+"_CITY", "San Diego")
	os.Setenv(envPrefix+"_AGE", "42")
	os.Setenv(envPrefix+"_SCORE", "4.5")
}

func Test_getEnv(t *testing.T) {
	testSetup()

	type args struct {
		name     string
		defValue string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"lowercase", args{"name", ""}, "Johandry"},
		{"any case", args{"CiTy", ""}, "San Diego"},
		{"default value", args{"LASTNAME", "Amador"}, "Amador"},
		{"env value", args{"NAME", "Pepe"}, "Johandry"},
		{"integer as string", args{"AGE", "foo"}, "42"},
		{"float as string", args{"SCORE", "bar"}, "4.5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.name, tt.args.defValue); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvFloat(t *testing.T) {
	testSetup()
	type args struct {
		name   string
		defVal float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"float", args{"SCORE", 0}, 4.5},
		{"default value", args{"RATE", 9.99}, 9.99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnvFloat(tt.args.name, tt.args.defVal); got != tt.want {
				t.Errorf("getEnvFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvInt(t *testing.T) {
	testSetup()
	type args struct {
		name   string
		defVal int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"integer", args{"AGE", 23}, 42},
		{"default value", args{"BOOKS", 123}, 123},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnvInt(tt.args.name, tt.args.defVal); got != tt.want {
				t.Errorf("getEnvInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_config_Read(t *testing.T) {
	type fields struct {
		sourceFileName string
		targetFileName string
		zero           string
		one            string
		percentage     float64
		delta          int
		output         string
		port           string
	}
	tests := []struct {
		name   string
		fields fields
		envFields map[string]string
		flagFields map[string]string
		want fields
	}{
		{
			"just one test to fit all", 
			fields{zero: " ", one: "*", percentage: 65.9, delta: 3}, 
			map[string]string{"PERCENTAGE": "70.3", "ON": "@"}, 
			map[string]string{"p": "80.5", "d": "2"}, 
			fields{zero: " ", one: "@", percentage: 80.5, delta: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for name, value := range tt.envFields {
				os.Setenv(envPrefix + "_" + name, value)
			}
			
			c := &config{
				sourceFileName: tt.fields.sourceFileName,
				targetFileName: tt.fields.targetFileName,
				zero:           tt.fields.zero,
				one:            tt.fields.one,
				percentage:     tt.fields.percentage,
				delta:          tt.fields.delta,
				output:         tt.fields.output,
				port:           tt.fields.port,
			}
			c.Init()
			for name, value := range tt.flagFields {
				if err := flag.Set(name, value); err != nil {
						t.Errorf("read() failed to set the flag --%s %s. %s", name, value, err)
				}
			}
			c.Read()
			wantC := &config{
				sourceFileName: tt.want.sourceFileName,
				targetFileName: tt.want.targetFileName,
				zero:           tt.want.zero,
				one:            tt.want.one,
				percentage:     tt.want.percentage,
				delta:          tt.want.delta,
				output:         tt.want.output,
				port:           tt.want.port,
			}
			if !reflect.DeepEqual(c, wantC) {
				t.Errorf("read() = %+v, want %+v", c, wantC)
			}
		})
	}
}
