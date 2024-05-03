import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { EmailEditorModule } from '@trippete/angular-email-editor';


@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    EmailEditorModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
