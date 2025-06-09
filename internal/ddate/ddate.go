package ddate

import (
	"fmt"
	"time"
)

type WeekDay int

const (
	Sweetmorn WeekDay = iota
	Boomtime
	Pungenday
	PricklePrickle
	SettingOrange
)

func (w WeekDay) String() string {
	return []string{"Sweetmorn", "Boomtime", "Pungenday", "Prickle-Prickle", "Setting Orange"}[w]
}

type Season int

const (
	Chaos Season = iota
	Discord
	Confusion
	Bureaucracy
	TheAftermath
)

func (s Season) String() string {
	return []string{"Chaos", "Discord", "Confusion", "Bureaucracy", "The Aftermath"}[s]
}

type HolyDay string

const (
	StTibsDay HolyDay = "St. Tib's Day"
	Mungday   HolyDay = "Mungday"
	Mojoday   HolyDay = "Mojoday"
	Syaday    HolyDay = "Syaday"
	Zaraday   HolyDay = "Zaraday"
	Maladay   HolyDay = "Maladay"
	Chaoflux  HolyDay = "Chaoflux"
	Discoflux HolyDay = "Discoflux"
	Confuflux HolyDay = "Confuflux"
	Bureflux  HolyDay = "Bureflux"
	Afflux    HolyDay = "Afflux"
)

type Date struct {
	Year      int
	Season    Season
	SeasonDay int
	WeekDay   WeekDay
	HolyDay   HolyDay
}

// Return a Discordian Date object from the given time
func FromTime(greg time.Time) Date {
	dis := Date{}
	dis.Year = greg.Year() + 1166

	disYearDay := greg.YearDay() - 1 // [0-364]
	if isLeapYear(greg.Year()) {
		if greg.Month() == time.February && greg.Day() == 29 {
			dis.HolyDay = StTibsDay
			return dis
		}
		if greg.Month() >= time.March {
			disYearDay -= 1 // keep it [60-364]
		}
	}

	// zero indexed season, [0-4]
	dis.Season = Season(disYearDay / 73)

	// one indexed day of the season, [1-73]
	dis.SeasonDay = disYearDay%73 + 1

	if dis.SeasonDay == 5 {
		// apostle days
		dis.HolyDay = []HolyDay{Mungday, Mojoday, Syaday, Zaraday, Maladay}[dis.Season]
	} else if dis.SeasonDay == 50 {
		// flux days
		dis.HolyDay = []HolyDay{Chaoflux, Discoflux, Confuflux, Bureflux, Afflux}[dis.Season]
	}

	// zero indexed day of the week, [0-4]
	dis.WeekDay = WeekDay(disYearDay % 5)

	return dis
}

// Return a discordian Date object corresponding to the current time in location
func NowIn(location *time.Location) Date {
	return FromTime(time.Now().In(location))
}

// Format Date as a string
func (d Date) Format(showHolydays bool) string {
	str := fmt.Sprintf("%s, %s %d, %d YOLD", d.WeekDay, d.Season, d.SeasonDay, d.Year)
	if d.HolyDay != "" {
		if d.HolyDay == StTibsDay {
			return fmt.Sprintf("%s, %d YOLD", StTibsDay, d.Year)
		}
		if showHolydays {
			return fmt.Sprintf("%s (%s)", str, d.HolyDay)
		}
	}
	return str
}

// Return true if year is a leap year
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
