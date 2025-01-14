package kustomersdk

import (
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/pkg/errors"
)

type ICSStatus string

const (
	ICSStatusConfirmed ICSStatus = "CONFIRMED"
	ICSStatusCancelled ICSStatus = "CANCELLED"
)

type ICSParams struct {
	ID          string
	Start       string
	startParsed time.Time
	End         string
	endParsed   time.Time
	TZ          string
	Title       string
	Description string
	Location    string
	Status      ICSStatus
	Sequence    uint8
	Attendees   []*ICSAttendee
}

type ICSAttendee struct {
	FirstName string
	LastName  string
	Email     string
}

type ICSFile struct {
	Params *ICSParams
}

var DefaultOrganizer = &ICSAttendee{
	FirstName: "Superpower",
	LastName:  "Health",
	Email:     "noreply@superpower.com",
}

// NewICSFile creates a new ICS file using the provided configuration and attendees list
func NewICSFile(p *ICSParams, attendees []*ICSAttendee) (*ICSFile, error) {
	if err := p.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate ICS params")
	}

	if attendees == nil {
		return nil, errors.New("attendees is required")
	}

	p.Attendees = append(attendees, DefaultOrganizer)

	return &ICSFile{
		Params: p,
	}, nil
}

// ToFile returns the ICS file as a byte slice
func (i *ICSFile) ToFile() []byte {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetCalscale("GREGORIAN")
	cal.SetXPublishedTTL("PT1H")

	event := cal.AddEvent(i.Params.ID)
	event.SetCreatedTime(time.Now().UTC())
	event.SetDtStampTime(time.Now().UTC())
	event.SetModifiedAt(time.Now().UTC())
	event.SetStartAt(i.Params.startParsed)
	event.SetEndAt(i.Params.endParsed)
	event.SetSummary(i.Params.Title)
	event.SetLocation(i.Params.Location)
	event.SetDescription(i.Params.Description)
	event.SetOrganizer(DefaultOrganizer.Email, ics.WithCN(DefaultOrganizer.FirstName+" "+DefaultOrganizer.LastName))
	event.SetStatus("CONFIRMED")

	if i.Params.Status == ICSStatusCancelled {
		event.SetSequence(1)
	} else {
		event.SetSequence(0)
	}

	for _, att := range i.Params.Attendees {
		event.AddAttendee(
			att.Email,
			ics.CalendarUserTypeIndividual,
			ics.ParticipationRoleReqParticipant,
			ics.ParticipationStatusAccepted,
			ics.WithRSVP(true),
			ics.WithCN(att.FirstName+" "+att.LastName),
		)
	}

	return []byte(cal.Serialize())
}

func (p *ICSParams) validate() error {
	if p.ID == "" {
		return errors.New("id is required")
	}

	if p.Start == "" {
		return errors.New("start is required")
	}

	// Attempt to parse the start time
	startParsed, err := time.Parse(time.RFC3339, p.Start)
	if err != nil {
		return errors.Wrap(err, "failed to parse start time")
	}
	p.startParsed = startParsed

	if p.End == "" {
		return errors.New("end is required")
	}

	// Attempt to parse the end time
	endParsed, err := time.Parse(time.RFC3339, p.End)
	if err != nil {
		return errors.Wrap(err, "failed to parse end time")
	}
	p.endParsed = endParsed

	if p.TZ == "" {
		return errors.New("tz is required")
	}

	if p.Status == ICSStatusConfirmed {
		// These params are only required if the status is CONFIRMED
		if p.Title == "" {
			return errors.New("title is required")
		}

		if p.Description == "" {
			return errors.New("description is required")
		}
	}

	if p.Status == ICSStatusCancelled {
		p.Description = ""
	}

	if p.Location == "" {
		return errors.New("location is required")
	}

	return nil
}
