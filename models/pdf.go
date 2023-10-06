package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MemberPopulated struct {
	Id        *primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`

	User        User               `bson:"user" json:"user"`
	Department  Department         `bson:"department" json:"department"`
	Identifier  string             `bson:"identifier,omitempty" json:"identifier"`
	SSN         string             `bson:"ssn,omitempty" json:"ssn"`
	Phone       string             `bson:"phone,omitempty" json:"phone"`
	Gender      string             `bson:"gender,omitempty" json:"gender"`
	Address     string             `bson:"address,omitempty" json:"address"`
	DateOfBirth *time.Time         `bson:"dateOfBirth,omitempty" json:"dateOfBirth"`
	JoiningDate *time.Time         `bson:"joiningDate,omitempty" json:"joiningDate"`
	Nationality string             `bson:"nationality,omitempty" json:"nationality"`
	CreatedBy   primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
}

type Department struct {
	Id *primitive.ObjectID `bson:"_id" json:"id"`

	Name string `bson:"name" json:"name"`
}

type User struct {
	Id        *primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updatedAt" json:"updatedAt"`

	Scope          []string `bson:"scope" json:"scope"`
	Email          string   `bson:"email" json:"email"`
	Password       string   `bson:"password" json:"password,omitempty"`
	FirstName      string   `bson:"firstName" json:"firstName"`
	LastName       string   `bson:"lastName" json:"lastName"`
	ProfilePicture string   `bson:"profilePicture" json:"profilePicture,omitempty"`
}

type Punch int

const (
	PunchIn     Punch = 0
	PunchOut    Punch = 1
	BreakOut    Punch = 2
	BreakIn     Punch = 3
	OvertimeIn  Punch = 4
	OvertimeOut Punch = 5
)

type AttendanceStatus int

const (
	AttendanceOnTime    AttendanceStatus = 0
	AttendanceLate      AttendanceStatus = 1
	AttendanceOff       AttendanceStatus = 2
	AttendanceLeave     AttendanceStatus = 3
	AttendanceLeftEarly AttendanceStatus = 4
	AttendanceOther     AttendanceStatus = 5
)

type Attendance struct {
	Id        *primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`

	PunchTime time.Time          `bson:"punchTime" json:"punchTime"`
	Punch     Punch              `bson:"punch" json:"punch"`
	Status    AttendanceStatus   `bson:"status" json:"status"`
	MemberId  primitive.ObjectID `bson:"memberId" json:"memberId"`
}
