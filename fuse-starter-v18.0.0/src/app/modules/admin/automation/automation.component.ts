
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit, ViewEncapsulation } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { TriggersComponent } from './triggers/triggers.component';
import { Router, RouterModule } from '@angular/router';
import { AsyncPipe, CommonModule, CurrencyPipe, NgClass, NgFor, NgIf, NgTemplateOutlet } from '@angular/common';
import { DataTransferService } from './datatransferservice';
import { AutomationService } from './automation.service';
import { MatProgressBar, MatProgressBarModule } from '@angular/material/progress-bar';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIcon, MatIconModule } from '@angular/material/icon';
import { FormBuilder, FormControlDirective, FormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { NgxMatDatetimePickerModule, NgxMatNativeDateModule, NgxMatTimepickerModule } from '@angular-material-components/datetime-picker';
import { DateTimePickerModule } from '@syncfusion/ej2-angular-calendars';
import { MatCardModule } from '@angular/material/card';
import { MatSortModule } from '@angular/material/sort';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatSelectModule } from '@angular/material/select';
import { MatOptionModule, MatRippleModule } from '@angular/material/core';
import { MatCheckboxModule } from '@angular/material/checkbox';



@Component({
  selector: 'automation',
  styleUrls: ['./automation.component.scss'],
  templateUrl: './automation.component.html',
  standalone   : true,
  imports :[CommonModule,NgIf, MatProgressBarModule, MatFormFieldModule, MatIconModule, MatInputModule, FormsModule, ReactiveFormsModule, MatButtonModule, MatSortModule, NgFor, NgTemplateOutlet, MatPaginatorModule, NgClass, MatSlideToggleModule, MatSelectModule, MatOptionModule, MatCheckboxModule, MatRippleModule, AsyncPipe, CurrencyPipe,
    NgxMatDatetimePickerModule,
    RouterModule,
      NgxMatTimepickerModule,
      NgxMatNativeDateModule,
      DateTimePickerModule,
      MatCardModule,
  ],
  encapsulation: ViewEncapsulation.None,
  
})
export class AutomationComponent implements OnInit {

  selectedAction: string = '';
  showWorkflow: boolean = false;
  searchControl: any;
  automations : any[] = [];
  selectedAutomation: any;
  selectedAutomationForm: FormGroup;
  flashMessage: string = '';


  constructor(
    private service: AutomationService,
    private router: Router,
    public dialog: MatDialog,
    private cdRef: ChangeDetectorRef,
    private dataTransferService: DataTransferService,
    private fb: FormBuilder,
    private _changeDetectorRef: ChangeDetectorRef
  ) {}

 

  createWorkflow() {
    const dialogRef = this.dialog.open(TriggersComponent);
    

    dialogRef.componentInstance.formSubmitted.subscribe((formData: any) => {
      const finalData = {
        name: formData.name,
        mailinglist_id: formData.mailinglist,
        trigger: formData.triggerType,
        trigger_data: {
            ...(formData.daysAfter !== undefined && { daysAfter: formData.daysAfter }),
            ...(formData.date !== undefined && { date: formData.date }),
        },
    };
      
      
      this.showWorkflow = true; 
      if(formData){
        this.dataTransferService.setWorkflowData(formData);
        this.service.createWorkflow(finalData)
      
        this.router.navigate(['automation/workflow'], { state: {workflowData:JSON.stringify(formData)  } }).then(() => {
          console.log("Workflow Component should now be loaded.");
      });;
      }  
      dialogRef.close();
      
    });
  }

 

  ngOnInit(): void {
    this.fetchAutomations();
    this.initForm();
  }

  // Initialize the form
  initForm(): void {
    this.selectedAutomationForm = this.fb.group({
      name: [''],
      status: [''],
     
    });
  }

  // Fetch automations
  fetchAutomations(): void {
    this.service.getAutomations().subscribe((data) => {
      this.automations = data.data.items;
    });
  }

  // Toggle details visibility
  toggleDetails(automationId: string): void {
    if (this.selectedAutomation?.id === automationId) {
      this.selectedAutomation = null;
    } else {
      this.selectedAutomation = this.automations.find(
        (automation) => automation.id === automationId
      );
      this.patchForm(this.selectedAutomation);
    }
  }

  // Patch the form with selected automation details
  patchForm(automation: any): void {
    this.selectedAutomationForm.patchValue({
      name: automation.name,
      status: automation.status,
     
    });
  }

  // Update selected automation
  updateSelectedAutomation(): void {
    if (this.selectedAutomationForm.valid) {
        const updatedData = this.selectedAutomationForm.value;

        this.service.updateAutomation(this.selectedAutomation.id, updatedData)
            .subscribe({
                next: () => {
                    // Update the selectedAutomation object with the updated form data
                    Object.assign(this.selectedAutomation, updatedData);

                    this.flashMessage = 'success';
                    this.showFlashMessage('success');
                    this._changeDetectorRef.detectChanges();
                },
                error: () => {
                    this.flashMessage = 'error';
                    this.showFlashMessage('error');
                },
            });
    }
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

  // Delete selected automation
  deleteSelectedAutomation(): void {
    this.service.deleteAutomation(this.selectedAutomation.id).subscribe({
      next: () => {
        this.selectedAutomation = null;
        this.fetchAutomations();
      },
      error: () => {
        this.flashMessage = 'error';
      },
    });
  }



}