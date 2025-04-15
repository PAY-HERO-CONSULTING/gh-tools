package custom_types

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertToProtobufTime(time time.Time) *timestamppb.Timestamp {
	return timestamppb.New(time)
}
