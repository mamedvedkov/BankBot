package internals

import (
	"testing"
)

func Test_getRowByTgId(t *testing.T) {
	type args struct {
		adapter *Repo
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
				adapter: NewRepo(),
				tgId:    72597934,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := getRowByTgId(tt.args.adapter, tt.args.tgId); got != tt.want {
				t.Errorf("\ngetRowByTgId()\n%v,\nwant\n%v", got, tt.want)
			}
		})
	}
}

func Test_aboutMyPayment(t *testing.T) {
	type args struct {
		adapter *Repo
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
				adapter: NewRepo(),
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
				t.Errorf("\naboutMyPayment()\n%v,\nwant\n%v", got, tt.want)
			}
		})
	}
}

func Test_status(t *testing.T) {
	type args struct {
		adapter *Repo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test status",
			args: args{NewRepo()},
			want: "Капитал - 1 320 950₽\nЗанято - 898 337₽\nЗапас - 150 208₽\nАктив - 272 405₽",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := status(tt.args.adapter); got != tt.want {
				t.Errorf("\nstatus()\n%v,\nwant\n%v", got, tt.want)
			}
		})
	}
}

func Test_cardHolders(t *testing.T) {
	type args struct {
		repo *Repo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test card holders",
			args: args{NewRepo()},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cardHolders(tt.args.repo); got != tt.want {
				t.Errorf("cardHolders() =\n%v\n, want\n%v", got, tt.want)
			}
		})
	}
}

func Test_debts(t *testing.T) {
	type args struct {
		repo *Repo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test for debs",
			args: args{NewRepo()},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := debts(tt.args.repo); got != tt.want {
				t.Errorf("debts() =\n%v\n, want\n%v", got, tt.want)
			}
		})
	}
}

func Test_getRowByName(t *testing.T) {
	type args struct {
		repo *Repo
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "тест поиск строки по фамилии",
			args: args{
				repo: NewRepo(),
				name: "авиро",
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "тест поиск строки по имени",
			args: args{
				repo: NewRepo(),
				name: "григор",
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "тест поиск строки по имени, слишком много результатов",
			args: args{
				repo: NewRepo(),
				name: "алексан",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRowByName(tt.args.repo, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRowByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRowByName() got =\n%v\n, want\n%v", got, tt.want)
			}
		})
	}
}

func Test_searchByRow(t *testing.T) {
	type args struct {
		repo   *Repo
		rowNum int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "поиск по номеру строки",
			args: args{
				repo:   NewRepo(),
				rowNum: 3,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchByRow(tt.args.repo, tt.args.rowNum); got != tt.want {
				t.Errorf("searchByRow() =\n%v\n, want\n%v", got, tt.want)
			}
		})
	}
}
