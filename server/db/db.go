package db

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	model "tracking/model"

	_ "github.com/denisenkom/go-mssqldb"
)

type DB interface {
	GetContainerByBookingNumber(bookingNumber string) ([]*model.LadenContainer, error)
	GetContainerContainerNumber(containerNumber string) ([]*model.LadenContainer, error)
	GetBookingByBookingNumber(bookingNumber string) (*model.Booking, error)
	GetBookingByContainerNumber(containerNumber string) (*model.Booking, error)
	Login(user, password string) (*model.User, error)
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
WHERE book_no = @BOOKNO
`

//

var queryBookingType = `
SELECT book_no,size,type,quantity,available,total_out
FROM booking_container_type
WHERE book_no = @BOOKNO
`

var queryBookingDetail = `
SELECT book_no,no,container_no,size,type,seal_no,trailer_name,license,gate_out_date
FROM booking_container_detail
WHERE book_no = @BOOKNO
`

var queryBookingbyContainer = `
SELECT TOP 1 d.book_no,d.no,d.container_no,d.size,d.type,d.seal_no,d.trailer_name,d.license,d.gate_out_date
FROM booking_container_detail d INNER JOIN booking_header h on d.book_no = h.book_no 
WHERE  container_no = @CONTAINER
ORDER BY create_date desc
`

/**

 */
func (d SQLDB) GetBookingByBookingNumber(bookingNumber string) (*model.Booking, error) {

	var booking *model.Booking

	booking = &model.Booking{}

	stmt, err := d.db.Prepare(queryBookingHeader)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("BOOKNO", bookingNumber))

	err = row.Scan(&booking.BookNo, &booking.Operator, &booking.Customer, &booking.VoyageNo, &booking.Destination, &booking.VesselName, &booking.PickupDate, &booking.GoodsDescription, &booking.Remark)
	if err != nil {

		if err != nil {
			return booking, nil
		}
	}

	rows, err := d.db.Query(queryBookingType, sql.Named("BOOKNO", bookingNumber))
	if err != nil {
		return booking, nil
		// log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		containerType := &model.BookingContainerType{}
		err := rows.Scan(&containerType.BookNo, &containerType.Size, &containerType.Type, &containerType.Quantity, &containerType.Available, &containerType.TotalOut)
		if err != nil {
			// log.Fatal(err)
			return booking, nil
		} else {
			booking.BookingContainerTypes = append(booking.BookingContainerTypes, *containerType)

		}
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return booking, nil
	}

	rows, err = d.db.Query(queryBookingDetail, sql.Named("BOOKNO", bookingNumber))
	if err != nil {
		// log.Fatal(err)
		return booking, nil
	}
	defer rows.Close()

	for rows.Next() {
		containerDetail := &model.BookingContainerDetail{}
		err := rows.Scan(&containerDetail.BookNo, &containerDetail.No, &containerDetail.ContainerNo, &containerDetail.Size, &containerDetail.Type, &containerDetail.SealNo, &containerDetail.TrailerName, &containerDetail.License, &containerDetail.GateOutDate)
		if err != nil {
			// log.Fatal(err)
			return booking, nil
		} else {
			booking.BookingContainerDetails = append(booking.BookingContainerDetails, *containerDetail)
		}
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return booking, nil
	}

	return booking, nil
}

func (d SQLDB) GetBookingByContainerNumber(containerNumber string) (*model.Booking, error) {

	var booking *model.Booking

	booking = &model.Booking{}

	rowBookingNumber := d.db.QueryRow(queryBookingbyContainer, sql.Named("CONTAINER", containerNumber))

	containerDetail := &model.BookingContainerDetail{}
	err := rowBookingNumber.Scan(&containerDetail.BookNo, &containerDetail.No, &containerDetail.ContainerNo, &containerDetail.Size, &containerDetail.Type, &containerDetail.SealNo, &containerDetail.TrailerName, &containerDetail.License, &containerDetail.GateOutDate)
	if err != nil {
		return booking, nil
	}

	stmt, err := d.db.Prepare(queryBookingHeader)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(sql.Named("BOOKNO", containerDetail.BookNo))

	err = row.Scan(&booking.BookNo, &booking.Operator, &booking.Customer, &booking.VoyageNo, &booking.Destination, &booking.VesselName, &booking.PickupDate, &booking.GoodsDescription, &booking.Remark)
	if err != nil {

		return booking, nil
	}

	rows, err := d.db.Query(queryBookingType, sql.Named("BOOKNO", containerDetail.BookNo))
	if err != nil {
		return booking, nil
		// log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		containerType := &model.BookingContainerType{}
		err := rows.Scan(&containerType.BookNo, &containerType.Size, &containerType.Type, &containerType.Quantity, &containerType.Available, &containerType.TotalOut)
		if err != nil {
			// log.Fatal(err)
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

SELECT [container_no]
           ,[size]
           ,[type]
           ,[book_no]
           ,[seal_no]
           ,[customer]
           ,[id_code]
           ,[destination]
           ,[vessel]
           ,[voyage_no]
           ,[renban]
           ,[cy_date]
           ,[gate_in_trailer_name]
           ,[gate_in_license]
           ,[gate_in_date]
           ,[gate_in_location]
           ,[gate_out_trailer_name]
           ,[gate_out_license]
		   ,[gate_out_date] 
FROM dbo.laden_container
WHERE [book_no] = @BOOKNO		 
ORDER BY gate_in_date DESC  
`

/*

 */
func (d SQLDB) GetContainerByBookingNumber(bookingNumber string) ([]*model.LadenContainer, error) {

	var ladenContainers []*model.LadenContainer

	rows, err := d.db.Query(queryContainerByBookingNo, sql.Named("BOOKNO", bookingNumber))
	if err != nil {
		// log.Fatal(err)
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
			// log.Fatal(err)
			return ladenContainers, nil
		} else {
			ladenContainers = append(ladenContainers, ladenContainer)
		}
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return ladenContainers, nil
	}

	return ladenContainers, nil
}

var queryContainerByContainerNo = ` 
SELECT top 1 [container_no]
           ,[size]
           ,[type]
           ,[book_no]
           ,[seal_no]
           ,[customer]
           ,[id_code]
           ,[destination]
           ,[vessel]
           ,[voyage_no]
           ,[renban]
           ,[cy_date]
           ,[gate_in_trailer_name]
           ,[gate_in_license]
           ,[gate_in_date]
           ,[gate_in_location]
           ,[gate_out_trailer_name]
           ,[gate_out_license]
		   ,[gate_out_date] 
FROM dbo.laden_container
WHERE [container_no] = @CONTAINERNO  
ORDER BY gate_in_date DESC
 
`

func (d SQLDB) GetContainerContainerNumber(containerNumber string) ([]*model.LadenContainer, error) {

	var ladenContainers []*model.LadenContainer

	rows, err := d.db.Query(queryContainerByContainerNo, sql.Named("CONTAINERNO", containerNumber))
	if err != nil {
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
			// log.Fatal(err)
			return ladenContainers, nil
		} else {
			ladenContainers = append(ladenContainers, ladenContainer)
		}
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return ladenContainers, nil
	}

	return ladenContainers, nil
}

var queryLogin = `
SELECT username
FROM dbo.laden_container_user
WHERE username = @USER AND password = @PASSWORD
`

func (d SQLDB) Login(userName, password string) (*model.User, error) {

	var user *model.User

	rows, err := d.db.Query(queryLogin, sql.Named("USER", userName), sql.Named("PASSWORD", password))
	if err != nil {
		fmt.Printf("error %s\n", err)
		return user, nil
		// log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		user = &model.User{}
		err := rows.Scan(&user.User)

		fmt.Printf("scan %s\n", user)
		if err != nil {
			// log.Fatal(err)
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
