USE tracking;
CREATE TABLE booking_header(
	book_no char(30) NOT NULL,
	operator char(10) NULL,
	customer char(180) NULL,
	yoyage_no char(10) NULL,
	vessel_name char(80) NULL,
	pickup_date datetime NULL,
	goods_description text NULL,
	remark text NULL
) ON PRIMARY ;
CREATE TABLE booking_container_type(
	book_no char(30) NOT NULL,
	size char(2) NULL,
	type char(2) NULL,
	quantity smallint NULL,
	available smallint NULL,
	total_out smallint NULL,
) ON PRIMARY;

CREATE TABLE booking_container_detail(
	book_no char(30) NOT NULL,
	no char(4) NULL,
	container_no char(7) NULL,
	size nchar(2) NULL,
	type nchar(2) NULL,
	seal_no nchar(20) NULL,
	trailer_name char(50) NULL,
	license char(7) NULL,
	gate_out_date datetime NULL
) ON PRIMARY;
