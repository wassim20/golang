import { ChangeDetectorRef, Component, CUSTOM_ELEMENTS_SCHEMA, Inject, OnInit, ViewEncapsulation } from '@angular/core';
import { MailinglistService } from './mailinglist.service';
import { MatButtonModule } from '@angular/material/button';
import { AsyncPipe, CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatStepperModule } from '@angular/material/stepper';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { DateTimePickerModule } from '@syncfusion/ej2-angular-calendars';
import { NgxMatDatetimePickerModule, NgxMatNativeDateModule, NgxMatTimepickerModule } from '@angular-material-components/datetime-picker';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatIconModule } from '@angular/material/icon';
import { MAT_DIALOG_DATA, MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { AuthService } from 'app/core/auth/auth.service';
import { get } from 'lodash';
import { forkJoin } from 'rxjs';


@Component({
  selector: 'mailinglist',
  standalone: true,
  imports: [
    MatButtonModule, CommonModule, FormsModule, ReactiveFormsModule, MatFormFieldModule,
    MatInputModule, AsyncPipe, CommonModule, MatStepperModule,MatIconModule,
    MatProgressBarModule, MatInputModule, MatSelectModule, MatCheckboxModule, MatDatepickerModule,
    DateTimePickerModule, NgxMatDatetimePickerModule, NgxMatTimepickerModule, NgxMatNativeDateModule
  ],
  encapsulation: ViewEncapsulation.None,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  templateUrl: './mailinglist.component.html',
  styleUrls: ['./mailinglist.component.scss']
})
export class MailinglistComponent implements OnInit {
onSearch($event: any) {
throw new Error('Method not implemented.');
}
  mailinglists: any[] = [];
  selectedMailingListForm: UntypedFormGroup;
  selectedMailingList: any = null;
  isLoading = false;
  flashMessage: 'success' | 'error' | null = null;


  constructor(
    private service: MailinglistService,
    private fb: UntypedFormBuilder,
    private _changeDetectorRef: ChangeDetectorRef,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    // Initialize the form
    this.selectedMailingListForm = this.fb.group({
      id: [''],
      name: ['', Validators.required],
      description: [''],
      companyId: [''],
      createdByUserId: [''],
      tags: ['']
    });

    // Fetch the mailing lists
    this.fetchMailingLists();
  }

  fetchMailingLists(): void {
    this.service.getMailingLists().subscribe({
      next: (data) => {
        if (data && data.data && data.data.items) {
          
          this.mailinglists = data.data.items;
          console.log('Mailing lists:', this.mailinglists);
          this.updateUser();
        } else {
          console.error('Invalid response structure:', data);
        }
      },
      error: (error) => {
        console.error('Error fetching mailing lists:', error);
      }
    });
    
  }
  updateUser(): void {
    const updateObservables = this.mailinglists.map(mailingList =>
        this.service.getlistcreator(mailingList.created_by_user)
    );

    forkJoin(updateObservables).subscribe(results => {
        results.forEach((result, index) => {
          
            const userData = result?.data;
            console.log('User data:', userData);
            
            if (userData) {
                this.mailinglists[index].created_by_user = userData.firstname+ ' ' + userData.lastname;
            } else {
                this.mailinglists[index].created_by_user = 'Unknown';
            }
        });

        console.log('Updated mailing lists:', this.mailinglists);
        this._changeDetectorRef.markForCheck(); // Trigger change detection if necessary
    });
}

  toggleDetails(mailingListId: string): void {
    console.log('Toggling details for mailing list:', mailingListId);
    
    this.isLoading = true;

    if (this.selectedMailingList && this.selectedMailingList.id === mailingListId) {
      this.closeDetails();
      return;
    }

    this.service.getMailingListByID(mailingListId).subscribe({
      next: (response) => {
        this.selectedMailingList = response.data;

        const formValue = {
          id: this.selectedMailingList.id,
          name: this.selectedMailingList.name,
          description: this.selectedMailingList.description,
          companyId: this.selectedMailingList.companyID,
          createdByUserId: this.selectedMailingList.createdByUserID,
          tags: this.selectedMailingList.tags
        };

        this.selectedMailingListForm.patchValue(formValue);
        console.log('Mailing List ID:', this.selectedMailingList.id);

        this._changeDetectorRef.markForCheck();
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error fetching mailing list details:', error);
        this.isLoading = false;
      }
    });
  }

  closeDetails(): void {
    this.selectedMailingList = null;
    this.isLoading = false;
  }
  updateSelectedMailingList(): void {
    if (this.selectedMailingListForm.valid) {
      const updatedMailingList = this.selectedMailingListForm.value;
      console.log('Updating mailing list:', updatedMailingList);

      this.service.updateMailingList(updatedMailingList).subscribe({
        next: () => {
          this.showFlashMessage('success');
        this.fetchMailingLists();
          this.closeDetails();
        },
        error: (error) => {
          console.error('Error updating mailing list:', error);
        }
      });
    }
  }
  showFlashMessage(type: 'success' | 'error'): void {
    // Show the flash message
    this.flashMessage = type;

    // Use change detection to update the view
    this._changeDetectorRef.detectChanges();

    // Hide the flash message after 3 seconds
    setTimeout(() => {
        this.flashMessage = null;
        this._changeDetectorRef.detectChanges(); // Update the view again after hiding
    }, 3000);
}
deleteSelectedMailingList(): void {
  if (this.selectedMailingList) {
    const mailingListId = this.selectedMailingList.id;
    console.log('Deleting mailing list:', mailingListId);

    this.service.deleteMailingList(mailingListId).subscribe({
      next: () => {
        this.showFlashMessage('success');
        this.fetchMailingLists();
        this.closeDetails();
      },
      error: (error) => {
        console.error('Error deleting mailing list:', error);
      }
    });
  }
}


  createMailingList(): void {
    const dialogRef = this.dialog.open(CreateMailingListDialogComponent, {
      width: '400px'
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // Re-fetch the mailing lists to include the newly created one
        this.fetchMailingLists();
      }
    });
  }

  addContacts(id:string): void {
    
    const dialogRef = this.dialog.open(AddContactsDialogComponent, {
      width: '400px',
      data: {id:id}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // Re-fetch the mailing lists to include the newly created one
        this.fetchMailingLists();
      }
    });
  }
}



@Component({
  selector: 'app-create-mailinglist-dialog',
  standalone: true,
  imports: [MatButtonModule, MatFormFieldModule, MatInputModule, MatDialogModule,FormsModule, ReactiveFormsModule],
  template: `
  <h2 mat-dialog-title>Create Mailing List</h2>
  
  <form [formGroup]="mailingListForm" (ngSubmit)="onSubmit()" class="stepper-form">

    <!-- Name -->
    <mat-form-field appearance="outline" class="full-width">
      <mat-label>Name</mat-label>
      <input matInput formControlName="name" placeholder="Enter mailing list name" required>
      <mat-error *ngIf="mailingListForm.controls['name'].touched && mailingListForm.controls['name'].invalid">
        <ng-container *ngIf="mailingListForm.controls['name'].hasError('required')">Name is required.</ng-container>
      </mat-error>
    </mat-form-field>

    <!-- Description -->
    <mat-form-field appearance="outline" class="full-width">
      <mat-label>Description</mat-label>
      <textarea matInput formControlName="description" placeholder="Enter mailing list description"></textarea>
    </mat-form-field>

    <div mat-dialog-actions>
      <button mat-button (click)="onCancel()">Cancel</button>
      <button mat-button color="primary" type="submit" [disabled]="mailingListForm.invalid">Create</button>
    </div>
  </form>
`

})
export class CreateMailingListDialogComponent {

  mailingListForm: UntypedFormGroup;

  constructor(
    private fb: UntypedFormBuilder,
    private dialogRef: MatDialogRef<CreateMailingListDialogComponent>,
    private service: MailinglistService,
    
  ) {
    this.mailingListForm = this.fb.group({
      name: ['', Validators.required],
      description: [''],
     
    });
  }

  onSubmit(): void {
    if (this.mailingListForm.valid) {
      
      const data={
        name:this.mailingListForm.value.name,
        description:this.mailingListForm.value.description,
        // company_id: this.service['getCompanyID'](),
        // createdByUserId: this.service.getuser(),
      }
      this.service.createMailingList(this.mailingListForm.value).subscribe({
        next: () => {
          this.dialogRef.close(true);
        },
        error: (error) => {
          console.error('Error creating mailing list:', error);
        }
      });
    }
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}

@Component({
  selector: 'app-add-contacts-dialog',
  standalone: true,
  imports: [MatButtonModule,CommonModule, MatFormFieldModule, MatInputModule, MatDialogModule,FormsModule, ReactiveFormsModule],
  template: `
  <h2 mat-dialog-title>Add Contacts to Mailing List</h2>
  
  <form [formGroup]="contactForm" (ngSubmit)="onSubmit()" class="stepper-form">

    <!-- Email -->
    <mat-form-field appearance="outline" class="full-width">
      <mat-label>Email</mat-label>
      <input matInput formControlName="email" placeholder="Enter contact's email" required>
      <mat-error *ngIf="contactForm.controls['email'].touched && contactForm.controls['email'].invalid">
        <ng-container *ngIf="contactForm.controls['email'].hasError('required')">Email is required.</ng-container>
        <ng-container *ngIf="contactForm.controls['email'].hasError('email')">Please enter a valid email.</ng-container>
      </mat-error>
    </mat-form-field>

    <!-- First Name -->
    <mat-form-field appearance="outline" class="full-width">
      <mat-label>First Name</mat-label>
      <input matInput formControlName="firstname" placeholder="Enter first name" required>
      <mat-error *ngIf="contactForm.controls['firstname'].touched && contactForm.controls['firstname'].invalid">
        <ng-container *ngIf="contactForm.controls['firstname'].hasError('required')">First name is required.</ng-container>
      </mat-error>
    </mat-form-field>

    <!-- Last Name -->
    <mat-form-field appearance="outline" class="full-width">
      <mat-label>Last Name</mat-label>
      <input matInput formControlName="lastname" placeholder="Enter last name" required>
      <mat-error *ngIf="contactForm.controls['lastname'].touched && contactForm.controls['lastname'].invalid">
        <ng-container *ngIf="contactForm.controls['lastname'].hasError('required')">Last name is required.</ng-container>
      </mat-error>
    </mat-form-field>

    <!-- Phone Number -->
    <mat-form-field appearance="outline" class="full-width">
      <mat-label>Phone Number</mat-label>
      <input matInput formControlName="phoneNumber" placeholder="Enter phone number" required>
      <mat-error *ngIf="contactForm.controls['phoneNumber'].touched && contactForm.controls['phoneNumber'].invalid">
        <ng-container *ngIf="contactForm.controls['phoneNumber'].hasError('required')">Phone number is required.</ng-container>
        <ng-container *ngIf="contactForm.controls['phoneNumber'].hasError('pattern')">Please enter a valid phone number.</ng-container>
        <ng-container *ngIf="contactForm.controls['phoneNumber'].hasError('minlength') || contactForm.controls['phoneNumber'].hasError('maxlength')">
          Phone number must be 8 digits long.
        </ng-container>
      </mat-error>
    </mat-form-field>

    <div mat-dialog-actions>
      <button mat-button (click)="onCancel()">Cancel</button>
      <button mat-button color="primary" type="submit" [disabled]="contactForm.invalid">Add Contact</button>
    </div>
  </form>
`
})
export class AddContactsDialogComponent {
contactForm: UntypedFormGroup;

constructor(
  private fb: UntypedFormBuilder,
  private dialogRef: MatDialogRef<AddContactsDialogComponent>,
  private service: MailinglistService,
  @Inject(MAT_DIALOG_DATA) public data: { id: string }  // Receive the mailing list ID
) {
  this.contactForm = this.fb.group({
    email: ['', [Validators.required, Validators.email,Validators.pattern('^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,4}$' )
    ]],
    firstname: ['', Validators.required],
    lastname: ['', Validators.required],
    phoneNumber: ['', [Validators.required, Validators.pattern('^[0-9]*$'), Validators.minLength(8), Validators.maxLength(8)]],
  });
}


  onSubmit(): void {
    if (this.contactForm.valid) {
      const newContact = {
        email: this.contactForm.value.email,
        first_name: this.contactForm.value.firstname,
        last_name: this.contactForm.value.lastname,
        phone_number: this.contactForm.value.phoneNumber,
        full_name :this.contactForm.value.firstname + ' ' + this.contactForm.value.lastname
         
      };

      this.service.addContactToMailingList(newContact,this.data.id).subscribe({
        next: () => {
          this.dialogRef.close(true);
        },
        error: (error) => {
          console.error('Error adding contact:', error);
        }
      });
    }
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}
