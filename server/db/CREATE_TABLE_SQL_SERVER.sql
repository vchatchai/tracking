USE [tracking]
GO
/****** Object:  Table [dbo].[booking_header]    Script Date: 08/15/2020 00:20:14 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[booking_header](
	[book_no] [char](30) NOT NULL,
	[operator] [char](10) NULL,
	[customer] [char](180) NULL,
	[yoyage_no] [char](10) NULL,
	[destination] [char](100) NULL,
	[vessel_name] [char](80) NULL,
	[pickup_date] [datetime] NULL,
	[goods_description] [text] NULL,
	[remark] [text] NULL,
	[create_date] [datetime] NULL
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
/****** Object:  Table [dbo].[booking_container_type]    Script Date: 08/15/2020 00:20:14 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[booking_container_type](
	[book_no] [char](30) NOT NULL,
	[size] [char](2) NULL,
	[type] [char](2) NULL,
	[quantity] [smallint] NULL,
	[available] [smallint] NULL,
	[total_out] [smallint] NULL,
) ON [PRIMARY]
GO
/****** Object:  Table [dbo].[booking_container_detail]    Script Date: 08/15/2020 00:20:14 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[booking_container_detail](
	[book_no] [char](30) NOT NULL,
	[no] [char](4) NULL,
	[container_no] [char](7w) NULL,
	[size] [nchar](2) NULL,
	[type] [nchar](2) NULL,
	[seal_no] [nchar](20) NULL,
	[trailer_name] [char](50) NULL,
	[license] [char](7) NULL,
	[gate_out_date] [datetime] NULL
) ON [PRIMARY]
GO





----------------
 
/****** Object:  Table [dbo].[laden_container]    Script Date: 08/20/2020 08:43:10 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

SET ANSI_PADDING ON
GO

CREATE TABLE [dbo].[laden_container](
	[container_no] [char](11) NOT NULL,
	[size] [char](2) NULL,
	[type] [char](2) NULL,
	[book_no] [char](30) NULL,
	[seal_no] [char](20) NULL,
	[customer] [char](100) NULL,
	[id_code] [char](20) NULL,
	[destination] [char](100) NULL,
	[vessel] [char](100) NULL,
	[voyage_no] [char](20) NULL,
	[renban] [char](20) NULL,
	[cy_date] [datetime] NULL,
	[gate_in_trailer_name] [char](100) NULL,
	[gate_in_license] [char](7) NULL,
	[gate_in_date] [datetime] NULL,
	[gate_in_location] [char](10) NULL,
	[gate_out_trailer_name] [char](100) NULL,
	[gate_out_license] [char](7) NULL,
	[gate_out_date] [datetime] NULL
) ON [PRIMARY]

GO

SET ANSI_PADDING OFF
GO
 
/****** Object:  Table [dbo].[laden_container_user]    Script Date: 08/20/2020 08:50:41 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[laden_container_user](
	[username] [nvarchar](50) NOT NULL,
	[password] [nvarchar](50) NULL
) ON [PRIMARY]

GO



