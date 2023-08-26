package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddText(t *testing.T) {
	storage, err := CreateRamStorage()
	require.NoError(t, err)

	type args struct {
		user string
		text string
	}
	type want struct {
		elementsCount int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "add first value",
			args: args{user: "asd", text: "123"},
			want: want{elementsCount: 1},
		},
		{
			name: "add second value with same data",
			args: args{user: "asd", text: "123"},
			want: want{elementsCount: 2},
		},
		{
			name: "add third value",
			args: args{user: "asd", text: "234"},
			want: want{elementsCount: 3},
		},
		{
			name: "add value for another user",
			args: args{user: "qwe", text: "234"},
			want: want{elementsCount: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.AddText(tt.args.user, tt.args.text)
			texts, err := storage.GetTexts(tt.args.user)
			require.NoError(t, err)
			assert.Equal(t, tt.want.elementsCount, len(texts))
		})
	}
}

func TestStorage_GetTexts(t *testing.T) {
	storage, err := CreateRamStorage()
	require.NoError(t, err)

	type args struct {
		user string
		text string
	}
	type add struct {
		data []args
	}
	type wantData struct {
		user  string
		texts []string
		res   bool
	}
	type want struct {
		data []wantData
	}
	tests := []struct {
		name string
		add  add
		want want
	}{
		{
			name: "valid add first value",
			add: add{[]args{
				{
					user: "asf",
					text: "qwe",
				},
				{
					user: "asf",
					text: "qwe",
				},
				{
					user: "asf",
					text: "qwe",
				},
				{
					user: "asf",
					text: "zxc",
				},
			}},
			want: want{[]wantData{
				{
					user: "asf",
					texts: []string{
						"qwe",
						"qwe",
						"qwe",
						"zxc",
					},
					res: true,
				},
			},
			},
		},
		{
			name: "valid add additional second value",
			add: add{[]args{
				{
					user: "asfd",
					text: "qwe",
				},
			}},
			want: want{[]wantData{
				{
					user: "asf",
					texts: []string{
						"qwe",
						"qwe",
						"qwe",
						"zxc",
					},
					res: true,
				},
				{
					user: "asfd",
					texts: []string{
						"qwe",
						"qwe",
						"qwe",
						"zxc",
					},
					res: false,
				},
				{
					user: "asfd",
					texts: []string{
						"qwe",
					},
					res: true,
				},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, argData := range tt.add.data {
				storage.AddText(argData.user, argData.text)
			}

			for _, wantData := range tt.want.data {
				texts, err := storage.GetTexts(wantData.user)
				require.NoError(t, err)
				assert.Equal(t, reflect.DeepEqual(texts, wantData.texts), wantData.res)
			}
		})
	}
}
