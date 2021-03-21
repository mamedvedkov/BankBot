package internals

import (
	"testing"
)

func Test_getRowByTgId(t *testing.T) {
	type args struct {
		adapter *Adapter
		tgId    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "testGetRowByTgId",
			args: args{
				adapter: NewAdapter(),
				tgId:    72597934,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := getRowByTgId(tt.args.adapter, tt.args.tgId); got != tt.want {
				t.Errorf("getRowByTgId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aboutMyPayment(t *testing.T) {
	type args struct {
		adapter *Adapter
		id      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "testAboutMyPayment",
			args: args{
				adapter: NewAdapter(),
				id:      69711013,
			},
			want: "Статистика по платежам\n" +
				"03.21 - 400 рублей\n" +
				"02.21 - 400 рублей\n" +
				"01.21 - 501 рублей\n" +
				"12.20 - 500 рублей\n" +
				"11.20 - 500 рублей\n" +
				"10.20 - 500 рублей\n" +
				"09.20 - 500 рублей\n" +
				"Ранее - 38436,52 рублей",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aboutMyPayment(tt.args.adapter, tt.args.id); got != tt.want {
				t.Errorf("aboutMyPayment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_status(t *testing.T) {
	type args struct {
		adapter *Adapter
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test status",
			args: args{NewAdapter()},
			want: "Капитал - 1 320 950₽\n" +
				"Занято - 898 337₽\n" +
				"Запас - 150 208₽\n" +
				"Актив - 272 405₽",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := status(tt.args.adapter); got != tt.want {
				t.Errorf("status() = %v, want %v", got, tt.want)
			}
		})
	}
}
