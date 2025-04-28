package fed25

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFetchSummary(t *testing.T) {
	_, err := fetchSummary()
	if err != nil {
		t.Fatal(err)
	}

	//spew.Dump(sum)
}

func TestFormatSummary(t *testing.T) {
	sum, err := fetchSummary()
	if err != nil {
		t.Fatal(err)
	}

	sum.OverviewPartyDetails[1].Elected = 5
	sum.OverviewPartyDetails[1].Leading = 1

	sum.OverviewPartyDetails[2].Elected = 5
	sum.OverviewPartyDetails[2].Leading = 2

	sum.OverviewPartyDetails[3].Elected = 0
	sum.OverviewPartyDetails[3].Leading = 8

	got, err := formatSummary(sum)
	if err != nil {
		t.Fatal(err)
	}

	want := `2025 Federal General Election, Ridings: 0/343, Polls: 0/75944
  
LIB: 0 elected, 8 leading, 0 votes, 0% popular vote
GRN: 5 elected, 2 leading, 0 votes, 0% popular vote
CON: 5 elected, 1 leading, 0 votes, 0% popular vote
BQ: 0 elected, 0 leading, 0 votes, 0% popular vote
NDP: 0 elected, 0 leading, 0 votes, 0% popular vote
PPC: 0 elected, 0 leading, 0 votes, 0% popular vote
OTH: 0 elected, 0 leading, 0 votes, 0% popular vote
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("template output mismatch (-want +got):\n%s", diff)
	}
}
