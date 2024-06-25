import { CUSTOM_ELEMENTS_SCHEMA, ChangeDetectorRef, Component, EventEmitter, Output } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatOptionModule, MatOptionSelectionChange } from '@angular/material/core';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatCardModule } from '@angular/material/card'; // Import MatCardModule
import { CommonModule } from '@angular/common';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatInputModule } from '@angular/material/input';
import { MatDialogRef } from '@angular/material/dialog';

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
export class TriggersComponent {
  @Output() formSubmitted = new EventEmitter<any>(); // Adjusted to emit just form data

  constructor(
    private fb: FormBuilder,
    private cdRef: ChangeDetectorRef,
    private dialogRef: MatDialogRef<TriggersComponent>
  ) {}
  
  public mailingLists: any = [
    { name: 'Mailing List 1', value: 'mailing_list_1' },
    { name: 'Mailing List 2', value: 'mailing_list_2' },
    { name: 'Mailing List 3', value: 'mailing_list_3' }
  ];

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

  ngOnInit() {
    this.triggerForm = this.fb.group({
      selectedTrigger: ['', Validators.required]
    });

    this.welcomeForm = this.fb.group({
      message: ['', Validators.required],
      sendImmediately: [false],
      
      mailinglist: ['', Validators.required]
    });

    this.specificDateForm = this.fb.group({
      date: ['', Validators.required],
      
      
      mailinglist: ['', Validators.required]
    });

    this.AddedSubDateForm = this.fb.group({
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
      this.formSubmitted.emit(formData );
      this.dialogRef.close();
    }
  }

  onSubmitAddedSubDateForm() {
    if (this.AddedSubDateForm.valid) {
      const formData = this.AddedSubDateForm.value;
      this.formSubmitted.emit(formData );

      this.dialogRef.close();
    }
  }
}