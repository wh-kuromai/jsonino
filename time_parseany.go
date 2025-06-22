package jsonino

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var vaguePatternRegex = regexp.MustCompile(`[0-9][0-9]?\.([0-9][0-9][0-9][0-9])`)
var delimRegex = regexp.MustCompile(`[^0-9a-zA-Z]+`)

var numRegex = regexp.MustCompile(`^[0-9]+$`)
var tzoffsetRegex = regexp.MustCompile(`([\+\-])([0-9][0-9]?)[^0-9]?([0-9]+)?`)
var yearRegex = regexp.MustCompile(`^[1-2][019][0-9][0-9]$`)

var yyyymmddRegex = regexp.MustCompile(`^([1-2][019][0-9][0-9])([0-9][0-9])([0-9][0-9])?$`)

var timeRegex = regexp.MustCompile(`[0-9]+|[a-zA-Z]+`)
var monthMap = map[string]int{
	"JAN": 1,
	"FEB": 2,
	"MAR": 3,
	"APR": 4,
	"MAY": 5,
	"JUN": 6,
	"JUL": 7,
	"AUG": 8,
	"SEP": 9,
	"OCT": 10,
	"NOV": 11,
	"DEC": 12,
}

var timezoneMap = map[string]float64{"DST": -12, "UTC": -11, "HST": -10, "AKDT": -8, "PDT": -7, "PST": -8, "UMST": -7, "MDT": -6, "CAST": -6, "CDT": -5, "CCST": -6, "SPST": -5, "EST": -5, "EDT": -4, "UEDT": -5, "VST": -4.5, "PYT": -4, "ADT": -3, "CBST": -4, "SWST": -4, "PSST": -4, "NDT": -2.5, "ESAST": -3, "AST": -3, "SEST": -3, "GDT": -3, "MST": -3, "BST": -3, "U": -2, "CVST": -1, "GMT": 0, "GST": 0, "WEDT": 2, "CEDT": 2, "RDT": 2, "WCAST": 1, "NST": 1, "MEDT": 3, "SDT": 3, "EEDT": 3, "SAST": 2, "FDT": 3, "TDT": 3, "JDT": 3, "LST": 2, "KST": 3, "EAST": 3, "MSK": 3, "SAMT": 4, "IDT": 4.5, "GET": 4, "CST": 4, "WAST": 5, "YEKT": 5, "PKT": 5, "IST": 5.5, "SLST": 5.5, "NCAST": 7, "NAST": 8, "MPST": 8, "TST": 8, "UST": 8, "NAEST": 8, "JST": 9, "ACST": 9.5, "AEST": 10, "WPST": 10, "YST": 9, "CPST": 11, "NZST": 12, "FST": 12, "KDT": 13, "SST": 13}

func ParseAnyTime(date string) time.Time {

	date = strings.ToUpper(date)

	tz := "UTC"
	tzoffset := 0.0
	tzoffsetFound := false

	// Find TZ Offset ( +0900 )
	//   TZ Offset must be after 5 ( DD-MM-HH-MM TZOFFSET)
	//   avoid parsing 2014-09-09 as -0909

	didx := delimRegex.FindAllStringIndex(date, 4)
	if len(didx) >= 4 {
		tzoffdate := date[didx[3][0]:]

		sm := tzoffsetRegex.FindStringSubmatchIndex(tzoffdate)
		if sm != nil {
			tzoffsetFound = true
			plusminus := tzoffdate[sm[2]:sm[3]]
			tzhour, _ := strconv.Atoi(strings.TrimLeft(tzoffdate[sm[4]:sm[5]], "0"))
			tzoffset += float64(tzhour)
			if sm[6] > 0 && sm[7] > 0 {
				tzmin, _ := strconv.Atoi(strings.TrimLeft(tzoffdate[sm[6]:sm[7]], "0"))
				tzoffset += float64(tzmin) / 60.0
			}
			if plusminus == "-" {
				tzoffset = -tzoffset
			}

			date = date[:didx[4][0]] + tzoffdate[:sm[0]] + tzoffdate[sm[1]:]
			tz = ""
		}

	}

	// add leading 0 to nsec like 00:00:00.2019
	vp := vaguePatternRegex.FindStringSubmatchIndex(date)
	if vp != nil {
		date = date[:vp[2]] + "0" + date[vp[2]:]
	}

	// Find Nums and Alphas
	dateidx := timeRegex.FindAllStringIndex(date, -1)
	tgt := make([]string, len(dateidx))
	for i := range dateidx {
		tgt[i] = date[dateidx[i][0]:dateidx[i][1]]
	}

	// Find TZ
	for i := range tgt {
		offset, ok := timezoneMap[tgt[i]]
		if ok {
			tz = tgt[i]
			if !tzoffsetFound {
				tzoffset = offset
			}

			tgt = removeElement(tgt, i)
			break
		}
	}

	// Month to Num (Jan -> 1)
	month := 0
	monthFound := false
	for i := range tgt {
		for k, v := range monthMap {
			if strings.HasPrefix(tgt[i], k) {
				month = v
				monthFound = true
				break
			}
		}
	}

	// REMOVE Non Nums
	tgt2 := make([]string, 0, len(tgt))
	for i := range tgt {
		if numRegex.MatchString(tgt[i]) {
			tgt2 = append(tgt2, tgt[i])
		}
	}

	// YYYYMMDD -> YYYY, MM, DD
	for i := range tgt2 {
		sm := yyyymmddRegex.FindStringSubmatchIndex(tgt2[i])
		if sm != nil {
			if sm[6] < 0 || sm[7] < 0 {
				tgt2 = append(tgt2[:i], tgt2[i][sm[2]:sm[3]], tgt2[i][sm[4]:sm[5]])
				tgt2 = append(tgt2, tgt2[i+1:]...)
			} else {
				tgt2 = append(tgt2[:i], tgt2[i][sm[2]:sm[3]], tgt2[i][sm[4]:sm[5]], tgt2[i][sm[6]:sm[7]])
				tgt2 = append(tgt2, tgt2[i+1:]...)
			}
			break
		}
	}

	// Find YYYY
	year := 0
	yearFound := false
	for i := range tgt2 {
		if yearRegex.MatchString(tgt2[i]) {
			year, _ = strconv.Atoi(strings.TrimLeft(tgt2[i], "0"))
			tgt2 = removeElement(tgt2, i)
			yearFound = true
			break
		}
	}

	// func shiftNum
	shiftNum := func(idx int) (int, bool) {
		if len(tgt2) > idx {
			num, _ := strconv.Atoi(strings.TrimLeft(tgt2[idx], "0"))
			tgt2 = removeElement(tgt2, 0)
			return num, true
		}
		return 0, false
	}

	day := 0
	if monthFound {
		// already found Alpha month.
		// DD MMM YY

		day, _ = shiftNum(0)

		if !yearFound {
			year, yearFound = shiftNum(0)
		}
	} else {
		// Num month.

		if yearFound {
			// (YYYY) MM DD
			// MM DD (YYYY)

			month, _ = shiftNum(0)
			day, _ = shiftNum(0)
		} else {
			// MM DD YY
			month, _ = shiftNum(0)
			day, _ = shiftNum(0)
			year, _ = shiftNum(0)

			//if month > 12 {
			// YY MM DD
			//year, month, day = month, day, year
			//}
		}

	}

	if month == 0 {
		month = 1
	}

	if day == 0 {
		day = 1
	}

	if year < 100 {
		if year > 50 {
			year = 1900 + year
		} else {
			year = 2000 + year
		}
	}

	hour, _ := shiftNum(0)
	min, _ := shiftNum(0)
	sec, _ := shiftNum(0)
	nsec, _ := shiftNum(0)

	//fmt.Println(year, month, day, hour, min, sec, nsec, tz, tzoffset)

	_ = yearFound
	return time.Date(year, (time.Month)(month), day, hour, min, sec, nsec, time.FixedZone(tz, int(tzoffset*60*60)))
	//panic("format error")
}

func removeElement[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
