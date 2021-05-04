package model

type Booking struct {
	// Name   MyNullString `json:"name"`
	// Price  int    `json:"price"`
	// Author MyNullString `json:"author"`

	BookNo           MyNullString `json:"book_no,string"`
	Operator         MyNullString `json:"operator,string"`
	Customer         MyNullString `json:"customer,string"`
	VoyageNo         MyNullString `json:"yoyage_no,string"`
	Destination      MyNullString `json:"destination,string"`
	VesselName       MyNullString `json:"vessel_name,string"`
	PickupDate       MyNullTime   `json:"pickup_date,omitempty"`
	GoodsDescription MyNullString `json:"goods_description,string"`
	Remark           MyNullString `json:"remark,string"`

	BookingContainerTypes   []BookingContainerType   `json:"bookingContainerTypes"`
	BookingContainerDetails []BookingContainerDetail `json:"bookingContainerDetails"`
}

type BookingContainerType struct {
	BookNo    MyNullString `json:"book_no,string"`
	Size      int          `json:"size"`
	Type      MyNullString `json:"type"`
	Quantity  int          `json:"quantity"`
	Available int          `json:"available"`
	TotalOut  int          `json:"total_out"`
}

type BookingContainerDetail struct {
	BookNo      MyNullString `json:"book_no,string"`
	No          MyNullString `json:"no"`
	ContainerNo MyNullString `json:"container_no,string"`
	Size        int          `json:"size"`
	Type        MyNullString `json:"type,string"`
	SealNo      MyNullString `json:"seal_no,string"`
	TrailerName MyNullString `json:"trailer_name,string"`
	License     MyNullString `json:"license,string"`
	GateOutDate MyNullTime   `json:"gate_out_date"`
}
