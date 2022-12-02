package template

import (
	"reflect"
	"strconv"
	"testing"
)

func TestTemplate(t *testing.T) {
	type args struct {
		templateContent string
		forwardBody     []byte
	}
	tests := []struct {
		args args
		want []string
	}{
		{args: args{
			templateContent: `function convert(origin){alert=origin.alerts[0];return{target:["https://api.day.app/asd/","https://api.day.app/123/"],payload:{title:"["+alert.status+"] "+alert.labels.alertname,body:"",}}};`,
			forwardBody:     []byte(`{"alerts":[{"status":"resolved","labels":{"alertname":"325i alert"}}]}`),
		}, want: []string{"https://api.day.app/asd/", "https://api.day.app/123/"}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			template, err := New(tt.args.templateContent, tt.args.forwardBody)
			if err != nil {
				t.Error(err)
			}
			r, err := template.RunJs()
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(r.Targets, tt.want) {
				t.Errorf("got %v, want %v", r.Targets, tt.want)
			}
		})
	}
}
