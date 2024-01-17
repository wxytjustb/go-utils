package storage

import (
	"testing"
)

var testMinioConfig = MinioConfig{
	Endpoints: "127.0.0.1:9000",
	AccessKey: "8wlGrILhHKEbmpliXSdk",
	SecretKey: "IqralHNa2kkFsEyL8vARHyZangtm0udMSEd53sd4",
	UseSSL:    false,
}

func TestMinioModuleInitialize(t *testing.T) {

	tests := []struct {
		name    string
		args    *MinioConfig
		wantErr bool
	}{
		{
			name:    "Test MinioModuleInitialize",
			args:    &testMinioConfig,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MinioModuleInitialize(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("MinioModuleInitialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMinioClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Test NewMinioClient",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMinioClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMinioClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("NewMinioClient() got = %v, want not nil", got)
			}
		})
	}
}

func TestMinioClient_Upload(t *testing.T) {
	type fields struct {
		client *MinioClient
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test MinioClient_Upload",
			fields: fields{
				client: GetSingleMinioClient(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//tt.fields.client.PutObject()
		})
	}
}
