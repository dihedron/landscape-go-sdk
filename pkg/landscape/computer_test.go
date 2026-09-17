package landscape

import "testing"

func TestComputerListOptionsQueryParams(t *testing.T) {
	got := (ComputerListOptions{
		Query: ComputerQuery{
			Keywords:                  []string{"appserver", "OR"},
			Tags:                      []string{"server"},
			Distributions:             []string{"jammy"},
			IDs:                       []int{1, 2},
			Profiles:                  []string{"wsl:1:compliant"},
			ContractExpiresWithinDays: []int{30},
			HasProManagement:          []bool{true},
		},
		Limit: 20,
	}).queryParams()

	want := "appserver OR tag:server distribution:jammy id:1 OR id:2 profile:wsl:1:compliant contract-expires-within-days:30 has-pro-management:true"
	if got["query"] != want {
		t.Fatalf("query = %q, want %q", got["query"], want)
	}
	if got["limit"] != "20" {
		t.Fatalf("limit = %q, want %q", got["limit"], "20")
	}
}

func TestComputerListOptionsOmitsEmptyQuery(t *testing.T) {
	if _, ok := (ComputerListOptions{}).queryParams()["query"]; ok {
		t.Fatal("empty query should be omitted")
	}
}

func TestComputerListOptionsDocumentedQueryExamples(t *testing.T) {
	tests := []struct {
		name  string
		query ComputerQuery
		want  string
	}{
		{
			name: "profile status",
			query: ComputerQuery{
				Profiles: []string{"wsl:1:compliant"},
			},
			want: "profile:wsl:1:compliant",
		},
		{
			name: "multiple IDs",
			query: ComputerQuery{
				IDs: []int{1, 2},
			},
			want: "id:1 OR id:2",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.query.String(); got != test.want {
				t.Fatalf("query = %q, want %q", got, test.want)
			}
		})
	}
}

func TestComputerListOptionsDocumentedSelectorPrefixes(t *testing.T) {
	query := ComputerQuery{
		Tags:                      []string{"server"},
		Distributions:             []string{"jammy"},
		Hostnames:                 []string{"appserv1"},
		Titles:                    []string{"Application Server 1"},
		Alerts:                    []string{"security-upgrades"},
		AccessGroups:              []string{"server"},
		IDs:                       []int{1},
		MACs:                      []string{"00:25:00:49:ef:80"},
		IPs:                       []string{"192.168.1.102"},
		Searches:                  []string{"production"},
		Needs:                     []string{"reboot", "license"},
		LicenseIDs:                []string{"123"},
		Annotations:               []string{"location:datacenter"},
		Profiles:                  []string{"security:7:pass", "usg:8:fail", "wsl:1:compliant"},
		ReleaseUpgrade:            []string{"available"},
		LicenseTypes:              []string{"pro"},
		Contracts:                 []string{"contract-123"},
		ContractExpiresWithinDays: []int{30},
		LicenseExpiresWithinDays:  []int{15},
		HasProManagement:          []bool{true, false},
	}

	want := "tag:server distribution:jammy hostname:appserv1 title:Application Server 1 alert:security-upgrades access-group:server id:1 mac:00:25:00:49:ef:80 ip:192.168.1.102 search:production needs:reboot OR needs:license license-id:123 annotation:location:datacenter profile:security:7:pass OR profile:usg:8:fail OR profile:wsl:1:compliant release-upgrade:available license-type:pro contract:contract-123 contract-expires-within-days:30 license-expires-within-days:15 has-pro-management:true OR has-pro-management:false"
	if got := query.String(); got != want {
		t.Fatalf("query = %q, want %q", got, want)
	}
}
