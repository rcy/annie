package fed25

import (
	"reflect"
	"strconv"
	"testing"
)

func TestFetchRidingSummary(t *testing.T) {
	ridings, err := fetchRidings()
	if err != nil {
		t.Fatal(err)
	}

	if len(ridings) != 343 {
		t.Errorf("expected 343 ridings, got %d", len(ridings))
	}
}

func TestFindRidingsByName(t *testing.T) {
	for _, tc := range []struct {
		text  string
		want  []int
		error bool
	}{
		{
			text: "koot",
			want: []int{327, 304},
		},
		{
			text: "e",
			want: []int{298, 33, 34, 12, 23},
		},
		{
			text: "bogus",
			want: []int{},
		},
	} {
		t.Run(tc.text, func(t *testing.T) {
			ridings, err := findRidingsByName(tc.text)
			if err != nil {
				if tc.error {
					return
				}
				t.Fatal(err)
			}

			ridingIDs := make([]int, 0, len(ridings))
			for _, riding := range ridings {
				ridingIDs = append(ridingIDs, riding.RidingID)
			}

			if !reflect.DeepEqual(ridings, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, ridingIDs)
			}
		})
	}
}

func TestFindRidingsByNameSummaries(t *testing.T) {
	for _, tc := range []struct {
		text  string
		want  []string
		error bool
	}{
		{
			text: "koot",
			want: []string{"304: Columbia – Kootenay – Southern Rockies (CON 0 0) (GRN 0 0) (LIB 0 0) (NDP 0 0) (PPC 0 0) (IND 0 0)", "327: Similkameen – South Okanagan – West Kootenay (CON 0 0) (GRN 0 0) (LIB 0 0) (NDP 0 0) (PPC 0 0)"},
		},
	} {
		t.Run(tc.text, func(t *testing.T) {
			got, err := findRidingsByNameSummaries(tc.text)
			if err != nil {
				if tc.error {
					return
				}
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, got)
			}
		})
	}
}

func TestFindRidingByID(t *testing.T) {
	for _, tc := range []struct {
		id    int
		want  string
		error bool
	}{
		{
			id:   123,
			want: "Brampton South",
		},
		{
			id:    999,
			error: true,
		},
		{
			id:   340,
			want: "West Vancouver – Sunshine Coast – Sea to Sky Country",
		},
	} {
		t.Run(strconv.Itoa(tc.id), func(t *testing.T) {
			riding, err := findRidingByID(tc.id)
			if err != nil {
				if tc.error {
					return
				}
				t.Fatal(err)
			}

			if riding.RidingNameEng != tc.want {
				t.Errorf("expected %s, got %s", tc.want, riding.RidingNameEng)
			}
		})
	}
}

func TestFindRidingByIDSummary(t *testing.T) {
	for _, tc := range []struct {
		id    int
		want  string
		error bool
	}{
		{
			id:   123,
			want: "Brampton South",
		},
		{
			id:    999,
			error: true,
		},
		{
			id:   304,
			want: "West Vancouver – Sunshine Coast – Sea to Sky Country",
		},
	} {
		t.Run(strconv.Itoa(tc.id), func(t *testing.T) {
			got, err := findRidingByIDSummary(tc.id)
			if err != nil {
				if tc.error {
					return
				}
				t.Fatal(err)
			}

			if got != tc.want {
				t.Errorf("expected %s, got %s", tc.want, got)
			}
		})
	}
}
