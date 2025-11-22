package dishwasher

import (
	"reflect"
	"testing"
)

func TestDishwasher_WithPath(t *testing.T) {
	type args struct {
		Path          string
		DimensionSize DishSize
	}

	tests := []struct {
		name    string
		args    args
		want    Dishwasher
		UsePath bool
		UseSize bool
	}{
		{
			name: "with path and size",
			args: args{
				Path:          "/test.db",
				DimensionSize: SizeMid,
			},
			want: Dishwasher{
				Path:          "/test.db",
				DimensionSize: SizeMid,
			},
			UsePath: true,
			UseSize: true,
		},
		{
			name: "use default path",
			args: args{
				DimensionSize: SizeMid,
			},
			want: Dishwasher{
				Path:          "/dishwasher.db",
				DimensionSize: SizeMid,
			},
			UseSize: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := New()

			if tt.UsePath {
				value = value.WithPath(tt.args.Path)
			}

			if tt.UseSize {
				value = value.WithDimensionSize(tt.args.DimensionSize)
			}

			if !reflect.DeepEqual(value, tt.want) {
				t.Errorf("WithPath() = %v, want %v", value, tt.want)
			}
		})
	}
}
