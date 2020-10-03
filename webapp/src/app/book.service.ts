import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from 'src/environments/environment';



export interface Book {
  book_no
  operator
  customer
  yoyage_no
  destination
  vessel_name
  pickup_date
  goods_description
  remark
  BookingContainerTypes: BookingContainerTypes
  bookingContainerDetails: BookingContainerDetails
}

interface BookingContainerTypes extends Array<BookingContainerType> { }
interface BookingContainerDetails extends Array<BookingContainerDetail> { }


export interface BookingContainerType {
  book_no
  size
  type
  quantity
  available
}

export interface BookingContainerDetail {
  no
  type
  seal_no
  trailer_name
  license
  gate_out_date: Date
}



export interface Container {
  container_no
  size
  type
  booking_no
  seal_no
  customer
  ld_code
  origin
  destination
  vessel
  voyage_no
  renban
  cy_date
  gate_in_trailer_name
  gate_in_license
  gate_in_date
  gate_in_location
  gate_out_trailer_name
  gate_out_license
  gate_out_date
}

export interface State  {
  searchBooking: boolean 
}

@Injectable({
  providedIn: 'root'
})
export class BookService {
  constructor(private http: HttpClient) { }




  // bookingUrl = "publish/booking/"
  getBookByContainerNumber(containerNumber) {
    // return this.http.get<any>('assets/books.json')
    // .toPromise()
    // .then(res => <Book[]> res.data)
    // .then(data => {return data;});
    console.log("containerNumber", containerNumber)
    const params = new HttpParams()
      .set('containerNumber', containerNumber);


    return this.http.get(environment.bookingUrl, { params })
  }
  getBookByBookingNumber(bookingNumber) {
    // return this.http.get<any>('assets/books.json')
    // .toPromise()
    // .then(res => <Book[]> res.data)
    // .then(data => {return data;});
    console.log("bookingNumber", bookingNumber)
    const params = new HttpParams()
      .set('bookingNumber', bookingNumber);


    return this.http.get(environment.bookingUrl, { params })
  }
  getContainerByBookingNumber(bookingNumber) {
    console.log(bookingNumber)
    const params = new HttpParams()
      .set('bookingNumber', bookingNumber);

    return this.http.get(environment.containerUrl, { params })
  }

  getContainerByContainerNumber(containerNumber) {
    console.log(containerNumber)
    const params = new HttpParams()
      .set('containerNumber', containerNumber);

    return this.http.get(environment.containerUrl, { params })
  }
  login(username, password) {
    console.log(username)
    const params = new HttpParams()
      .set('username', username)
      .set('password', password)
      ;

    return this.http.get(environment.loginUrl, { params })
  }



  // getPosts(data:Post[] ) { 
  //   this.http.get(this.url).subscribe((data) => {


  //     console.warn(data)
  //   })
  // return this.http.get<any>('https://jsonplaceholder.typicode.com/posts')
  // .toPromise()
  // .then(res => <Post[]> res.data)
  // .then(data => {return data;});
  // }
}
