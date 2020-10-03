import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { BookService, Book, State } from '../book.service';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { FilterPipe } from '../filter.pipe'
import { formatDate } from '@angular/common';
@Component({
  selector: 'app-book-data',
  templateUrl: './book-data.component.html',
  styleUrls: ['./book-data.component.css']
})
export class BookDataComponent implements OnInit {

  book: Book;
  bookTemp: Book;
  cols: any[];


  @Output() showData: EventEmitter<any> = new EventEmitter<any>();

  @Input()
  bookingNumber: string = "";
  containerNumber: string = "";
  showResult = false;
  isEmpty = false;

  inputSize = 30;
  emptyMessage;
  searchText: string = "";
  state: State = { searchBooking: false };

  constructor(private bookService: BookService, private router: ActivatedRoute, private route: Router) {
    this.router.queryParams.subscribe(params => {
      console.log(params)
      if (params['booking'] != null) {
        this.bookingNumber = params['booking'];
        this.searchByBookingNumber(null);
      }

      if (params['container'] != null) {
        this.containerNumber = params['container'];
        this.searchBookByContainerNumber(null);

      }
    });

  }

  ngOnInit(): void {



    this.cols = [
      { field: 'no', header: 'No' },
      { field: 'container_no', header: 'Container No' },
      { field: 'size', header: 'Size/Type' },
      { field: 'seal_no', header: 'Seal No' },
      { field: 'trailer_name', header: 'Trailer' },
      { field: 'license', header: 'License' },
      { field: 'gate_out_date', header: 'Gate Out Date' },
    ];



  }
  clearContainer(event) {
    this.containerNumber = "";
    this.showResult = false;
    this.searchText = "";
    // this.route.navigate([`/track/container/ `]);
  }
  clearBooking(event) {
    this.bookingNumber = "";
    this.showResult = false;
    this.searchText = "";
    // this.route.navigate([`/track/booking/ `]);
  }

  countLine = 0;

  countValue(rowData) {
    console.log(rowData)
    if (rowData.count == undefined) {
      rowData.count = 0
    }


    return rowData.count + 1
  }



  searchByBookingNumber(event) {
    // alert(this.bookingNumber);
    // this.router.snapshot.params.id = this.bookingNumber;
    this.searchText = "";
    // this.route.navigate([`/track/booking/${this.bookingNumber}`]);
    if (this.bookingNumber != null && this.bookingNumber.length > 0) {


      this.book = {} as Book;
      this.bookService.getBookByBookingNumber(this.bookingNumber).subscribe((data) => {


        // data = {
        //   book_no: "booking123",
        //   operator: "operator",
        //   customer: "customer01",
        //   yoyage_no: "voyageNo",
        //   vessel_name: "vesselName",
        //   pickup_date: "2020-08-15T22:22:20.671832922+07:00",
        //   goods_description: "Goood Description",
        //   remark: "remark",
        //   bookingContainerTypes: [
        //     {
        //       book_no: "booking123",
        //       size: 10,
        //       type: "MB3",bx
        //       quantity: 10,
        //       available: 2
        //     }
        //   ],
        //   bookingContainerDetails: [
        //     {
        //       no: 1,
        //       type: "MB3",
        //       seal_no: "sealNo",
        //       trailer_name: "trialerName",
        //       license: "license",
        //       gate_out_date: "2020-08-15T22:22:20.671835499+07:00"
        //     }
        //   ]

        // <Post[]>this.post
        this.book = <Book>data;
        this.bookTemp = <Book>JSON.parse(JSON.stringify(data));

        this.globalSearch()

        this.state.searchBooking = true;

        console.log(this.book)

        this.isEmpty = this.book['book_no'] == null;
        this.showResult = true;
        console.log('this.isEmpty', this.isEmpty)
        console.log('this.showResult', this.showResult)
        if (this.isEmpty) {
          this.emptyMessage = "Record not found !";
        }
      });


    } else {
      this.isEmpty = true;
      this.showResult = true;
      this.emptyMessage = "Please fill in the Booking No.";
    }


    this.showData.emit();


  }


  searchBookByContainerNumber(event) {
    this.searchText = "";
    if (this.containerNumber != null && this.containerNumber.length > 0) {

      this.book = {} as Book;
      this.bookService.getBookByContainerNumber(this.containerNumber).subscribe((data) => {


        // data = {
        //   book_no: "booking123",
        //   operator: "operator",
        //   customer: "customer01",
        //   yoyage_no: "voyageNo",
        //   vessel_name: "vesselName",
        //   pickup_date: "2020-08-15T22:22:20.671832922+07:00",
        //   goods_description: "Goood Description",
        //   remark: "remark",
        //   bookingContainerTypes: [
        //     {
        //       book_no: "booking123",
        //       size: 10,
        //       type: "MB3",bx
        //       quantity: 10,
        //       available: 2
        //     }
        //   ],
        //   bookingContainerDetails: [
        //     {
        //       no: 1,
        //       type: "MB3",
        //       seal_no: "sealNo",
        //       trailer_name: "trialerName",
        //       license: "license",
        //       gate_out_date: "2020-08-15T22:22:20.671835499+07:00"
        //     }
        //   ]

        // <Post[]>this.post
        this.book = <Book>data;
        this.bookTemp = <Book>JSON.parse(JSON.stringify(data));
        this.globalSearch()
        this.state.searchBooking = false;

        console.log(this.book)
        this.isEmpty = this.book['book_no'] == null;
        this.showResult = true;
 
        if (this.isEmpty) {
          this.emptyMessage = "Record not found !";
        }
      });

    } else {
      this.isEmpty = true;
      this.showResult = true;
      this.emptyMessage = "Please fill in the Container No.";
    }



    this.showData.emit();


  }


  globalSearch() {

    if(this.bookTemp.bookingContainerDetails == null) { 
      return
    }

    this.book.bookingContainerDetails = this.bookTemp.bookingContainerDetails.filter((item, index, items) => {

      this.state.searchBooking = true
      if (this.searchText === undefined || !this.searchText) {
        return true
      }

      this.state.searchBooking = false

      let searchText = this.searchText.toLowerCase();
      let result: boolean = false;
      for (let key of Object.keys(item)) {
        if (key == 'no') {
        } else if (key != 'size' && key != 'gate_out_date' && key != 'gate_in_date') { 
          if (item[key] != null) {
            if (!result) {
              result = item[key].toLocaleLowerCase().includes(searchText);
            }
          }

        } else if (key == 'size') {
          if (item[key] != null) {
            let value = `${item[key]}/${item['type']}`
            if (!result) {
              result = value.toLocaleLowerCase().includes(searchText)
            }
          }
        } else if (key == 'gate_out_date') {
          if (!result) {
            result = formatDate(item[key], 'dd/MM/yyyy', 'en-US').toLocaleLowerCase().includes(searchText);
          }
        } else if (key == 'gate_in_date') {
          if (!result) {
            result = formatDate(item[key], 'dd/MM/yyyy', 'en-US').toLocaleLowerCase().includes(searchText);
          }
        }


      }
      return result

    }).sort((n1,n2) =>   new Date(n1.gate_out_date).getTime()   -  new Date(n2.gate_out_date).getTime() )
 
     


  }




}
