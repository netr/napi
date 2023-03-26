package napi

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"it works as expected and returns lower case letters", "HEY DOES THIS WORK", "hey_does_this_work"},
		{"it removes multiple white spaces and returns a single underscore", "HEY      DOES      THIS       WORK", "hey_does_this_work"},
		{"it removes contiguous underscores and replaces them with one", "hey_does_this____________work", "hey_does_this_work"},
		{"it removes periods", "talk about . this bro", "talk_about_this_bro"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.arg); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
