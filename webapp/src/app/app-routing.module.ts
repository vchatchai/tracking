import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { BookDataComponent } from './book-data/book-data.component';
import { AppComponent } from './app.component';
import {LoginComponent} from './login/login.component'


const routes: Routes = [ 
  // {path: '', component: AppComponent},
  
  // {path: 'track/:type/:id', component: BookDataComponent},  
  {path: '', component: BookDataComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
