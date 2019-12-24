package algorithm

import (
	"fmt"
	"github.com/reschedulize/school_course_data"
	"sync"
	"time"
)

func generateLCCombos(sets [][]*school_course_data.Class) (output []*Schedule) {
	prev := make([]*Schedule, len(sets[0]))

	for i, c := range sets[0] {
		sch := &Schedule{}
		prev[i] = sch.permutate(c)
	}

	for i := 1; i < len(sets); i++ {
		var good []*Schedule

		for _, x := range sets[i] {
			for _, y := range prev {
				err := y.check(x)

				if err == nil {
					good = append(good, y.permutate(x))
				}
			}
		}

		prev = good
	}

	return prev
}

func generatePossibleSchedules(sets [][]*Schedule) []*Schedule {
	prev := sets[0]

	for i := 1; i < len(sets); i++ {
		var good []*Schedule

		for _, x := range sets[i] {
				for _, p := range prev {
					err := p.checkSchedule(x)

					if err == nil {
						good = append(good, p.permutateSchedule(x))
					}
				}
		}

		prev = good
	}

	return prev
}

func sortLinkedClasses(classes []*school_course_data.Class) []*Schedule {
	// Sort classes into link groups and link IDs
	// Ex. L2. Link group = 2, link ID = L2
	groups := make(map[string]map[string][]*school_course_data.Class)

	for i, c := range classes {
		linkGroup := c.LinkID[1:]

		_, exists := groups[linkGroup]

		if !exists {
			groups[linkGroup] = make(map[string][]*school_course_data.Class)
		}

		_, exists = groups[linkGroup][c.LinkID]

		if !exists {
			groups[linkGroup][c.LinkID] = []*school_course_data.Class{}
		}

		groups[linkGroup][c.LinkID] = append(groups[linkGroup][c.LinkID], classes[i])
	}

	// Filter out groups that have classes with a missing schedule
	for linkGroup, group := range groups {
		for linkID := range group {
			newIndex := 0

			for _, class := range group[linkID] {
				if class.WeekMask != 0 {
					group[linkID][newIndex] = class
					newIndex++
				}
			}

			group[linkID] = group[linkID][:newIndex]

			// If there is not at least 1 class from each class type, delete the entire group
			if len(group[linkID]) == 0 {
				delete(groups, linkGroup)
				break
			}
		}
	}

	// Convert map to slice
	linkedSets := make([]*Schedule, 0, len(groups))

	for _, group := range groups {
		sets := make([][]*school_course_data.Class, 0, len(group))

		for _, c := range group {
			sets = append(sets, c)
		}

		linkedSets = append(linkedSets, generateLCCombos(sets)...)
	}

	return linkedSets
}

func makeLCCombinations(api school_course_data.IAPI, term string, courses []string) ([]*Schedule, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error

	result := make([]*Schedule, 0, len(courses))

	for _, course := range courses {
		wg.Add(1)

		go func(course string) {
			defer wg.Done()

			var classes []*school_course_data.Class
			classes, err = api.Classes(term, course, 100)

			if err != nil {
				return
			}

			sortedClasses := sortLinkedClasses(classes)

			mutex.Lock()
			result = append(result, sortedClasses...)
			mutex.Unlock()
		}(course)
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func Solve(api school_course_data.IAPI, term string, courseGroups [][]string) ([][]string, error) {
	start := time.Now().UnixNano()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var err error

	var scheduleSlots [][]*Schedule

	for _, courses := range courseGroups {
		wg.Add(1)

		go func(courses []string) {
			defer wg.Done()

			var linkedGroups []*Schedule
			linkedGroups, err = makeLCCombinations(api, term, courses)

			if err != nil {
				return
			}

			mutex.Lock()
			scheduleSlots = append(scheduleSlots, linkedGroups)
			mutex.Unlock()
		}(courses)
	}

	wg.Wait()

	if err != nil {
		return nil, err
	}

	possibilities := generatePossibleSchedules(scheduleSlots)

	fmt.Println(time.Now().UnixNano() - start)

	output := make([][]string, len(possibilities))

	for i, schedule := range possibilities {
		classes := schedule.Classes
		output[i] = make([]string, len(classes))

		for ii, class := range classes {
			output[i][ii] = class.CRN
		}
	}

	return output, nil
}
