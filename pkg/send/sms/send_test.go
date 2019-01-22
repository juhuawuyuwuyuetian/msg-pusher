/* ====================================================
#   Copyright (C)2019 All rights reserved.
#
#   Author        : domchan
#   Email         : 814172254@qq.com
#   File Name     : send.go
#   Created       : 2019/1/7 14:33
#   Last Modified : 2019/1/7 14:33
#   Describe      :
#
# ====================================================*/
package sms

import (
	"net/http"
	"testing"

	"uuabc.com/sendmsg/config"
	"uuabc.com/sendmsg/pkg/send"
)

var (
	ec  *config.Aliyun
	cli *Client
)

func init() {
	err := config.Init("../../../conf.yaml")
	if err != nil {
		panic(err)
	}
	ec = config.AliyunConf()
	cli = NewClient(Config{
		AccessKeyId:  ec.AccessKeyId,
		AccessSecret: ec.AccessSecret,
		GatewayURL:   ec.GatewayURL,
	})
}

func TestClient_Send(t *testing.T) {
	type fields struct {
		cfg    Config
		client *http.Client
	}
	type args struct {
		msg send.Message
		do  send.DoRes
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "send_case_1",
			fields: fields{
				cfg:    cli.cfg,
				client: cli.client,
			},
			args: args{
				msg: NewRequest("13423234442", "test", "te", "te", "123"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				cfg:    tt.fields.cfg,
				client: tt.fields.client,
			}
			if err := c.Send(tt.args.msg, tt.args.do); (err != nil) != tt.wantErr {
				t.Errorf("Client.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}