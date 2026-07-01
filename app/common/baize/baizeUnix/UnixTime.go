package baizeUnix

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

type BaiZeTime time.Time

// MarshalJSON implements json.Marshaler.
func (t BaiZeTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(stamp), nil
}

func (t *BaiZeTime) UnmarshalJSON(b []byte) error {
	parseInt, _ := strconv.ParseInt(string(b), 10, 64)
	*t = BaiZeTime(time.Unix(parseInt, 0))
	return nil
}
func (t BaiZeTime) ToString() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

func (t *BaiZeTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*t = BaiZeTime(v)
		return nil
	case []byte:
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", string(v), time.Local)
		if err != nil {
			return err
		}
		*t = BaiZeTime(parsed)
		return nil
	case string:
		parsed, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if err != nil {
			return err
		}
		*t = BaiZeTime(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported BaiZeTime scan type: %T", value)
	}
}

func (t BaiZeTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}
