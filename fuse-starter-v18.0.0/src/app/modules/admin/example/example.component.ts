import { Component, ViewEncapsulation, ElementRef, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { EmailEditorComponent,UnlayerOptions,EmailEditorModule  } from '@trippete/angular-email-editor';
import {MatButtonModule} from '@angular/material/button';
import { AsyncPipe, CommonModule } from '@angular/common';
import { FormBuilder, FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import {StepperOrientation, MatStepperModule} from '@angular/material/stepper';
import { CampaignService } from './campaign.service';
import { BreakpointObserver } from '@angular/cdk/layout';
import { Observable, map, switchMap } from 'rxjs';
import {MatCheckboxModule} from '@angular/material/checkbox';
import { MatSelectModule } from '@angular/material/select';
import { MatDatepickerModule } from '@angular/material/datepicker';
import {STEPPER_GLOBAL_OPTIONS} from '@angular/cdk/stepper';
import { DateTimePickerModule } from "@syncfusion/ej2-angular-calendars";
import { MAT_DATE_LOCALE } from '@angular/material/core';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NgxMatDatetimePickerModule, NgxMatNativeDateModule, NgxMatTimepickerModule } from '@angular-material-components/datetime-picker';





@Component({
    selector     : 'example',
    standalone   : true,
    imports     : [EmailEditorModule,MatButtonModule,CommonModule,FormsModule,ReactiveFormsModule,MatFormFieldModule,
      MatInputModule,AsyncPipe,CommonModule,MatStepperModule,
      MatSelectModule,MatCheckboxModule,MatDatepickerModule,
      DateTimePickerModule,
      NgxMatDatetimePickerModule,
      NgxMatTimepickerModule,
      NgxMatNativeDateModule,
      

      ],
      
    
    templateUrl  : './example.component.html',
    styleUrls : ['./example.component.css'],
    encapsulation: ViewEncapsulation.None,
    providers: [
      {
        provide: STEPPER_GLOBAL_OPTIONS,
        
        useValue: { showError: true }
      },
      { 
        provide: MAT_DATE_LOCALE, useValue: 'en-US' 
      },
    ]
})
export class ExampleComponent implements OnInit
{
  campaignForm: FormGroup;
  additionalCampaignForm: any;
  newcampain: any;
  firstFormGroup: FormGroup;
  campaign: any;
  campaignIn: any;
  MailingList: any[] = [];
  companyID : any;
  stepperOrientation: Observable<StepperOrientation>;
  campaigncreated: boolean = false;
  campaingcreatedID: any;

  constructor(private fb: FormBuilder,
    private service: CampaignService,
    breakpointObserver: BreakpointObserver,
  ) {
    this.stepperOrientation = breakpointObserver
      .observe('(min-width: 800px)')
      .pipe(map(({matches}) => (matches ? 'horizontal' : 'vertical')));
   }

   public dateControl = new FormControl(new Date(2021,9,4,5,6,7));
  ngOnInit(): void {
     this.companyID = 'afa35ff6-4de5-4806-9a21-e0c2453d2834'; // Use dynamic companyID as needed

    this.firstFormGroup = this.fb.group({
      mailingList: [null, Validators.required],
    });
    this.campaignForm = this.fb.group({
      name: ['', Validators.required],
      subject: ['', Validators.required],
      type: ['', Validators.required],
      fromEmail: ['', [Validators.required, Validators.email]],
      fromName: ['', Validators.required],
      replyTo: ['', [Validators.required, Validators.email]],
      });
    this.additionalCampaignForm = this.fb.group({
      trackOpen: [false, Validators.required], // Default value false
  trackClick: [false, Validators.required], // Default value false
      resend: [false],
      deliveryAt: [null, Validators.required],
    });
  
    this.service.getMailingLists().subscribe({
      next: (data) => {

        if (data && data.data && data.data.items) {
          this.MailingList = data.data.items; // Assign mailing list items
         
        } else {
          console.error('Invalid response structure:', data);
          // Handle unexpected response structure if needed
        }
        console.log(data);
        
      },
      error: (error) => {
        console.error('Error fetching mailing lists:', error);
        // Handle error scenario if needed
      }
    });
    
    
  }
  
  
    
  @ViewChild(EmailEditorComponent)
  private emailEditor!: any;

  theme : string = 'dark';

  unlayerOptions!: UnlayerOptions;

  
  exportHtml() {
    this.emailEditor.exportHtml((data: any) => {
      const { design, html } = data;
      console.log('exportHtml', html);
  
      // Fetch campaign data and update in a single request sequence
      this.service.getCampaignByID( this.campaingcreatedID)
        .pipe(
          // Use switchMap to chain an update request based on successful fetch
          switchMap((response) => {
            this.campaignIn = {
              type: 'html', // Assuming this is the default type for your application
              name: response.data.name,
              subject: response.data.subject,
              html: html,
              fromEmail: response.data.fromEmail,
              fromName: response.data.fromName,
              deliveryAt: response.data.deliveryAt ? new Date(response.data.deliveryAt) : undefined,
              trackOpen: response.data.trackOpen,
              trackClick: response.data.trackClick,
              replyTo: response.data.replyTo
            };
  
            this.campaign = response.data;
            console.log('Campaign data:', this.campaign);
  
            // Update the campaign with the fetched data and extracted HTML
            return this.service.updateCampaign( this.campaingcreatedID, this.campaignIn);
          })
        )
        .subscribe(
          (updateResponse) => {
            console.log('Campaign updated successfully!', updateResponse.data);
            // Handle successful response (e.g., display success message)
          },
          (error) => {
            console.error('Error fetching or updating campaign:', error);
            // Handle errors here (e.g., display error message to the user)
          }
        );
    }, {
      cleanup: true,
    });
  }
  

  saveDesign() {
    this.emailEditor.saveDesign((design :any) => {
      console.log('saveDesign', design);
    });
  }

  loadDesign() {
    const design = {}; // Your saved design JSON
    this.emailEditor.loadDesign(design);
  }

  onEditorLoad() {
    // Handle any actions after the editor loads (optional)
  }
  
  onSubmit() {
    if (this.campaignForm.valid && this.firstFormGroup.valid) {
      // Convert to 'YYYY-MM-DD HH:mm:ss' format
      const deliveryAt = this.formatDate(this.campaignForm.value.deliveryAt);
      //const runAt = this.formatDate(this.campaignForm.value.runAt);
  
      // OR, if you want to use Unix timestamps
      // const deliveryAt = Math.floor(new Date(this.campaignForm.value.deliveryAt).getTime() / 1000);
      // const runAt = Math.floor(new Date(this.campaignForm.value.runAt).getTime() / 1000);
  
      this.newcampain = {
        mailingListId: this.firstFormGroup.value.mailingList,
        type: this.campaignForm.value.type,
        name: this.campaignForm.value.name,
        subject: this.campaignForm.value.subject,
        fromEmail: this.campaignForm.value.fromEmail,
        fromName: this.campaignForm.value.fromName,
        deliveryAt: deliveryAt, // Formatted date
        trackOpen: !!this.campaignForm.value.trackOpen, // Ensure boolean
        trackClick: !!this.campaignForm.value.trackClick, // Ensure boolean
        replyTo: this.campaignForm.value.replyTo,
        resend: this.campaignForm.value.resend
      };
      
      this.service.createCampaign(this.newcampain).subscribe(
        response => {
          console.log('Campaign created successfully!', response.data.ID);
          this.campaigncreated = true;
          this.campaingcreatedID = response.data.ID;
          // Handle successful response (e.g., display success message)
        },
        error => {
          console.error('Error creating campaign:', error);
          // Handle errors here (e.g., display error message to the user)
        }
      );
      // Handle form submission, e.g., send data to the backend
    } else {
      console.log('Form is not valid');
      // Log invalid fields for campaignForm
      Object.keys(this.campaignForm.controls).forEach(key => {
        if (this.campaignForm.controls[key].invalid) {
          console.log(`Invalid Field in campaignForm: ${key}`);
        }
      });
      // Log invalid fields for firstFormGroup
      Object.keys(this.firstFormGroup.controls).forEach(key => {
        if (this.firstFormGroup.controls[key].invalid) {
          console.log(`Invalid Field in firstFormGroup: ${key}`);
        }
      });
    }
  }
  
  // Utility function to format date to 'YYYY-MM-DD HH:mm:ss'
  formatDate(date: any): string {
    const d = new Date(date);
    const year = d.getFullYear();
    const month = ('0' + (d.getMonth() + 1)).slice(-2);
    const day = ('0' + d.getDate()).slice(-2);
    const hours = ('0' + d.getHours()).slice(-2);
    const minutes = ('0' + d.getMinutes()).slice(-2);
    const seconds = ('0' + d.getSeconds()).slice(-2);
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  }
  


   
} 
