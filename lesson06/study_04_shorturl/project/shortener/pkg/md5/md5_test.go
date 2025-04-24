package md5

import "testing"

func TestSum(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Test hello", args: args{[]byte("hello")}, want: "5d41402abc4b2a76b9719d911017c592"},
		{name: "Test world", args: args{[]byte("world")}, want: "7d793037a0760186574b0282f2f435e7"},
		{name: "空字符串示例", args: args{[]byte("")}, want: "d41d8cd98f00b204e9800998ecf8427e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.data); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
