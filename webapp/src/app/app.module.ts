import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BookDataComponent } from './book-data/book-data.component';
import { TableModule }from 'primeng/table';
import {FormsModule} from '@angular/forms';
import {InputTextModule} from 'primeng/inputtext';
import { HttpRequest, HttpClientModule } from '@angular/common/http';
import {ButtonModule} from 'primeng/button';
import { LoginComponent } from './login/login.component';
import {DialogModule,} from 'primeng/dialog';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {PasswordModule} from 'primeng/password';
import {MenubarModule} from 'primeng/menubar';
import {MenuItem} from 'primeng/api';
import {TabViewModule} from 'primeng/tabview';
import {PanelModule} from 'primeng/panel';
import { ContainerComponent } from './container/container.component';
import {GoTopButtonModule} from 'ng-go-top-button'; 
import { FilterPipe } from './filter.pipe';
import {MessagesModule} from 'primeng/messages';
import {MessageModule} from 'primeng/message';
import {DataViewModule} from 'primeng/dataview';
import {CardModule} from 'primeng/card';
 

@NgModule({
  declarations: [
    AppComponent,
    BookDataComponent,
    LoginComponent,
    ContainerComponent,
    FilterPipe
  ],
  imports: [
    BrowserAnimationsModule,
    BrowserModule,
    AppRoutingModule,
    TableModule,
    FormsModule,
    HttpClientModule,
    InputTextModule,
    ButtonModule,
    DialogModule,
    PasswordModule,
    MenubarModule, 
    TabViewModule,
    PanelModule,
    MessagesModule,
    MessageModule,
    GoTopButtonModule,
    DataViewModule,
    CardModule
    
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
