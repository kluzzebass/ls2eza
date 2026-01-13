package dig2doggo

import (
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "simple query",
			args: []string{"example.com"},
			want: []string{"-q", "example.com"},
		},
		{
			name: "query with type",
			args: []string{"example.com", "MX"},
			want: []string{"-q", "example.com", "-t", "MX"},
		},
		{
			name: "query with nameserver",
			args: []string{"@8.8.8.8", "example.com"},
			want: []string{"-q", "example.com", "-n", "8.8.8.8"},
		},
		{
			name: "query with type and nameserver",
			args: []string{"@8.8.8.8", "example.com", "A"},
			want: []string{"-q", "example.com", "-t", "A", "-n", "8.8.8.8"},
		},
		{
			name: "short flag -q",
			args: []string{"-q", "example.com"},
			want: []string{"-q", "example.com"},
		},
		{
			name: "short flag -t",
			args: []string{"-t", "AAAA", "example.com"},
			want: []string{"-q", "example.com", "-t", "AAAA"},
		},
		{
			name: "short flag -t combined",
			args: []string{"-tAAAA", "example.com"},
			want: []string{"-q", "example.com", "-t", "AAAA"},
		},
		{
			name: "ipv4 only",
			args: []string{"-4", "example.com"},
			want: []string{"-4", "-q", "example.com"},
		},
		{
			name: "ipv6 only",
			args: []string{"-6", "example.com"},
			want: []string{"-6", "-q", "example.com"},
		},
		{
			name: "reverse lookup",
			args: []string{"-x", "8.8.8.8"},
			want: []string{"-x", "-q", "8.8.8.8"},
		},
		{
			name: "plus short option",
			args: []string{"+short", "example.com"},
			want: []string{"--short", "-q", "example.com"},
		},
		{
			name: "plus tcp option",
			args: []string{"+tcp", "example.com"},
			want: []string{"-n", "@tcp://", "-q", "example.com"},
		},
		{
			name: "plus dnssec option",
			args: []string{"+dnssec", "example.com"},
			want: []string{"--do", "-q", "example.com"},
		},
		{
			name: "plus recurse option",
			args: []string{"+recurse", "example.com"},
			want: []string{"--rd", "-q", "example.com"},
		},
		{
			name: "plus aa option",
			args: []string{"+aa", "example.com"},
			want: []string{"--aa", "-q", "example.com"},
		},
		{
			name: "plus ad option",
			args: []string{"+ad", "example.com"},
			want: []string{"--ad", "-q", "example.com"},
		},
		{
			name: "plus cd option",
			args: []string{"+cd", "example.com"},
			want: []string{"--cd", "-q", "example.com"},
		},
		{
			name: "plus nsid option",
			args: []string{"+nsid", "example.com"},
			want: []string{"--nsid", "-q", "example.com"},
		},
		{
			name: "plus cookie option",
			args: []string{"+cookie", "example.com"},
			want: []string{"--cookie", "-q", "example.com"},
		},
		{
			name: "plus padding option",
			args: []string{"+padding", "example.com"},
			want: []string{"--padding", "-q", "example.com"},
		},
		{
			name: "plus ede option",
			args: []string{"+ede", "example.com"},
			want: []string{"--ede", "-q", "example.com"},
		},
		{
			name: "plus search option",
			args: []string{"+search", "example.com"},
			want: []string{"--search", "-q", "example.com"},
		},
		{
			name: "plus timeout option",
			args: []string{"+timeout=5", "example.com"},
			want: []string{"--timeout", "5s", "-q", "example.com"},
		},
		{
			name: "plus ndots option",
			args: []string{"+ndots=2", "example.com"},
			want: []string{"--ndots", "2", "-q", "example.com"},
		},
		{
			name: "plus subnet option",
			args: []string{"+subnet=192.0.2.0/24", "example.com"},
			want: []string{"--ecs", "192.0.2.0/24", "-q", "example.com"},
		},
		{
			name: "query with class",
			args: []string{"example.com", "IN", "A"},
			want: []string{"-q", "example.com", "-t", "A", "-c", "IN"},
		},
		{
			name: "query with class flag",
			args: []string{"-c", "CH", "example.com"},
			want: []string{"-q", "example.com", "-c", "CH"},
		},
		{
			name: "complex query",
			args: []string{"-4", "+short", "+dnssec", "@1.1.1.1", "example.com", "MX"},
			want: []string{"-4", "--short", "--do", "-q", "example.com", "-t", "MX", "-n", "1.1.1.1"},
		},
		{
			name: "nameserver with protocol",
			args: []string{"@tcp://8.8.8.8", "example.com"},
			want: []string{"-q", "example.com", "-n", "@tcp://8.8.8.8"},
		},
		{
			name: "nameserver with https",
			args: []string{"@https://cloudflare-dns.com/dns-query", "example.com"},
			want: []string{"-q", "example.com", "-n", "@https://cloudflare-dns.com/dns-query"},
		},
		{
			name: "debug mode",
			args: []string{"-m", "example.com"},
			want: []string{"--debug", "-q", "example.com"},
		},
		{
			name: "multiple plus options",
			args: []string{"+short", "+dnssec", "+nsid", "example.com"},
			want: []string{"--short", "--do", "--nsid", "-q", "example.com"},
		},
		{
			name: "negated plus option",
			args: []string{"+norecurse", "example.com"},
			want: []string{"-q", "example.com"},
		},
		{
			name: "combined short flags",
			args: []string{"-4m", "example.com"},
			want: []string{"-4", "--debug", "-q", "example.com"},
		},
	}

	tr := &Translator{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tr.Translate(tt.args, "")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidQueryType(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{"A record", "A", true},
		{"AAAA record", "AAAA", true},
		{"MX record", "MX", true},
		{"lowercase a", "a", true},
		{"lowercase mx", "mx", true},
		{"invalid type", "INVALID", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidQueryType(tt.arg)
			if got != tt.want {
				t.Errorf("isValidQueryType(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}

func TestIsValidQueryClass(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{"IN class", "IN", true},
		{"CH class", "CH", true},
		{"HS class", "HS", true},
		{"CHAOS class", "CHAOS", true},
		{"HESIOD class", "HESIOD", true},
		{"lowercase in", "in", true},
		{"invalid class", "INVALID", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidQueryClass(tt.arg)
			if got != tt.want {
				t.Errorf("isValidQueryClass(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}
