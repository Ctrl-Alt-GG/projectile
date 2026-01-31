package openttd

// Based on https://github.com/OpenTTD/OpenTTD/blob/1dd3d655747df41cda76bad54d06e86d6efa35ed/src/timer/timer_game_common.cpp#L55-L109

const (
	DAYS_IN_YEAR      = 365 ///< days per year
	DAYS_IN_LEAP_YEAR = 366 ///< sometimes, you need one day more...

)

func IsLeapYear(year uint32) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func FigureOutYear(date uint32) int {

	/* Year determination in multiple steps to account for leap
	 * years. First do the large steps, then the smaller ones.
	 */

	/* There are 97 leap years in 400 years */
	yr := 400 * (date / (DAYS_IN_YEAR*400 + 97))

	rem := date % (DAYS_IN_YEAR*400 + 97)

	if rem >= DAYS_IN_YEAR*100+25 {
		/* There are 25 leap years in the first 100 years after
		 * every 400th year, as every 400th year is a leap year */
		yr += 100
		rem -= DAYS_IN_YEAR*100 + 25

		/* There are 24 leap years in the next couple of 100 years */
		yr += 100 * (rem / (DAYS_IN_YEAR*100 + 24))
		rem = rem % (DAYS_IN_YEAR*100 + 24)
	}

	if !IsLeapYear(yr) && rem >= DAYS_IN_YEAR*4 {
		/* The first 4 year of the century are not always a leap year */
		yr += 4
		rem -= DAYS_IN_YEAR * 4
	}

	/* There is 1 leap year every 4 years */
	yr += 4 * (rem / (DAYS_IN_YEAR*4 + 1))
	rem = rem % (DAYS_IN_YEAR*4 + 1)

	/* The last (max 3) years to account for; the first one
	 * can be, but is not necessarily a leap year */
	var checkVal uint32 = DAYS_IN_YEAR
	if IsLeapYear(yr) {
		checkVal = DAYS_IN_LEAP_YEAR
	}
	for rem >= checkVal {
		yr++
		if IsLeapYear(yr) {
			rem -= DAYS_IN_LEAP_YEAR
			checkVal = DAYS_IN_LEAP_YEAR
		} else {
			rem -= DAYS_IN_YEAR
			checkVal = DAYS_IN_YEAR
		}
	}

	return int(yr)
}
