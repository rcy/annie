package fed25

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			want: []int{304, 327},
		},
		{
			text: "e",
			want: []int{298, 33, 34, 12, 23, 35, 261, 36, 112, 37, 113, 114, 115, 262, 247, 116, 117, 38, 39, 40, 24, 42, 43, 44, 263, 118, 119, 121, 122, 124, 46, 47, 126, 128, 299, 300, 41, 264, 265, 267, 268, 270, 271, 272, 274, 129, 13, 2, 301, 130, 248, 3, 14, 48, 9, 131, 50, 302, 234, 49, 303, 304, 51, 306, 15, 53, 52, 16, 132, 308, 249, 133, 134, 54, 135, 275, 276, 277, 278, 279, 280, 281, 282, 283, 136, 10, 137, 235, 309, 138, 139, 141, 140, 310, 285, 25, 56, 57, 286, 143, 145, 18, 146, 147, 149, 150, 58, 59, 60, 151, 152, 61, 62, 311, 154, 313, 155, 156, 158, 160, 159, 63, 64, 67, 65, 287, 161, 314, 68, 69, 70, 288, 162, 71, 289, 163, 165, 164, 5, 73, 74, 75, 76, 72, 27, 11, 77, 166, 168, 290, 169, 170, 79, 28, 171, 172, 173, 174, 176, 29, 81, 83, 250, 78, 177, 179, 317, 178, 180, 182, 318, 319, 184, 342, 84, 185, 186, 320, 189, 192, 191, 85, 86, 292, 195, 196, 197, 87, 88, 321, 237, 90, 198, 251, 323, 238, 91, 294, 252, 253, 254, 92, 324, 325, 95, 96, 97, 20, 30, 98, 99, 100, 101, 102, 200, 257, 256, 201, 203, 206, 204, 241, 103, 104, 295, 209, 208, 327, 328, 258, 21, 329, 296, 240, 211, 6, 212, 214, 330, 331, 259, 22, 215, 7, 105, 217, 218, 106, 32, 219, 107, 222, 332, 333, 334, 335, 336, 337, 108, 223, 338, 109, 224, 225, 340, 227, 229, 228, 242, 243, 244, 245, 246, 297, 230, 232, 260},
		},
		{
			text: "bogus",
			want: []int{},
		},
		{
			text: "quadra",
			want: []int{337},
		},
		{
			text: "Vancouver Quadra",
			want: []int{337},
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

			if !reflect.DeepEqual(ridingIDs, tc.want) {
				t.Errorf("expected %#v, got %#v", tc.want, ridingIDs)
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
			want: `Brampton South`,
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

			got := riding.RidingNameEng

			if got != tc.want {
				t.Errorf("expected %s, got %s", tc.want, got)
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
			id: 123,
			want: `[123] Brampton South, Ontario, 67956 votes, 0/169 polls reporting, 2025-04-28 13:54:12
[123] 

[123] CON Sukhdeep Kang - 0 votes, 0 margin 

[123] LIB Sonia Sidhu [incumbent] - 0 votes, 0 margin 

[123] NDP Rajni Sharma - 0 votes, 0 margin 

[123] PPC Vijay Kumar - 0 votes, 0 margin 

[123] IND Manmohan Khroud - 0 votes, 0 margin 

`,
		},
		{
			id:    999,
			error: true,
		},
		// {
		// 	id:   304,
		// 	want: "West Vancouver – Sunshine Coast – Sea to Sky Country",
		// },
	} {
		t.Run(strconv.Itoa(tc.id), func(t *testing.T) {
			got, err := findRidingByIDSummary(tc.id)
			if err != nil {
				if tc.error {
					return
				}
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("output mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
