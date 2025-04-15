package serviceutils

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ParseTime(
	timeString string,
) (time.Time, error) {

	parsedTime, err := time.ParseInLocation(timeLayout, timeString, time.UTC)
	if err != nil {

		parsedTime, err = time.ParseInLocation(timeandtimezonelayout, timeString, time.UTC)
		if err != nil {
			return parsedTime, fmt.Errorf(
				fmt.Sprintf("unable to parse time err = [%+v]", err),
			)
		}
	}

	return parsedTime, nil
}

func FormatDate(timeToFormat time.Time) string {
	return timeToFormat.In(time.UTC).Format(dateLayout)
}

func FormatTime(
	timeToFormat time.Time,
) string {
	return timeToFormat.In(time.UTC).Format(timeLayout)
}

func FormatDateTime(
	timeToFormat time.Time,
) string {
	return timeToFormat.In(time.UTC).Format(safaricomTimeLayout)
}

func ConvertToProtobufTime(time time.Time) *timestamppb.Timestamp {
	return timestamppb.New(time)
}
