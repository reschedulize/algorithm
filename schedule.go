package algorithm

import "github.com/reschedulize/school_course_data"

type Schedule struct {
	Classes  []*school_course_data.Class
	WeekMask uint8
	DayMasks [7][4]uint64
}

func (s *Schedule) check(class *school_course_data.Class) error {
	if s.WeekMask & class.WeekMask == 0 {
		return nil
	}

	for i := 0; i < 7; i++ {
		for j := 0; j < 4; j++ {
			if s.DayMasks[i][j] & class.DayMasks[i][j] != 0 {
				return timeConflictError
			}
		}
	}

	return nil
}

func (s *Schedule) checkSchedule(schedule *Schedule) error {
	if s.WeekMask & schedule.WeekMask == 0 {
		return nil
	}

	for i := 0; i < 7; i++ {
		for j := 0; j < 4; j++ {
			if s.DayMasks[i][j] & schedule.DayMasks[i][j] != 0 {
				return timeConflictError
			}
		}
	}

	return nil
}

func (s *Schedule) permutate(class *school_course_data.Class) *Schedule {
	newSch := &Schedule{}

	newSch.WeekMask = s.WeekMask | class.WeekMask

	for i := 0; i < 7; i++ {
		for j := 0; j < 4; j++ {
			newSch.DayMasks[i][j] = s.DayMasks[i][j] | class.DayMasks[i][j]
		}
	}

	newSch.Classes = append(s.Classes, class)

	return newSch
}

func (s *Schedule) permutateSchedule(schedule *Schedule) *Schedule {
	newSch := &Schedule{}

	newSch.WeekMask = s.WeekMask | schedule.WeekMask

	for i := 0; i < 7; i++ {
		for j := 0; j < 4; j++ {
			newSch.DayMasks[i][j] = s.DayMasks[i][j] | schedule.DayMasks[i][j]
		}
	}

	newSch.Classes = append(s.Classes, schedule.Classes...)

	return newSch
}
