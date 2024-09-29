import { AsyncPipe,CommonModule, CurrencyPipe, NgClass, NgFor, NgIf, NgTemplateOutlet } from '@angular/common';
import { AfterViewInit, ChangeDetectionStrategy, ChangeDetectorRef,NgModule, Component, OnDestroy, OnInit, ViewChild, ViewEncapsulation, ElementRef, TemplateRef } from '@angular/core';
import { FormControl, FormsModule, ReactiveFormsModule, UntypedFormBuilder, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxChange, MatCheckboxModule } from '@angular/material/checkbox';
import { MatOptionModule, MatRippleModule } from '@angular/material/core';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { fuseAnimations } from '@fuse/animations';
import { FuseConfirmationService } from '@fuse/services/confirmation';
import { DateTimePickerModule } from "@syncfusion/ej2-angular-calendars";
import * as htmlToImage from 'html-to-image';


import { catchError, debounceTime, map, merge, Observable, of, Subject, switchMap, takeUntil } from 'rxjs';
import { CampaignService } from '../campaign.service';
import { NgxMatDatetimePickerModule, NgxMatNativeDateModule, NgxMatTimepickerModule } from '@angular-material-components/datetime-picker';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { MatCardModule } from '@angular/material/card';
import html2canvas from 'html2canvas';
import { NavigationExtras, Router, RouterModule } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';

@Component({
  selector: 'campaignlist',
  standalone: true,
  templateUrl: './campaign-list.component.html',
  styleUrls: ['./campaign-list.component.css'],
  imports        : [CommonModule,NgIf, MatProgressBarModule, MatFormFieldModule, MatIconModule, MatInputModule, 
    FormsModule, ReactiveFormsModule, MatButtonModule, MatSortModule, NgFor, NgTemplateOutlet, MatPaginatorModule,
     NgClass, MatSlideToggleModule, MatSelectModule, MatOptionModule, MatCheckboxModule, MatRippleModule, AsyncPipe, CurrencyPipe,
    NgxMatDatetimePickerModule,
    RouterModule,
      NgxMatTimepickerModule,
      NgxMatNativeDateModule,
      DateTimePickerModule,
      MatCardModule,
  ],

})
export class CampaignListComponent implements OnInit {



  @ViewChild('htmlContent', { static: false }) htmlContent: ElementRef;
  @ViewChild('previewDialog') previewDialog: TemplateRef<any>;

  imageSrc: string = null;
  campaigns: any;
  filteredCampaigns
  selectedCampaign: any;
  flashMessage: 'success' | 'error' | null = null;
  selectedCampaignForm: UntypedFormGroup;
  isLoading: boolean = false; 
  searchControl = new FormControl('');

  constructor( private dialog: MatDialog,private router: Router,private sanitizer: DomSanitizer,private service: CampaignService, private _changeDetectorRef: ChangeDetectorRef,private fb : UntypedFormBuilder) {
    this.searchControl.valueChanges.subscribe(searchText => {
      this.filterCampaigns(searchText);
      console.log(searchText);
      
    });
   }
   filterCampaigns(searchText: string): void {
    if (!searchText) {
      // If no search text, reset the filtered campaigns to the full list
      this.filteredCampaigns = this.campaigns;
    } else {
      // Otherwise, filter the campaigns based on the search text
      this.filteredCampaigns = this.campaigns.filter(campaign =>
        campaign.name.toLowerCase().includes(searchText.toLowerCase())
      );
    }
  }

  

ngOnInit(): void {
  this.selectedCampaign=null;

  this.selectedCampaignForm = this.fb.group({
    id: [''],
    customOrder: [0],
    deliveryAt: [''],
    fromEmail: [''],
    fromName: [''],
    type: [''],
    name: [''],
    plain: [''],
    replyTo: [''],
    resend: [false],
    runAt: [''],
    signDKIM: [false],
    status: [''],
    subject: [''],
    trackClick: [true],
    trackOpen: [true],
});

  this.service.getCampaigns(1, 10).subscribe({
    next: (data) => {
      if (data && data.data && data.data.items) {
        this.campaigns = data.data.items;
        this.filteredCampaigns = this.campaigns;
        console.log(this.campaigns); // Correctly log the campaigns after they are fetched
      } else {
        console.error('Invalid response structure:', data);
        // Handle unexpected response structure if needed
      }
    },
    error: (error) => {
      console.error('Error fetching mailing lists:', error);
      // Handle error scenario if needed
    }
  });
}
navigateToDashboard(campaignId: string) {
  const navigationExtras: NavigationExtras = {
    state: { campaignId: campaignId }
  };
  this.router.navigate(['/dashboard'], navigationExtras);
}

deleteSelectedCampaign() {
  this.service.deleteCampaign(this.selectedCampaign.id)
  .subscribe(() =>
    {
      // Show a success message
        this.showFlashMessage('success');
        this.selectedCampaign = null;
        this.service.getCampaigns(1, 10).subscribe({
          next: (data) => {
            if (data && data.data && data.data.items) {
              this.campaigns = data.data.items; // Assign mailing list items
              console.log(this.campaigns); // Correctly log the campaigns after they are fetched
            } else {
              console.error('Invalid response structure:', data);
              // Handle unexpected response structure if needed
            }
          },
          error: (error) => {
            console.error('Error fetching mailing lists:', error);
            // Handle error scenario if needed
          }
        });
    }
  );
}

createCampaign() {
  this.router.navigate(['/example']);
}

openPreviewDialog(html): void {
 
  
  
  this.generateImage(html);
  this.dialog.open(this.previewDialog,{
    data: {
      html: html,
      imageSrc: this.imageSrc,
      isLoading: this.isLoading
    }
  });

}
  



deleteCampaign() {
throw new Error('Method not implemented.');
}
updateSelectedCampaign() {
  // Prepare the updated campaign object
  const campaignIn = {
    type: '',
    name: '',
    subject: '',
    html: '',
    fromEmail: '',
    fromName: '',
    deliveryAt: null, // This will be updated with the formatted date
    trackOpen: false,
    trackClick: false,
    replyTo: ''
  };

  // Assign values from form to the campaign object
  campaignIn.type = this.selectedCampaignForm.value.type;
  campaignIn.name = this.selectedCampaignForm.value.name;
  campaignIn.subject = this.selectedCampaignForm.value.subject;
  campaignIn.fromEmail = this.selectedCampaignForm.value.fromEmail;
  campaignIn.fromName = this.selectedCampaignForm.value.fromName;

  // Format the date to 'YYYY-MM-DD HH:mm:ss' or use Unix timestamp
  campaignIn.deliveryAt = this.selectedCampaignForm.value.deliveryAt 
    ? this.formatDate(this.selectedCampaignForm.value.deliveryAt) 
    : null; // Ensure null is passed if the delivery date is not set

  campaignIn.trackOpen = this.selectedCampaignForm.value.trackOpen;
  campaignIn.trackClick = this.selectedCampaignForm.value.trackClick;
  campaignIn.replyTo = this.selectedCampaignForm.value.replyTo;

  console.log(campaignIn);

  // Update the campaign using the service
  this.service.updateCampaign(this.selectedCampaign.id, campaignIn)
    .subscribe(
      () => {
        // Show a success message
        this.showFlashMessage('success');
      },
      error => {
        console.error('Error updating campaign:', error);
        // Handle errors here (e.g., display error message to the user)
      }
    );
}
// Utility function to format date to 'YYYY-MM-DDTHH:mm:ssZ' (RFC 3339 format)
formatDate(date: any): string {
  const d = new Date(date);
  const year = d.getFullYear();
  const month = ('0' + (d.getMonth() + 1)).slice(-2);
  const day = ('0' + d.getDate()).slice(-2);
  const hours = ('0' + d.getHours()).slice(-2);
  const minutes = ('0' + d.getMinutes()).slice(-2);
  const seconds = ('0' + d.getSeconds()).slice(-2);

  // Return in 'YYYY-MM-DDTHH:mm:ssZ' format (UTC time)
  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}Z`;
}


showFlashMessage(type: 'success' | 'error'): void
{
    // Show the message
    this.flashMessage = type;

    // Mark for check
    this._changeDetectorRef.markForCheck();

    // Hide it after 3 seconds
    setTimeout(() =>
    {
        this.flashMessage = null;

        // Mark for check
        this._changeDetectorRef.markForCheck();
    }, 3000);
}
closeDetails(): void
{
    this.selectedCampaign = null;
    this.isLoading = false;
    this.imageSrc = null; // Clear the image source
}
 toggleDetails(CampaignId: string): void {
    // If the Campaign is already selected...
    this.imageSrc = null; // Clear the image source
    this.isLoading = true;
    
    
    if (this.selectedCampaign && this.selectedCampaign.id === CampaignId) {
        // Close the details
        this.closeDetails();
        return;
    }

    // Get the Campaign by id
    this.service.getCampaignByID( CampaignId)
        .subscribe((Campaign) => {
            // Set the selected Campaign
            this.selectedCampaign = Campaign.data;

            // Prepare the data for the form
            const formValue = {
                id: this.selectedCampaign.id,
                customOrder: this.selectedCampaign.customOrder,
                deliveryAt: this.selectedCampaign.deliveryAt,
                fromEmail: this.selectedCampaign.fromEmail,
                fromName: this.selectedCampaign.fromName,
                html: this.selectedCampaign.html,
                type: this.selectedCampaign.type,
                name: this.selectedCampaign.name,
                plain: this.selectedCampaign.plain,
                replyTo: this.selectedCampaign.replyTo,
                resend: this.selectedCampaign.resend,
                runAt: this.selectedCampaign.runAt,
                signDKIM: this.selectedCampaign.signDKIM,
                status: this.selectedCampaign.status,
                subject: this.selectedCampaign.subject,
                trackClick: this.selectedCampaign.trackClick,
                trackOpen: this.selectedCampaign.trackOpen,
                
              };
              

            //Fill the form with the prepared data
            this.selectedCampaignForm.patchValue(formValue);
            console.log(this.selectedCampaign.id,CampaignId);
            
            

            // Mark for check
            this._changeDetectorRef.markForCheck();
        });
}

generateImage(html): void {
  
  
  if (html) {
    this.isLoading = true;
    // Create a hidden iframe
    let iframe = document.createElement('iframe');
    iframe.style.width = '100%';
    iframe.style.height = '100%';
    iframe.style.position = 'absolute';
    iframe.style.top = '-9999px';
    iframe.style.left = '-9999px';
    iframe.style.visibility = 'visible';
    document.body.appendChild(iframe);
    // Set the iframe's content
    iframe.contentDocument.open();
    iframe.contentDocument.write(html);
    iframe.contentDocument.close();

    // Wait for the iframe content to load
    iframe.onload = () => {
      html2canvas(iframe.contentDocument.body, { useCORS: true, logging: true })
        .then((canvas) => {
          this.imageSrc = canvas.toDataURL();
          console.log(this.imageSrc);
        })
        .catch((error) => {
          console.error('oops, something went wrong!', error);
        })
        .finally(() => {
          // Remove the iframe from the DOM
          document.body.removeChild(iframe);
          this.isLoading = false;
        });
    };
  }
  else{
    this.isLoading = false;
  
  }
}

sanitizeHtml(html: string): SafeHtml {

  

  return this.sanitizer.bypassSecurityTrustHtml(html);
}
    

adjustScale(iframe: HTMLIFrameElement): void {
  if (!iframe) return;
  
  const contentWidth = iframe.contentWindow.document.body.scrollWidth;
  const contentHeight = iframe.contentWindow.document.body.scrollHeight;
  
  const containerWidth = iframe.offsetWidth;
  const containerHeight = iframe.offsetHeight;
  
  const scaleX = containerWidth / contentWidth;
  const scaleY = containerHeight / contentHeight;
  
  const scale = Math.min(scaleX, scaleY);
  
  iframe.style.transform = `scale(${scale})`;
}





}
