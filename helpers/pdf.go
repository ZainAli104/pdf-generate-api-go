package helpers

import "pdf-generate/models"

func PunchToString(p models.Punch) string {
	switch p {
	case models.PunchIn:
		return "Punch In"
	case models.PunchOut:
		return "Punch Out"
	case models.BreakOut:
		return "Break Out"
	case models.BreakIn:
		return "Break In"
	case models.OvertimeIn:
		return "Overtime In"
	case models.OvertimeOut:
		return "Overtime Out"
	default:
		return "Unknown"
	}
}

func AttendanceStatusToString(status models.AttendanceStatus) string {
	switch status {
	case models.AttendanceOnTime:
		return "On Time"
	case models.AttendanceLate:
		return "Late"
	case models.AttendanceOff:
		return "Off"
	case models.AttendanceLeave:
		return "Leave"
	case models.AttendanceLeftEarly:
		return "Left Early"
	case models.AttendanceOther:
		return "Other"
	default:
		return "Unknown"
	}
}
