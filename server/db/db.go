package db

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vchatchai/tracking/server/model"
)

type DB interface {
	GetContainerByBookingNumber(bookingNumber string) ([]*model.LadenContainer, error)
	GetContainerContainerNumber(containerNumber string) ([]*model.LadenContainer, error)
	GetBookingByBookingNumber(bookingNumber string) (*model.Booking, error)
	GetBookingByContainerNumber(containerNumber string) (*model.Booking, error)
	Login(user, password string) (*model.User, error)
	ClearBooking(tx *sql.Tx) error
	ClearContainer(tx *sql.Tx) error
	ClearUser(tx *sql.Tx) error
	InsertBooking(tx *sql.Tx, bookings []*model.Booking) error
	InsertContainer(tx *sql.Tx, ladenContainers []*model.LadenContainer) error
	InsertUser(tx *sql.Tx, users []*model.User) error

	RefreshBooking([]*model.Booking) error
	RefreshContainer([]*model.LadenContainer) error
	RefreshUser([]*model.User) error
}

type SQLDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) DB {
	return SQLDB{db: db}
}

var queryBookingHeader = `
SELECT 	book_no, operator, customer, yoyage_no, destination, vessel_name, pickup_date, goods_description, remark
FROM booking_header
WHERE book_no = ?
`

//

var queryBookingType = `
SELECT book_no,size,type,quantity,available,total_out
FROM booking_container_type
WHERE book_no = ?
`

var queryBookingDetail = `
SELECT book_no,no,container_no,size,type,seal_no,trailer_name,license,gate_out_date
FROM booking_container_detail
WHERE book_no = ?
`

var queryBookingbyContainer = `
SELECT  d.book_no,d.no,d.container_no,d.size,d.type,d.seal_no,d.trailer_name,d.license,d.gate_out_date
FROM booking_container_detail d INNER JOIN booking_header h on d.book_no = h.book_no 
WHERE  container_no = ?
ORDER BY create_date desc limit 1
`

/**

 */
func (d SQLDB) GetBookingByBookingNumber(bookingNumber string) (*model.Booking, error) {

	var booking *model.Booking

	booking = &model.Booking{}

	stmtHeader, err := d.db.Prepare(queryBookingHeader)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmtHeader.Close()

	row := stmtHeader.QueryRow(bookingNumber)

	err = row.Scan(&booking.BookNo, &booking.Operator, &booking.Customer, &booking.VoyageNo, &booking.Destination, &booking.VesselName, &booking.PickupDate, &booking.GoodsDescription, &booking.Remark)
	if err != nil {

		if err != nil {
			return booking, nil
		}
	}

	rows, err := d.db.Query(queryBookingType, bookingNumber)
	if err != nil {
		return booking, nil
	}
	defer rows.Close()

	for rows.Next() {
		containerType := &model.BookingContainerType{}
		err := rows.Scan(&containerType.BookNo, &containerType.Size, &containerType.Type, &containerType.Quantity, &containerType.Available, &containerType.TotalOut)
		if err != nil {
			return booking, nil
		} else {
			booking.BookingContainerTypes = append(booking.BookingContainerTypes, *containerType)

		}
	}
	err = rows.Err()
	if err != nil {
		return booking, nil
	}

	rows, err = d.db.Query(queryBookingDetail, bookingNumber)
	if err != nil {

		return booking, nil
	}
	defer rows.Close()

	for rows.Next() {
		containerDetail := &model.BookingContainerDetail{}
		err := rows.Scan(&containerDetail.BookNo, &containerDetail.No, &containerDetail.ContainerNo, &containerDetail.Size, &containerDetail.Type, &containerDetail.SealNo, &containerDetail.TrailerName, &containerDetail.License, &containerDetail.GateOutDate)
		if err != nil {

			return booking, nil
		} else {
			booking.BookingContainerDetails = append(booking.BookingContainerDetails, *containerDetail)
		}
	}
	err = rows.Err()
	if err != nil {

		return booking, nil
	}

	return booking, nil
}

func (d SQLDB) GetBookingByContainerNumber(containerNumber string) (*model.Booking, error) {

	var booking *model.Booking

	booking = &model.Booking{}

	rowBookingNumber := d.db.QueryRow(queryBookingbyContainer, containerNumber)

	containerDetail := &model.BookingContainerDetail{}
	err := rowBookingNumber.Scan(&containerDetail.BookNo, &containerDetail.No, &containerDetail.ContainerNo, &containerDetail.Size, &containerDetail.Type, &containerDetail.SealNo, &containerDetail.TrailerName, &containerDetail.License, &containerDetail.GateOutDate)
	if err != nil {
		fmt.Println("rowBookingNumber error", err)
		return booking, nil
	}

	stmtHeader, err := d.db.Prepare(queryBookingHeader)
	if err != nil {
		fmt.Println("Prepare stmtHeader failed:", err.Error())
	}
	defer stmtHeader.Close()

	row := stmtHeader.QueryRow(containerDetail.BookNo.Value())

	err = row.Scan(&booking.BookNo, &booking.Operator, &booking.Customer, &booking.VoyageNo, &booking.Destination, &booking.VesselName, &booking.PickupDate, &booking.GoodsDescription, &booking.Remark)
	if err != nil {

		fmt.Println("Prepare row failed:", err.Error(), booking)
		return booking, nil
	}

	fmt.Println("booking", booking)

	rows, err := d.db.Query(queryBookingType, containerDetail.BookNo.Value())
	if err != nil {
		return booking, nil

	}
	defer rows.Close()

	for rows.Next() {
		containerType := &model.BookingContainerType{}
		err := rows.Scan(&containerType.BookNo, &containerType.Size, &containerType.Type, &containerType.Quantity, &containerType.Available, &containerType.TotalOut)
		if err != nil {

			return booking, nil
		} else {
			booking.BookingContainerTypes = append(booking.BookingContainerTypes, *containerType)

		}
	}
	err = rows.Err()
	if err != nil {
		return booking, nil
	}

	booking.BookingContainerDetails = append(booking.BookingContainerDetails, *containerDetail)

	return booking, nil
}
func query(db *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	return rows, err

}

var queryContainerByBookingNo = ` 

SELECT container_no
           ,size
           ,type
           ,book_no
           ,seal_no
           ,customer
           ,id_code
           ,destination
           ,vessel
           ,voyage_no
           ,renban
           ,cy_date
           ,gate_in_trailer_name
           ,gate_in_license
           ,gate_in_date
           ,gate_in_location
           ,gate_out_trailer_name
           ,gate_out_license
		   ,gate_out_date
FROM laden_container
WHERE book_no = ?		 
ORDER BY gate_in_date DESC  
`

/*

 */
func (d SQLDB) GetContainerByBookingNumber(bookingNumber string) ([]*model.LadenContainer, error) {

	var ladenContainers []*model.LadenContainer

	rows, err := d.db.Query(queryContainerByBookingNo, bookingNumber)
	if err != nil {
		fmt.Println("GetContainerByBookingNumber query error", err)
		return ladenContainers, nil
	}
	defer rows.Close()

	for rows.Next() {
		ladenContainer := &model.LadenContainer{}
		err = rows.Scan(&ladenContainer.ContainerNo,
			&ladenContainer.Size,
			&ladenContainer.Type,
			&ladenContainer.BookingNo,
			&ladenContainer.SealNo,
			&ladenContainer.Customer,
			&ladenContainer.LDCode,
			&ladenContainer.Destination,
			&ladenContainer.Vessel,
			&ladenContainer.VoyageNo,
			&ladenContainer.Renban,
			&ladenContainer.CYDate,
			&ladenContainer.GateInTrailerName,
			&ladenContainer.GateInLicense,
			&ladenContainer.GateInDate,
			&ladenContainer.GateInLocation,
			&ladenContainer.GateOutTrailerName,
			&ladenContainer.GateOutLicense,
			&ladenContainer.GateOutDate,
		)
		if err != nil {
			fmt.Println("GetContainerByBookingNumber scan error", err)

			return ladenContainers, nil
		} else {
			ladenContainers = append(ladenContainers, ladenContainer)
		}
	}
	err = rows.Err()
	if err != nil {

		return ladenContainers, nil
	}

	return ladenContainers, nil
}

var queryContainerByContainerNo = ` 
SELECT container_no
           ,size
           ,type
           ,book_no
           ,seal_no
           ,customer
           ,id_code
           ,destination
           ,vessel
           ,voyage_no
           ,renban
           ,cy_date
           ,gate_in_trailer_name
           ,gate_in_license
           ,gate_in_date
           ,gate_in_location
           ,gate_out_trailer_name
           ,gate_out_license
		   ,gate_out_date
FROM tracking.laden_container
WHERE container_no = ?  
ORDER BY gate_in_date DESC limit 1
 
`

func (d SQLDB) GetContainerContainerNumber(containerNumber string) ([]*model.LadenContainer, error) {

	var ladenContainers []*model.LadenContainer

	rows, err := d.db.Query(queryContainerByContainerNo, containerNumber)
	if err != nil {
		fmt.Println("GetContainerContainerNumber query error", err)
		return ladenContainers, nil
	}
	defer rows.Close()

	for rows.Next() {
		ladenContainer := &model.LadenContainer{}
		err = rows.Scan(&ladenContainer.ContainerNo,
			&ladenContainer.Size,
			&ladenContainer.Type,
			&ladenContainer.BookingNo,
			&ladenContainer.SealNo,
			&ladenContainer.Customer,
			&ladenContainer.LDCode,
			&ladenContainer.Destination,
			&ladenContainer.Vessel,
			&ladenContainer.VoyageNo,
			&ladenContainer.Renban,
			&ladenContainer.CYDate,
			&ladenContainer.GateInTrailerName,
			&ladenContainer.GateInLicense,
			&ladenContainer.GateInDate,
			&ladenContainer.GateInLocation,
			&ladenContainer.GateOutTrailerName,
			&ladenContainer.GateOutLicense,
			&ladenContainer.GateOutDate,
		)
		if err != nil {
			fmt.Println("GetContainerContainerNumber row error", err)
			return ladenContainers, nil
		} else {
			ladenContainers = append(ladenContainers, ladenContainer)
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("GetContainerContainerNumber final error", err)
		return ladenContainers, nil
	}

	return ladenContainers, nil
}

var queryLogin = `
SELECT username
FROM tracking.laden_container_user
WHERE username = ? AND password = ?
`

func (d SQLDB) Login(userName, password string) (*model.User, error) {

	var user *model.User

	rows, err := d.db.Query(queryLogin, userName, password)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return user, nil

	}
	defer rows.Close()

	for rows.Next() {
		user = &model.User{}
		err := rows.Scan(&user.User)

		fmt.Println("scan ", user)
		if err != nil {

			return user, nil
		} else {
			h := sha1.New()
			h.Write([]byte(userName + password))

			user.ID = model.MyNullString{sql.NullString{string(h.Sum(nil)), true}}

		}

		break
	}

	return user, nil
}

var clearLogin = `DELETE FROM laden_container_user`

func (d SQLDB) ClearUser(tx *sql.Tx) error {

	_, err := d.db.Exec(clearLogin) // OK

	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}

	return nil
}

var clearBookingHeader = `DELETE FROM booking_header`
var clearBookingDetail = `DELETE FROM booking_container_detail`
var clearBookingType = `DELETE FROM booking_container_type`

func (d SQLDB) ClearBooking(tx *sql.Tx) error {

	_, err := d.db.Exec(clearBookingHeader) // OK

	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}

	_, err = d.db.Exec(clearBookingDetail) // OK

	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}
	_, err = d.db.Exec(clearBookingType) // OK

	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}

	return nil

}

var insertBookingHeader = `INSERT INTO booking_header(	book_no,
	operator,
	customer,
	yoyage_no,
	destination,
	vessel_name,
	pickup_date,
	goods_description,
	remark,
	create_date) VALUES(?,?,?,?,?,?,?,?,?,NOW())`
var insertBookingDetail = `INSERT INTO booking_container_detail(book_no,
	no,
	container_no,
	size,
	type,
	seal_no,
	trailer_name,
	license,
	gate_out_date
	) VALUES(?,?,?,?,?,?,?,?,?)`
var insertBookingType = `INSERT INTO booking_container_type(book_no,
	size,
	type,
	quantity,
	available,
	total_out
	) VALUES(?,?,?,?,?,?)`

func (d SQLDB) InsertBooking(tx *sql.Tx, bookings []*model.Booking) error {

	stmtHeader, err := d.db.Prepare(insertBookingHeader)
	if err != nil {
		log.Fatal(err)
	}

	stmtDetail, err := d.db.Prepare(insertBookingDetail)
	if err != nil {
		log.Fatal(err)
	}

	stmtType, err := d.db.Prepare(insertBookingType)
	if err != nil {
		log.Fatal(err)
	}

	for _, booking := range bookings {
		// fmt.Println(booking)

		book_no := booking.BookNo.Value()
		operator := booking.Operator.Value()
		customer := booking.Customer.Value()
		yoyage_no := booking.VoyageNo.Value()
		destination := booking.Destination.Value()
		vessel_name := booking.VesselName.Value()
		pickup_date := booking.PickupDate.Value()
		goods_description := booking.GoodsDescription.Value()
		remark := booking.Remark.Value()
		res, err := stmtHeader.Exec(book_no,
			operator,
			customer,
			yoyage_no,
			destination,
			vessel_name,
			pickup_date,
			goods_description,
			remark)
		if err != nil {
			log.Fatal(err)
		}
		_, err = res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		_, err = res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("%s ID = %d, affected = %d\n", book_no, lastId, rowCnt)

		for _, detail := range booking.BookingContainerDetails {
			book_no := detail.BookNo.Value()
			no := detail.No.Value()
			container_no := detail.ContainerNo.Value()
			size := detail.Size
			typeValue := detail.Type.Value()
			seal_no := detail.SealNo.Value()
			trailer_name := detail.TrailerName.Value()
			license := detail.License.Value()
			gate_out_date := detail.GateOutDate.Value()

			res, err := stmtDetail.Exec(book_no,
				no,
				container_no,
				size,
				typeValue,
				seal_no,
				trailer_name,
				license,
				gate_out_date)
			if err != nil {
				log.Fatal(err)
			}
			_, err = res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			_, err = res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, bookingtype := range booking.BookingContainerTypes {
			book_no := bookingtype.BookNo.Value()
			size := bookingtype.Size
			typeValue := bookingtype.Type.Value()
			quantity := bookingtype.Quantity
			available := bookingtype.Available
			total_out := bookingtype.TotalOut

			// fmt.Println(typeValue)
			res, err := stmtType.Exec(book_no,
				size,
				typeValue,
				quantity,
				available,
				total_out)
			if err != nil {
				log.Fatal(err)
			}
			_, err = res.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}
			_, err = res.RowsAffected()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	return nil
}

var clearContainer = `DELETE FROM laden_container`

func (d SQLDB) ClearContainer(tx *sql.Tx) error {
	fmt.Println("ClearContainer")
	_, err := d.db.Exec(clearContainer) // OK

	if err != nil {
		fmt.Printf("error %s\n", err)
		return err
	}

	return nil

}

var insertContainer = `INSERT INTO laden_container(
	container_no,
	size,
	type,
	book_no,
	seal_no,
	customer,
	id_code,
	destination,
	vessel,
	voyage_no,
	renban,
	cy_date,
	gate_in_trailer_name,
	gate_in_license,
	gate_in_date,
	gate_in_location,
	gate_out_trailer_name,
	gate_out_license,
	gate_out_date
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

func (d SQLDB) InsertContainer(tx *sql.Tx, ladenContainers []*model.LadenContainer) error {
	stmtHeader, err := tx.Prepare(insertContainer)
	if err != nil {
		log.Fatal(err)
	}
	defer stmtHeader.Close() // danger!
	for _, container := range ladenContainers {

		container_no := container.ContainerNo.Value()
		size := container.Size

		typeValue := container.Type.Value()
		book_no := container.BookingNo.Value()
		seal_no := container.SealNo.Value()
		customer := container.Customer.Value()
		id_code := container.LDCode.Value()
		destination := container.Destination.Value()
		vessel := container.Vessel.Value()
		voyage_no := container.VoyageNo.Value()
		renban := container.Renban.Value()
		cy_date := container.CYDate.Value()
		gate_in_trailer_name := container.GateInTrailerName.Value()
		gate_in_license := container.GateInLicense.Value()
		gate_in_date := container.GateInDate.Value()
		gate_in_location := container.GateInLocation.Value()
		gate_out_trailer_name := container.GateOutTrailerName.Value()
		gate_out_license := container.GateOutLicense.Value()
		gate_out_date := container.GateOutDate.Value()
		// fmt.Println(container)
		// for i := 0; i < len(container.ContainerNo.String); i++ {
		// 	fmt.Printf("%+q ", container.ContainerNo.String[i])
		// }
		// fmt.Println("container_no", container.ContainerNo.String+"'", len(string(container.ContainerNo.String)))
		// fmt.Println("size", size, "'")
		// fmt.Println("typeValue", typeValue, "'")
		// fmt.Println("book_no", book_no, "'")
		// fmt.Println("seal_no", seal_no, "'")
		// fmt.Println("customer", customer, "'")
		// fmt.Println("id_code", id_code, "'")
		// fmt.Println("destination", destination, "'")
		// fmt.Println("vessel", vessel, "'")
		// fmt.Println("voyage_no", voyage_no, "'")
		// fmt.Println("renban", renban, "'")
		// fmt.Println("cy_date", cy_date, "'")
		// fmt.Println("gate_in_trailer_name", gate_in_trailer_name, "'")
		// fmt.Println("gate_in_license", gate_in_license, "'")
		// fmt.Println("gate_in_date", gate_in_date, "'")
		// fmt.Println("gate_in_location", gate_in_location, "'")
		// fmt.Println("gate_out_trailer_name", gate_out_trailer_name, "'")
		// fmt.Println("gate_out_license", gate_out_license, "'")
		// fmt.Println("gate_out_dat", gate_out_date, "'")
		_, err = stmtHeader.Exec(
			container_no,
			size,
			typeValue,
			book_no,
			seal_no,
			customer,
			id_code,
			destination,
			vessel,
			voyage_no,
			renban,
			cy_date,
			gate_in_trailer_name,
			gate_in_license,
			gate_in_date,
			gate_in_location,
			gate_out_trailer_name,
			gate_out_license,
			gate_out_date,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

var insertUser = `INSERT INTO laden_container_user(username, password) VALUES(?, ?)`

func (d SQLDB) InsertUser(tx *sql.Tx, users []*model.User) error {
	stmtHeader, err := tx.Prepare(insertUser)
	if err != nil {
		log.Fatal(err)
	}
	defer stmtHeader.Close() // danger!
	for _, user := range users {

		var username string
		var password string

		if user.User.Valid {
			username = user.User.String
		}

		if user.User.Valid {
			password = user.Password.String
		}

		_, err = stmtHeader.Exec(username, password)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (d SQLDB) RefreshBooking(bookings []*model.Booking) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	d.ClearBooking(tx)
	d.InsertBooking(tx, bookings)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
func (d SQLDB) RefreshContainer(ladenContainers []*model.LadenContainer) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	d.ClearContainer(tx)
	d.InsertContainer(tx, ladenContainers)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
func (d SQLDB) RefreshUser(users []*model.User) error {
	tx, err := d.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	d.ClearUser(tx)
	d.InsertUser(tx, users)

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
