import { CUSTOM_ELEMENTS_SCHEMA, ChangeDetectorRef, Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatOptionModule, MatOptionSelectionChange } from '@angular/material/core';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatCardModule } from '@angular/material/card'; // Import MatCardModule
import { CommonModule } from '@angular/common';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatInputModule } from '@angular/material/input';
import { MatDialogRef } from '@angular/material/dialog';
import { CampaignService } from '../../example/campaign.service';

@Component({
  selector: 'app-triggers',
  standalone: true,
  templateUrl: './triggers.component.html',
  styleUrls: ['./triggers.component.scss'],
  imports: [
    MatCardModule, // Include MatCardModule in imports
    MatFormFieldModule,
    MatOptionModule,
    MatSelectModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    CommonModule,
    ReactiveFormsModule
  ],
  schemas: [CUSTOM_ELEMENTS_SCHEMA]
})
export class TriggersComponent implements OnInit {
  @Output() formSubmitted = new EventEmitter<any>(); // Adjusted to emit just form data

  constructor(
    private service : CampaignService,
    private fb: FormBuilder,
    private cdRef: ChangeDetectorRef,
    private dialogRef: MatDialogRef<TriggersComponent>
  ) {}
  
  public mailingLists: any = [];

  public triggers: any[] = [
    { name: 'Welcome Subscriber', value: 'welcome_subscriber' },
    { name: 'Specific Date', value: 'specific_date' },
    { name: 'Subscriber Added Date', value: 'subscriber_added_date' }
  ];

  selectedTrigger: string = '';
  selectedMailingList: string = '';
  triggerForm: FormGroup;
  welcomeForm: FormGroup;
  specificDateForm: FormGroup;
  AddedSubDateForm: FormGroup;
  companyID : any;
  ngOnInit() {
    this.companyID = 'afa35ff6-4de5-4806-9a21-e0c2453d2834'; // Use dynamic companyID as needed

    this.service.getMailingLists(this.companyID).subscribe({
      next: (data) => {

        if (data && data.data && data.data.items) {
          this.mailingLists = data.data.items; // Assign mailing list items
          
          
         
        } else {
          console.error('Invalid response structure:', data);
          // Handle unexpected response structure if needed
        }
        console.log(this.mailingLists);
        
      },
      error: (error) => {
        console.error('Error fetching mailing lists:', error);
        // Handle error scenario if needed
      }
    });
    this.triggerForm = this.fb.group({
      selectedTrigger: ['', Validators.required]
    });

    this.welcomeForm = this.fb.group({
      name: ['', Validators.required],
      mailinglist: ['', Validators.required]
    });

    this.specificDateForm = this.fb.group({
      name: ['', Validators.required],
      date: ['', Validators.required],
      mailinglist: ['', Validators.required]
    });

    this.AddedSubDateForm = this.fb.group({
      name: ['', Validators.required],
      daysAfter: ['', [Validators.required, Validators.min(1)]],
      mailinglist: ['', Validators.required]
    });
  }

  onTriggerChange(event: any) {
    this.selectedTrigger = event.value;
  }

  selectTrigger(trigger: { value: string }) {
    this.selectedTrigger = trigger.value;
    this.triggerForm.get('selectedTrigger')?.setValue(trigger.value);
    this.cdRef.detectChanges();
  }

  onSubmitWelcomeForm() {
    if (this.welcomeForm.valid) {
      const formData = this.welcomeForm.value;
      formData.triggerType = 'welcome'; // Add the trigger type
      this.formSubmitted.emit(formData);
      this.dialogRef.close();
    }
  }

  onSubmitSpecificDateForm() {
    
    
    if (this.specificDateForm.valid) {
      const formData = this.specificDateForm.value;
      formData.triggerType = 'specific date'; // Add the trigger type
      this.formSubmitted.emit(formData );
      this.dialogRef.close();
    }
  }

  onSubmitAddedSubDateForm() {
    if (this.AddedSubDateForm.valid) {
      const formData = this.AddedSubDateForm.value;
      formData.triggerType = 'days after'; // Add the trigger type
      this.formSubmitted.emit(formData );

      this.dialogRef.close();
    }
  }
}