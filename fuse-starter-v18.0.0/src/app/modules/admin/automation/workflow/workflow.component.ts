import { Component, ElementRef, Inject, OnInit, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { DataTransferService } from '../datatransferservice';
import { AutomationService } from '../automation.service';
import { MatOptionModule } from '@angular/material/core';
import { duration } from 'html2canvas/dist/types/css/property-descriptors/duration';
@Component({
  selector: 'workflow',
  templateUrl: './workflow.component.html',
  standalone: true,
  styleUrls: ['./workflow.component.scss'],
  imports: [
    MatCardModule,
    CommonModule,
    MatDialogModule,
    MatInputModule,
    MatListModule,
    MatIconModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    FormsModule,
  ],
})
export class WorkflowComponent implements OnInit {
  @ViewChild('workflowContainer') workflowContainer!: ElementRef;
  @ViewChild('workflowArea') workflowArea!: ElementRef;


  workflowdata: any;
  actions: any[] = [];
  zoomLevel: number = 1;
  zoomTransform: string = 'scale(1)';
  previousAction: { id: any; title: any; type: any; parent_id: any; };
  constructor(private service: AutomationService,private datatransfer:DataTransferService,public dialog: MatDialog) {}

  ngOnInit(): void {
    // Get workflow data from the dataTransferService
    this.workflowdata = this.datatransfer.getWorkflowData();

    // Load the actions if there is existing workflow data
    if (this.workflowdata) {
        this.loadActions(this.workflowdata.id);
    }
    
    
}

loadActions(workflowId: string): void {
  this.service.getActions(workflowId).subscribe((response: any) => {
    const actions = response.data.items;

    // Initialize a map to store actions by their ID
    const actionMap = new Map<string, any>();

    // Step 1: Populate the map and initialize branches
    actions.forEach((action: any) => {
      action.yesBranch = [];
      action.noBranch = [];
      action.subActions = []; // Initialize subActions for non-condition actions
      actionMap.set(action.id, action);
    });

    // Step 2: Build the action tree
    const rootActions: any[] = [];

    actions.forEach((action: any) => {
      if (action.parent_id && action.parent_id !== "00000000-0000-0000-0000-000000000000") {
        const parentAction = actionMap.get(action.parent_id);
        if (parentAction) {
          const actionData = action.data ? JSON.parse(action.data) : {};

          if (parentAction.type === 'condition') {
            // Place in the correct branch based on the action's data
            if (actionData.branch === 'yes') {
              parentAction.yesBranch.push(action);
            } else if (actionData.branch === 'no') {
              parentAction.noBranch.push(action);
            }
          } else {
            // For non-condition parent actions, add to subActions
            parentAction.subActions.push(action);
          }
        }
      } else {
        // Root actions without parent_id
        rootActions.push(action);
      }
    });

    // Flatten the structured root actions into a single array for rendering
    this.actions = rootActions.flatMap((action: any) => this.flattenActionTree(action));

    // Optional: Log the structured actions to verify
    console.log('Structured Actions:', this.actions);
  });
}

private flattenActionTree(action: any, isBranch: boolean = false): any[] {
  // If this action is part of a branch (yesBranch or noBranch), do not add it to the flattened array
  if (isBranch) {
    return [];
  }

  const flattened = [action];

  // Flatten the subActions, which are direct children of non-condition parents
  if (action.subActions && action.subActions.length > 0) {
    action.subActions.forEach((subAction: any) => {
      flattened.push(...this.flattenActionTree(subAction));
    });
  }

  // Flatten the branches, but mark them as part of a branch
  if (action.yesBranch && action.yesBranch.length > 0) {
    action.yesBranch.forEach((subAction: any) => {
      flattened.push(...this.flattenActionTree(subAction, true));
    });
  }

  if (action.noBranch && action.noBranch.length > 0) {
    action.noBranch.forEach((subAction: any) => {
      flattened.push(...this.flattenActionTree(subAction, true));
    });
  }

  return flattened;
}



  addAction(parentAction?: any, branch?: 'yes' | 'no'): void {
    const dialogRef = this.dialog.open(ActionDialogComponent, {
        width: '300px',
        data: { parentAction }
    });

    dialogRef.afterClosed().subscribe(result => {
        if (result) {
            // Set parent_id based on the correct logic:
            if (parentAction && branch) {
                // Get the last action in the selected branch
                const lastActionInBranch = this.getLastActionInBranch(parentAction, branch);
                result.parent_id = lastActionInBranch ? lastActionInBranch.id : parentAction.id;
            } else {
                // No branch, so set parent_id to the parentAction's id or null if it's the first action
                result.parent_id = parentAction ? parentAction.id : null;
            }

            // Save the new action to the backend and update the branch accordingly
            this.saveAction(result, branch, parentAction);
        }
    });
}
// Define the isNested function
isNested(action: any): boolean {
  return action?.parent_id !== null && action?.parent_id !== undefined;
}
saveAction(actionData: any, branch?: 'yes' | 'no', parentAction?: any): void {
  const workflowData = this.datatransfer.getWorkflowData();
  
  this.service.createAction(actionData, workflowData.id).subscribe(
      (savedAction: any) => {
          // Create the new action object with the saved ID
          const newAction = {
              id: savedAction.data.ID,
              title: actionData.title,
              type: actionData.type,
              parent_id: savedAction.data.parent_id,
              yesBranch: [],
              noBranch: []
          };

          if (parentAction) {
              if (branch === 'yes') {
                  if (!parentAction.yesBranch) {
                      parentAction.yesBranch = [];
                  }
                  parentAction.yesBranch.push(newAction);
              } else if (branch === 'no') {
                  if (!parentAction.noBranch) {
                      parentAction.noBranch = [];
                  }
                  parentAction.noBranch.push(newAction);
              } else {
                  this.actions.push(newAction);
              }
          } else {
              this.actions.push(newAction);
          }
      },
      (error) => {
          console.error('Error saving action:', error);
      }
  );
  
}

// Helper function to get the last action in a branch (Yes or No)
getLastActionInBranch(parentAction: any, branch: 'yes' | 'no'): any {
    if (branch === 'yes' && parentAction.yesBranch.length > 0) {
        return parentAction.yesBranch[parentAction.yesBranch.length - 1];
    } else if (branch === 'no' && parentAction.noBranch.length > 0) {
        return parentAction.noBranch[parentAction.noBranch.length - 1];
    }
    return null;
}

isActionComplete(action: any): boolean {
  const data = action.data ? JSON.parse(action.data) : {}; // Parse the JSON string
  if (action.type === 'email') {
    return data.subject && data.track_open !== undefined && data.track_click !== undefined && data.HTML && data.from && data.reply_to;
  } else if (action.type === 'wait') {
    return data.duration;
  } else if (action.type === 'condition') {
    return data.criteria && data.duration;
  }
  return true;
}
getIncompleteMessage(action: any): string {
  if (action.type === 'email') {
      return `double click to edit`;
  } else if (action.type === 'wait') {
      return `double click to edit`;
  } else if (action.type === 'condition') {
      return `double click to edit`;
  }
  return 'Incomplete action';
}







editAction(event: MouseEvent, action: any, branch?: 'yes' | 'no'): void {
  event.stopPropagation();

  // Parse the action data if it exists
  let parsedData = {};
  if (action.data) {
    try {
      parsedData = JSON.parse(action.data);
    } catch (error) {
      console.error("Error parsing action data:", error);
    }
  }

  // Open the edit dialog and pass the parsed data along with branch information
  const dialogRef = this.dialog.open(EditActionDialogComponent, {
    width: '600px',
    data: { ...action, ...parsedData, branch } // Merge action with parsed data and branch info
  });

  const workflowData = this.datatransfer.getWorkflowData();

  dialogRef.afterClosed().subscribe(result => {
    if (result) {
      action.data = JSON.stringify(result); // Save the updated data back to the action

      const editedAction = {
        id: action.id,
        type: action.type,
        data: action.data
      };

      // Call backend to update the action
      this.service.updateAction(workflowData.id, editedAction).subscribe(
        (updatedAction) => {
          console.log('Action updated successfully:', updatedAction);
        },
        (error) => {
          console.error('Error updating action:', error);
        }
      );
    }
  });
}


areAllActionsComplete(): boolean {
  return this.actions.every(action => this.isActionComplete(action));
}
startWorkflow() {
  const workflowData = this.datatransfer.getWorkflowData();
  this.service.startWorkflow(workflowData.id).subscribe(
      (response) => {
          console.log('Workflow started successfully:', response);
      },
      (error) => {
          console.error('Error starting workflow:', error);
      }
  );

}

  removeAction(index: number, parentAction?: any, branch?: 'yes' | 'no'): void {
    if (parentAction) {
      if (branch) {
        parentAction[branch + 'Branch'].splice(index, 1);
      } else {
        parentAction.children.splice(index, 1);
      }
    } else {
      this.actions.splice(index, 1);
    }
  }

  // Method to check if there is any condition at a given level
  hasCondition(actions = this.actions): boolean {
    return actions.some(action => action.type === 'condition');
  }

  zoomIn(): void {
    this.zoomLevel += 0.1;
    this.updateZoom();
  }

  zoomOut(): void {
    this.zoomLevel -= 0.1;
    this.updateZoom();
  }

  resetZoom(): void {
    this.zoomLevel = 1;
    this.updateZoom();
  }

  updateZoom(): void {
    this.zoomTransform = `scale(${this.zoomLevel})`;
  }

  toggleFullScreen(): void {
    const elem = this.workflowContainer.nativeElement;

    if (!document.fullscreenElement) {
      elem.requestFullscreen().catch((err: any) => {
        console.log(`Error attempting to enable full-screen mode: ${err.message}`);
      });
    } else {
      document.exitFullscreen();
    }
  }
}

@Component({
  selector: 'app-edit-action-dialog',
  standalone: true,
  styleUrls: ['./workflow.component.scss'],
  imports: [
    MatFormFieldModule,
    MatInputModule,
    
    ReactiveFormsModule,
    MatSelectModule,
    MatOptionModule,
    MatIconModule,
    MatButtonModule,
    CommonModule,
    FormsModule
  ],
  template: `
    <div mat-dialog-content class="edit-container" *ngIf="data.type === 'email'">
      <h2  mat-dialog-title>Edit Action</h2>
      <form  class="stepper-form">
      <ng-container *ngIf="data.type === 'email'" [formGroup]="form" class="form-container">
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Subject</mat-label>
          <input matInput formControlName="subject" required>
          <mat-error *ngIf="form.controls['subject'].touched && form.controls['subject'].invalid">
            Please enter a valid subject.
          </mat-error>
        </mat-form-field>
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Track Open</mat-label>
          <mat-select formControlName="track_open" required>
            <mat-option [value]="true">Yes</mat-option>
            <mat-option [value]="false">No</mat-option>
          </mat-select>
          <mat-error *ngIf="form.controls['track_open'].touched && form.controls['track_open'].invalid">
            Please select a valid option.
          </mat-error>
        </mat-form-field>
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Track Click</mat-label>
          <mat-select formControlName="track_click" required>
            <mat-option [value]="true">Yes</mat-option>
            <mat-option [value]="false">No</mat-option>
          </mat-select>
          <mat-error *ngIf="form.controls['track_click'].touched && form.controls['track_click'].invalid">
            Please select a valid option.
          </mat-error>
        </mat-form-field>
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>HTML</mat-label>
          <textarea matInput formControlName="HTML" required></textarea>
          <mat-error *ngIf="form.controls['HTML'].touched && form.controls['HTML'].invalid">
            Please enter HTML content.
          </mat-error>
        </mat-form-field>
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>From</mat-label>
          <input matInput formControlName="from" required>
          <mat-error *ngIf="form.controls['from'].touched && form.controls['from'].invalid">
            Please enter a valid sender email.
          </mat-error>
        </mat-form-field>
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Reply-To</mat-label>
          <input matInput formControlName="reply_to" required>
          <mat-error *ngIf="form.controls['reply_to'].touched && form.controls['reply_to'].invalid">
            Please enter a valid reply-to email.
          </mat-error>
        </mat-form-field>
      </ng-container>
      </form>
    </div>

    <div mat-dialog-content class="edit-container" *ngIf="data.type === 'wait'">
    <h2  mat-dialog-title>Edit Action</h2>
      <!-- Wait Action Fields -->
      <form  class="stepper-form">
      <ng-container *ngIf="data.type === 'wait'" [formGroup]="form" class="form-container">
        <mat-form-field appearance="outline" class="triyple-full-width">
          <mat-label>Duration</mat-label>
          <div matPrefix>
            <input matInput formControlName="durationValue" placeholder="Enter duration" type="number" min="0" required>
          </div>
          <mat-select formControlName="durationUnit" required>
            <mat-option value="s">Seconds</mat-option>
            <mat-option value="m">Minutes</mat-option>
            <mat-option value="h">Hours</mat-option>
            <mat-option value="d">Days</mat-option>
          </mat-select>
          <mat-error *ngIf="form.controls['durationValue'].touched && form.controls['durationValue'].invalid">
            Please enter a valid duration.
          </mat-error>
        </mat-form-field>
      </ng-container>
      </form>
    </div>

      <!-- Condition Action Fields -->
    <div mat-dialog-content class="edit-container" *ngIf="data.type === 'condition'">
    <h2  mat-dialog-title>Edit Action</h2>
      <form  class="stepper-form">
      <ng-container *ngIf="data.type === 'condition'" [formGroup]="form" class="form-container">
        <mat-form-field appearance="outline" class="double-full-width">
          <mat-label>Criteria</mat-label>
          <mat-select formControlName="criteria" required>
            <mat-option [value]="'read'">Read</mat-option>
            <mat-option [value]="'click'">Click</mat-option>
          </mat-select>
          <mat-error *ngIf="form.controls['criteria'].touched && form.controls['criteria'].invalid">
            Please select a valid criteria.
          </mat-error>
        </mat-form-field>
        <mat-form-field appearance="outline" class="full-width">
          <mat-label>Duration</mat-label>
          <div matPrefix>
            <input matInput formControlName="durationValue" placeholder="Enter duration" type="number" min="0" required>
          </div>
          <mat-select formControlName="durationUnit" required>
            <mat-option value="s">Seconds</mat-option>
            <mat-option value="m">Minutes</mat-option>
            <mat-option value="h">Hours</mat-option>
            <mat-option value="d">Days</mat-option>
          </mat-select>
          <mat-error *ngIf="form.controls['durationValue'].touched && form.controls['durationValue'].invalid">
            Please enter a valid duration.
          </mat-error>
        </mat-form-field>
      </ng-container>
      </form>
    </div>
    <div mat-dialog-actions class="stepper-buttons">
      <button mat-button (click)="onCancel()">Cancel</button>
      <button mat-flat-button color="primary" (click)="onSave()">Save</button>
    </div>
  `,
})
export class EditActionDialogComponent implements OnInit {
  form: FormGroup;

  constructor(
    public dialogRef: MatDialogRef<EditActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private fb: FormBuilder
  ) {}

  ngOnInit(): void {
    console.log("Data passed to dialog:", this.data);

    if (this.data.type === 'email') {
      this.form = this.fb.group({
        subject: [this.data.subject || '', Validators.required],
        track_open: [this.data.track_open ?? null, Validators.required],
        track_click: [this.data.track_click ?? null, Validators.required],
        HTML: [this.data.HTML || '', Validators.required],
        from: [this.data.from || '', Validators.required],
        reply_to: [this.data.reply_to || '', Validators.required],
      });
    } else if (this.data.type === 'wait') {
      const duration = this.parseDuration(this.data.duration || '0s');
      this.form = this.fb.group({
        durationValue: [duration.value, [Validators.required, Validators.min(0)]],
        durationUnit: [duration.unit, Validators.required],
      });
    } else if (this.data.type === 'condition') {
      this.form = this.fb.group({
        criteria: [this.data.criteria || '', Validators.required],
        durationValue: [this.parseDuration(this.data.duration || '0s').value, [Validators.required, Validators.min(0)]],
        durationUnit: [this.parseDuration(this.data.duration || '0s').unit, Validators.required],
      });
    }
  }

  // Helper function to parse the duration from the format '6s', '2m', etc.
  parseDuration(duration: string): { value: number, unit: string } {
    const match = duration.match(/(\d+)([smhd])/);
    return {
      value: match ? parseInt(match[1], 10) : 0,
      unit: match ? match[2] : 's',
    };
  }

  onSave(): void {
    if (this.form.valid) {
      const actionData = this.form.value;
      let branchData = {};

      if (this.data.branch) {
        branchData = { branch: this.data.branch };
      }

      if (this.data.type === 'email') {
        const emailData = {
          type: 'email',
          subject: actionData.subject,
          track_open: actionData.track_open,
          track_click: actionData.track_click,
          HTML: actionData.HTML,
          from: actionData.from,
          reply_to: actionData.reply_to,
          ...branchData,
        };
        this.dialogRef.close(emailData);
      } else if (this.data.type === 'wait') {
        const waitData = {
          type: 'wait',
          duration: actionData.durationValue + actionData.durationUnit,
          ...branchData,
        };
        this.dialogRef.close(waitData);
      } else if (this.data.type === 'condition') {
        const conditionData = {
          criteria: actionData.criteria,
          duration: actionData.durationValue + actionData.durationUnit,
          ...branchData,
        };
        this.dialogRef.close(conditionData);
      }
    }
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}


@Component({
  standalone: true,
  imports: [
    MatListModule,
    MatButtonModule,
    CommonModule
  ],
  template: `
    <h2 mat-dialog-title>Select Action</h2>
    <div mat-dialog-content>
      <mat-list>
        <mat-list-item (click)="selectAction('Email')">Email</mat-list-item>
        <mat-list-item (click)="selectAction('Wait')">Wait</mat-list-item>
        <mat-list-item (click)="selectAction('Condition')">Condition</mat-list-item>
      </mat-list>
    </div>
    <div mat-dialog-actions>
      <button mat-button (click)="onCancel()">Cancel</button>
    </div>
  `,
})
export class ActionDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<ActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}

  selectAction(actionType: string): void {
   
    
    const action = {
      id:'',
      type: actionType.toLowerCase(),
      title: actionType,
      parent_id: this.data.parentAction?.id ?? null,
      yesBranch: [],
      noBranch: []
    };
    console.log("action",action);
    this.dialogRef.close(action);
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}
