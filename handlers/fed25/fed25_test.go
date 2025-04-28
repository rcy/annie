package fed25

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

func TestFetchSummary(t *testing.T) {
	sum, err := fetchSummary()
	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(sum)
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

	want := `prty seat lead vote pop%
LIB: 0000 0008 0000 00
GRN: 0005 0002 0000 00
CON: 0005 0001 0000 00
 BQ: 0000 0000 0000 00
NDP: 0000 0000 0000 00
PPC: 0000 0000 0000 00
OTH: 0000 0000 0000 00
Ridings: 0/343, Polls: 0/78953
`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("template output mismatch (-want +got):\n%s", diff)
	}

	// if got != want {
	// 	t.Errorf("want: %s\ngot: %s", want, got)
	// }
}
