package webconnectivity_test

import (
	"io"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ooni/probe-engine/experiment/urlgetter"
	"github.com/ooni/probe-engine/experiment/webconnectivity"
	"github.com/ooni/probe-engine/netx/archival"
	"github.com/ooni/probe-engine/netx/modelx"
)

func TestAnalyzeInvalidURL(t *testing.T) {
	out := webconnectivity.Analyze("\t\t\t", nil)
	if out.DNSConsistency != "" {
		t.Fatal("unexpected DNSConsistency")
	}
	if out.BodyLengthMatch != nil {
		t.Fatal("unexpected BodyLengthMatch")
	}
	if out.HeadersMatch != nil {
		t.Fatal("unexpected HeadersMatch")
	}
	if out.StatusCodeMatch != nil {
		t.Fatal("unexpected StatusCodeMatch")
	}
	if out.TitleMatch != nil {
		t.Fatal("unexpected TitleMatch")
	}
	if out.Accessible != nil {
		t.Fatal("unexpected Accessible")
	}
	if out.Blocking != nil {
		t.Fatal("unexpected Blocking")
	}
}

func TestDNSConsistency(t *testing.T) {
	measurementFailure := modelx.FailureDNSNXDOMAINError
	controlFailure := webconnectivity.ControlDNSNameError
	eofFailure := io.EOF.Error()
	type args struct {
		URL *url.URL
		tk  *webconnectivity.TestKeys
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{{
		name: "when the URL contains an IP address",
		args: args{
			URL: &url.URL{
				Host: "10.0.0.1",
			},
			tk: &webconnectivity.TestKeys{
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						Failure: &controlFailure,
					},
				},
			},
		},
		wantOut: "consistent",
	}, {
		name: "when the failures are not compatible",
		args: args{
			URL: &url.URL{
				Host: "www.kerneltrap.org",
			},
			tk: &webconnectivity.TestKeys{
				DNSExperimentFailure: &eofFailure,
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						Failure: &controlFailure,
					},
				},
			},
		},
		wantOut: "inconsistent",
	}, {
		name: "when the failures are compatible",
		args: args{
			URL: &url.URL{
				Host: "www.kerneltrap.org",
			},
			tk: &webconnectivity.TestKeys{
				DNSExperimentFailure: &measurementFailure,
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						Failure: &controlFailure,
					},
				},
			},
		},
		wantOut: "consistent",
	}, {
		name: "when the ASNs are equal",
		args: args{
			URL: &url.URL{
				Host: "fancy.dns",
			},
			tk: &webconnectivity.TestKeys{
				TestKeys: urlgetter.TestKeys{
					Queries: []archival.DNSQueryEntry{{
						Answers: []archival.DNSAnswerEntry{{
							ASN: 15169,
						}, {
							ASN: 13335,
						}},
					}},
				},
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						ASNs: []int64{13335, 15169},
					},
				},
			},
		},
		wantOut: "consistent",
	}, {
		name: "when the ASNs overlap",
		args: args{
			URL: &url.URL{
				Host: "fancy.dns",
			},
			tk: &webconnectivity.TestKeys{
				TestKeys: urlgetter.TestKeys{
					Queries: []archival.DNSQueryEntry{{
						Answers: []archival.DNSAnswerEntry{{
							ASN: 15169,
						}, {
							ASN: 13335,
						}},
					}},
				},
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						ASNs: []int64{13335, 13335},
					},
				},
			},
		},
		wantOut: "consistent",
	}, {
		name: "when the ASNs do not overlap",
		args: args{
			URL: &url.URL{
				Host: "fancy.dns",
			},
			tk: &webconnectivity.TestKeys{
				TestKeys: urlgetter.TestKeys{
					Queries: []archival.DNSQueryEntry{{
						Answers: []archival.DNSAnswerEntry{{
							ASN: 15169,
						}, {
							ASN: 15169,
						}},
					}},
				},
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						ASNs: []int64{13335, 13335},
					},
				},
			},
		},
		wantOut: "inconsistent",
	}, {
		name: "when ASNs lookup fails but IPs overlap",
		args: args{
			URL: &url.URL{
				Host: "fancy.dns",
			},
			tk: &webconnectivity.TestKeys{
				TestKeys: urlgetter.TestKeys{
					Queries: []archival.DNSQueryEntry{{
						Answers: []archival.DNSAnswerEntry{{
							IPv4: "8.8.8.8",
							ASN:  0,
						}, {
							IPv6: "2001:4860:4860::8844",
							ASN:  0,
						}},
					}},
				},
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						Addrs: []string{"8.8.4.4", "2001:4860:4860::8844"},
						ASNs:  []int64{0, 0},
					},
				},
			},
		},
		wantOut: "consistent",
	}, {
		name: "when ASNs lookup fails and IPs do not overlap",
		args: args{
			URL: &url.URL{
				Host: "fancy.dns",
			},
			tk: &webconnectivity.TestKeys{
				TestKeys: urlgetter.TestKeys{
					Queries: []archival.DNSQueryEntry{{
						Answers: []archival.DNSAnswerEntry{{
							IPv4: "8.8.8.8",
							ASN:  0,
						}, {
							IPv6: "2001:4860:4860::8844",
							ASN:  0,
						}},
					}},
				},
				Control: webconnectivity.ControlResponse{
					DNS: webconnectivity.ControlDNSResult{
						Addrs: []string{"8.8.4.4", "2001:4860:4860::8888"},
						ASNs:  []int64{0, 0},
					},
				},
			},
		},
		wantOut: "inconsistent",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := webconnectivity.DNSConsistency(tt.args.URL, tt.args.tk)
			if diff := cmp.Diff(tt.wantOut, gotOut); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
