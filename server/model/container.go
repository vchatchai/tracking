package model

type User struct {
	User MyNullString `json:"user,string"`
	ID   MyNullString `json:"id,string"`
}

type LadenContainer struct {
	ContainerNo MyNullString `json:"container_no,string"`
	Size        int          `json:"size"`
	Type        MyNullString `json:"type,string"`
	BookingNo   MyNullString `json:"booking_no,string"`
	SealNo      MyNullString `json:"seal_no,string"`
	Customer    MyNullString `json:"customer,string"`
	LDCode      MyNullString `json:"ld_code,string"`
	Origin      MyNullString `json:"origin,string"`
	Destination MyNullString `json:"destination,string"`
	Vessel      MyNullString `json:"vessel,string"`
	VoyageNo    MyNullString `json:"voyage_no,string"`
	Renban      MyNullString `json:"renban,string"`
	CYDate      MyNullTime   `json:"cy_date"`

	GateInTrailerName MyNullString `json:"gate_in_trailer_name,string"`
	GateInLicense     MyNullString `json:"gate_in_license,string"`
	GateInDate        MyNullTime   `json:"gate_in_date"`
	GateInLocation    MyNullString `json:"gate_in_location,string"`

	GateOutTrailerName MyNullString `json:"gate_out_trailer_name,string"`
	GateOutLicense     MyNullString `json:"gate_out_license,string"`
	GateOutDate        MyNullTime   `json:"gate_out_date"`
}
