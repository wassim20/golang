import { ChangeDetectorRef, Component, CUSTOM_ELEMENTS_SCHEMA, Inject, OnInit, ViewEncapsulation } from '@angular/core';
import { MailinglistService } from './mailinglist.service';
import { MatButtonModule } from '@angular/material/button';
import { AsyncPipe, CommonModule, CurrencyPipe, NgClass, NgFor, NgIf } from '@angular/common';
import { FormControl, FormsModule, ReactiveFormsModule, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
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
import { MatOptionModule, MatRippleModule } from '@angular/material/core';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import * as XLSX from 'xlsx';
import { MatTabsModule } from '@angular/material/tabs';
import { TranslocoModule } from '@ngneat/transloco';
import { MatMenuModule } from '@angular/material/menu';


@Component({
  selector: 'mailinglist',
  standalone: true,
  imports: [
    MatButtonModule, CommonModule, FormsModule, ReactiveFormsModule, MatFormFieldModule,
    MatInputModule, AsyncPipe, CommonModule, MatStepperModule,MatIconModule,
    MatProgressBarModule, MatInputModule, MatSelectModule, MatCheckboxModule, MatDatepickerModule,
    DateTimePickerModule, NgxMatDatetimePickerModule, NgxMatTimepickerModule, NgxMatNativeDateModule,
    MatSnackBarModule,
  ],
  encapsulation: ViewEncapsulation.None,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  templateUrl: './mailinglist.component.html',
  styleUrls: ['./mailinglist.component.scss']
})
export class MailinglistComponent implements OnInit {

  mailinglists: any[] = [];
  filteredMailingLists = this.mailinglists;
  selectedMailingListForm: UntypedFormGroup;
  selectedMailingList: any = null;
  isLoading = false;
  flashMessage: 'success' | 'error' | null = null;
  mailingListSearchControl = new FormControl('');


  constructor(
    private snackBar: MatSnackBar,
    private service: MailinglistService,
    private fb: UntypedFormBuilder,
    private _changeDetectorRef: ChangeDetectorRef,
    private dialog: MatDialog
  ) {
    this.mailingListSearchControl.valueChanges.subscribe(searchText => {
      this.filterMailingLists(searchText);
    });
  }
  filterMailingLists(searchText: string): void {
    if (!searchText) {
      // If no search text, reset the filtered mailing lists to the full list
      this.filteredMailingLists = this.mailinglists;
    } else {
      // Otherwise, filter the mailing lists based on the search text
      this.filteredMailingLists = this.mailinglists.filter(mailingList =>
        mailingList.name.toLowerCase().includes(searchText.toLowerCase())
      );
    }
  }

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
          this.filteredMailingLists = this.mailinglists;
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
            //console.log('User data:', userData);
            
            if (userData) {
                this.mailinglists[index].created_by_user = userData.firstname+ ' ' + userData.lastname;
            } else {
                this.mailinglists[index].created_by_user = 'Unknown';
            }
        });

        //console.log('Updated mailing lists:', this.mailinglists);
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
      width: '600px'
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        // Show a success message
        this.snackBar.open('Mailing list created successfully!', 'Close', {
          duration: 3000,
          verticalPosition: 'top',
        });

        // Re-fetch the mailing lists to include the newly created one
        this.fetchMailingLists();
      }
    });
  }

  addContacts(id:string): void {
    
    const dialogRef = this.dialog.open(AddContactsDialogComponent, {
      width: '800px',height:'600px',
      data: {id:id}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result && result.success) {
        console.log(result);
        this.snackBar.open('Contact added successfully!', 'Close', {
          duration: 3000,
          verticalPosition: 'top',
        });
        // Re-fetch the mailing lists to include the newly created one
        this.fetchMailingLists();
      } else if (result.error.ConstraintName==="contacts_email_key") {
        console.log(result);
        
        this.snackBar.open(`Error: Email already exists in the Mailing list`, 'Close', {
          duration: 5000,
          verticalPosition: 'top',
        });
      } else {
        this.snackBar.open('Contact not added! Unkown error', 'Close', {
          duration: 3000,
          verticalPosition: 'top',
        });
      }
    });
  }
}


@Component({
  selector: 'app-create-mailinglist-dialog',
  standalone: true,
  styleUrls: ['./mailinglist.component.scss'],
  imports: [
    MatButtonModule, MatFormFieldModule, MatInputModule, MatDialogModule,
    FormsModule, ReactiveFormsModule, MatSelectModule, MatOptionModule,
    MatIconModule, CommonModule
  ],
  template: `
    <h2 mat-dialog-title class="title">Create Mailing List</h2>
    <div mat-dialog-content>
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
        <mat-form-field appearance="outline" class="double-full-width">
          <mat-label>Description</mat-label>
          <textarea matInput formControlName="description" placeholder="Enter mailing list description" rows="4"></textarea>
        </mat-form-field>
      </form>
    </div>
    <div mat-dialog-actions class="stepper-buttons">
      <button mat-flat-button color="primary" type="submit" (click)="onSubmit()" class="save-button" [disabled]="mailingListForm.invalid">Create</button>
      <mat-error *ngIf="mailingListForm.invalid && mailingListForm.touched" class="form-error">
        Please check for form errors and fill out all required fields.
      </mat-error>
      <button mat-button class="cancel-button" (click)="onCancel()">Cancel</button>
    </div>
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
      const data = {
        name: this.mailingListForm.value.name,
        description: this.mailingListForm.value.description,
      };
      console.log('Submitting:', data); // Log submission data

      this.service.createMailingList(data).subscribe({
        next: () => {
          console.log('Mailing list created successfully'); // Confirm successful creation
          this.dialogRef.close(true);
        },
        error: (error) => {
          console.error('Error creating mailing list:', error);
        }
      });
    } else {
      console.log('Form is invalid:', this.mailingListForm.errors); // Log form errors
    }
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}


@Component({
  selector: 'app-add-contacts-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    MatButtonModule,
    MatTabsModule,
    MatStepperModule,
    MatFormFieldModule,
    MatInputModule,
    MatDialogModule,
    TranslocoModule, MatRippleModule, MatMenuModule,
     NgFor, NgIf,  NgClass, CurrencyPipe,
     MatSelectModule,MatIconModule
  ],
  template: `<mat-tab-group>

  <!-- Tab 1: Manually Add Contact -->
  <mat-tab label="Manually Add Contact">
    <div class="settings-section">
      <h2 class="title">Add a New Contact</h2>
      <div class="shadow-lg overflow-hidden">
        <form class="mailinglist-form" [formGroup]="contactForm" (ngSubmit)="onSubmit()">
          <div class="form-container">
            <!-- First Name -->
            <mat-form-field class="input-group half-width">
              <mat-label>First Name</mat-label>
              <input matInput formControlName="firstname" required>
              <mat-error *ngIf="contactForm.controls['firstname'].touched && contactForm.controls['firstname'].invalid">
                Please enter a valid first name.
              </mat-error>
            </mat-form-field>

            <!-- Last Name -->
            <mat-form-field class="input-group half-width">
              <mat-label>Last Name</mat-label>
              <input matInput formControlName="lastname" required>
              <mat-error *ngIf="contactForm.controls['lastname'].touched && contactForm.controls['lastname'].invalid">
                Please enter a valid last name.
              </mat-error>
            </mat-form-field>

            <!-- Email -->
            <mat-form-field class="input-group full-width">
              <mat-label>Email</mat-label>
              <input matInput formControlName="email" type="email" required>
              <mat-error *ngIf="contactForm.controls['email'].touched && contactForm.controls['email'].invalid">
                Please enter a valid email.
              </mat-error>
            </mat-form-field>

            <!-- Phone Number -->
            <mat-form-field class="input-group full-width">
              <mat-label>Phone Number</mat-label>
              <input matInput formControlName="phoneNumber" required>
              <mat-error *ngIf="contactForm.controls['phoneNumber'].touched && contactForm.controls['phoneNumber'].invalid">
                Please enter a valid phone number.
              </mat-error>
            </mat-form-field>
          </div>

          <div mat-dialog-actions class="stepper-buttons">
            <button mat-flat-button class="cancel-button" color="primary" type="submit" [disabled]="contactForm.invalid">
              Add Contact
            </button>
          </div>
        </form>
      </div>
    </div>
  </mat-tab>

  <!-- Tab 2: Bulk Upload Contacts -->
  <mat-tab label="Bulk Upload Contacts">
    <div class="settings-section">
      <h2 class="title">Upload Multiple Contacts</h2>
      <div class="shadow-lg overflow-hidden">
        <form (ngSubmit)="onSubmitBulk()">
          <!-- File Input for Bulk Upload without mat-form-field -->
          <div class="input-group full-width">
            <label>Upload Contacts</label>
            <input type="file" (change)="onFileChange($event)" accept=".csv, .xlsx, .xls" required />
            <div *ngIf="fileError" class="mat-error">{{ fileError }}</div>
          </div>


          <div mat-dialog-actions class="stepper-buttons">
            <button mat-flat-button class="cancel-button" color="primary" type="submit" [disabled]="!fileData">
              Upload Contacts
            </button>
          </div>
        </form>
      </div>
    </div>
  </mat-tab>

</mat-tab-group>
`
})
export class AddContactsDialogComponent {
contactForm: UntypedFormGroup;
fileData: any[] = []; // Holds the processed file data
fileError: string | null = null;
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
          this.dialogRef.close({ success: true });
        },
        error: (error) => {
          console.error('Error adding contact:', error.error.data);
          this.dialogRef.close({ success: false, error: error.error.data });
        }
      });
    }
  }
  // Handle file change
  onFileChange(event: any): void {
    const file = event.target.files[0];
    if (file) {
      const fileExtension = file.name.split('.').pop();
      if (fileExtension === 'csv') {
        this.readCSVFile(file);
      } else if (fileExtension === 'xlsx' || fileExtension === 'xls') {
        this.readExcelFile(file);
      } else {
        this.fileError = 'Please upload a valid CSV or Excel file.';
      }
    }
  }

  

  readExcelFile(file: File): void {
    const reader = new FileReader();
    reader.onload = (e: any) => {
      const binaryStr = e.target.result;
      const wb = XLSX.read(binaryStr, { type: 'binary' });
      const ws = wb.Sheets[wb.SheetNames[0]];
      const data = XLSX.utils.sheet_to_json(ws, { header: 1 });
      this.processData(data);
    };
    reader.readAsBinaryString(file);
  }

  readCSVFile(file: File): void {
    const reader = new FileReader();
    reader.onload = (e: any) => {
      const text = e.target.result;
      const rows = text.split('\n').map(row => row.split(',').map(cell => cell.trim())); // Split rows by newline and cells by comma
      
      console.log('Raw rows from CSV:', rows); // Log raw rows for debugging
      this.processData(rows);
    };
    reader.readAsText(file);
  }
  
  processData(data: any[]): void {
    if (!Array.isArray(data) || data.length === 0) {
      this.fileError = 'No data found in the file.';
      return;
    }
  
    const headers = data[0]; // First row is assumed to be headers
  
    // Check if headers is an array
    if (!Array.isArray(headers)) {
      this.fileError = 'Invalid header format in CSV.';
      console.error('Headers:', headers); // Log headers for debugging
      return;
    }
  
    this.fileData = data.slice(1) // Skip the header row
      .map(row => {
        const contact: any = {};
        headers.forEach((header: string, index: number) => {
          contact[header] = row[index]; // Map each header to the corresponding value
        });
        
        // Create full_name by combining first_name and last_name
        contact.full_name = `${contact.first_name}`+` `+` ${contact.last_name}`;
  
        return contact;
      })
      .filter(item => item.email && item.first_name && item.last_name && item.phone_number);
  
    console.log('Processed file data:', this.fileData);
    
    if (this.fileData.length === 0) {
      this.fileError = 'No valid contacts found in the file.';
    } else {
      this.fileError = null;
    }
  }
  
  

  onSubmitBulk(): void {
    if (this.fileData.length > 0) {
      this.service.addContactToMailingList(this.fileData, this.data.id).subscribe({
        next: () => {
          this.dialogRef.close({ success: true });
        },
        error: (error) => {
          this.dialogRef.close({ success: false, error: error.error.data });
        }
      });
    }
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}