package repository

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAdapter_GetValues(t *testing.T) {
	type args struct {
		readRange string
	}
	tests := []struct {
		name    string
		adapter *Repo
		args    args
		want    string
	}{
		{
			name:    "testGetValues",
			adapter: NewRepo(),
			args:    args{readRange: "Расходы/Доходы!A1:A4"},
			want:    "[[Наименование] [Невозвратный капитал] [Проценты] [Проценты]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%v", tt.adapter.getValues(tt.args.readRange)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
