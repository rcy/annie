package day

import (
	"strings"
	"testing"
)

func TestStripPhrases(t *testing.T) {
	words := stripPhrases([]string{
		"National Tree Day",
		"National Crime Prevention Month",
		"Auntie's Day",
		"National Soccer Day",
		"World Nature Conservation Day",
		"National Milk Chocolate Day",
		"World Day for Grandparents and the Elderly",
		"National Parents' Day",
		"World Hepatitis Day",
		"National Marine Week",
		"Festival of British Archaeology",
		"National Moth Week",
		"Love Parks Week",
		"Beans Month",
		"National Anti-Boredom Month",
		"Lasagna Awareness Month",
		"Good Care Month",
		"Wild About Wildlife Month",
		"National Culinary Arts Month",
		"National Hemp Month",
		"National Powersports Month",
		"Independent Retailer Month",
		"National Cell Phone Courtesy Month",
		"National Horseradish Month",
		"National Picnic Month",
		"National Ice Cream Month",
		"Plastic Free July",
		"Bank Account Bonus Month",
		"Sarcoma Awareness Month",
		"World Watercolor Month",
	})

	got := strings.Join(words, ", ")
	want := "tree, crime prevention, auntie's, soccer, nature conservation, milk chocolate, grandparents elderly, parents', hepatitis, marine, festival of british archaeology, moth, love parks, beans, anti-boredom, lasagna awareness, good care, wild about wildlife, culinary arts, hemp, powersports, independent retailer, cell phone courtesy, horseradish, picnic, ice cream, plastic free july, bank account bonus, sarcoma awareness, watercolor"

	if got != want {
		t.Errorf("expected: %s\n got: %s", want, got)
	}
}
