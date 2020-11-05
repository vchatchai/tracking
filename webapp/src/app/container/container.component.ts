import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core';
import { BookService, Container } from '../book.service';

@Component({
  selector: 'app-container',
  templateUrl: './container.component.html',
  styleUrls: ['./container.component.css']
})
export class ContainerComponent implements OnInit {

  constructor(private bookService: BookService) { }

  ngOnInit(): void {
    this.cols = [
      { field: 'ld_code', header: 'LDCode' },
      { field: 'customer', header: 'Customer name' },
      { field: 'destination', header: 'Destination/ Origin country' },
      { field: 'booking_no', header: 'Booking No/BL No' },
      { field: 'container_no', header: 'Container No.' },
      { field: 'seal_no', header: 'Seal No.' },
      { field: 'type', header: 'Size/Type' },
      { field: 'vessel', header: 'Vessel' },
      { field: 'voyage_no', header: 'VoyageNo' },
      { field: 'gate_in_date', header: 'Gate in' },
      { field: 'gate_out_date', header: 'Gate out' },
      { field: 'DIY', header: 'DIY' },
      { field: 'gate_in_location', header: 'Location' },
    ];




  }



  @Output() showData: EventEmitter<any> = new EventEmitter<any>();

  @Input()
  bookingNumber: string = "";
  containerNumber: string = "";

  cols: any[];
  showResult = false;
  isEmpty = false;
  containers: Container[] = [];
  emptyMessage;
  searchText;
  clearContainer(event) {
    this.containerNumber = null;
    this.showResult = false;
    this.searchText = "";
  }
  clearBooking(event) {
    this.bookingNumber = null;
    this.showResult = false;
    this.searchText = "";
  }

  searchByBooking(event) {
    this.searchText = "";
    if ( this.bookingNumber != null && this.bookingNumber.length > 0) {
      this.bookService.getContainerByBookingNumber(this.bookingNumber).subscribe((data) => {
        console.warn(data);
        this.containers = <Container[]>data;
        this.showResult = true;
        this.isEmpty = this.containers == null;
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
  searchByContainer(event) { 
    this.searchText = "";
    if (this.containerNumber != null && this.containerNumber.length > 0) {
      this.bookService.getContainerByContainerNumber(this.containerNumber).subscribe((data) => {
        console.warn(data);
        this.containers = <Container[]>data;
        this.showResult = true;
        this.isEmpty = this.containers == null;
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

  calculateDate(date1: Date, date2: Date) {
    
    console.log('date1',date1)

    console.log('date2',date2)
    if(date1 == null) {
      
      date1 = new Date();
    }

    if(date2 == null) {
      date2 = new Date();
    }
    
    date1 = new Date(date1);
    
    date2 = new Date(date2);
    
        console.log('nextdate1',date1)
    
        console.log('nextdate2',date2)
 
    // var diffDays: any =  new Date(date1.getTime() - date2.getTime()).;  // Math.floor((Number(date1) - Number(date2)) / (1000 * 60 * 60 * 24));
    var diffDays: any = Math.floor((Date.UTC(date1.getFullYear(), date1.getMonth(), date1.getDate()) - Date.UTC(date2.getFullYear(), date2.getMonth(), date2.getDate()) ) /(1000 * 60 * 60 * 24));
    diffDays = diffDays + 1;
    console.log(diffDays)
    return diffDays;
  }
}
