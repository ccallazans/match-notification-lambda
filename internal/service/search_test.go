package service

import (
	"reflect"
	"testing"
)

func TestGetUsersByTopic(t *testing.T) {
	type args struct {
		topics []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUsersByTopic(tt.args.topics)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUsersByTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUsersByTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
