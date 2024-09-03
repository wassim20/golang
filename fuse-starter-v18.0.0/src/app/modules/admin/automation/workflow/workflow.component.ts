import { Component, ElementRef, Inject, OnInit, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatListModule } from '@angular/material/list';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { DataTransferService } from '../datatransferservice';
import { AutomationService } from '../automation.service';
import { MatOptionModule } from '@angular/material/core';
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
    console.log('loaded worflow');
    
    
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
  if (action.type === 'email') {
    return action.subject && action.track_open !== undefined && action.track_click !== undefined && action.HTML && action.from && action.reply_to;
  } else if (action.type === 'wait') {
    return action.duration;
  } else if (action.type === 'condition') {
    return action.criteria && action.duration;
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







editAction(action: any): void {
  const dialogRef = this.dialog.open(EditActionDialogComponent, {
      width: '300px',
      data: action
  });
  const workflowData = this.datatransfer.getWorkflowData();
  console.log("workflowData",workflowData);
  
  dialogRef.afterClosed().subscribe(result => {
    
      if (result) {
        console.log("result",result);
        

          action.data = result;
          console.log("action",action);
          

          // Call backend to update the action
          this.service.updateAction(workflowData.id,action ).subscribe(
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
  standalone: true,
  imports: [ MatButtonModule, CommonModule,MatFormFieldModule,MatInputModule,MatSelectModule,FormsModule, MatOptionModule,MatSelectModule,MatInputModule],
  template: `
   <div mat-dialog-content>
  <!-- Email Action Fields -->
  <ng-container *ngIf="data.type === 'email'">
    <mat-form-field>
      <mat-label>Subject</mat-label>
      <input matInput [(ngModel)]="data.subject">
    </mat-form-field>
    <mat-form-field>
      <mat-label>Track Open</mat-label>
      <mat-select [(ngModel)]="data.track_open">
        <mat-option [value]="true">Yes</mat-option>
        <mat-option [value]="false">No</mat-option>
      </mat-select>
    </mat-form-field>
    <mat-form-field>
      <mat-label>Track Click</mat-label>
      <mat-select [(ngModel)]="data.track_click">
        <mat-option [value]="true">Yes</mat-option>
        <mat-option [value]="false">No</mat-option>
      </mat-select>
    </mat-form-field>
    <mat-form-field>
      <mat-label>HTML</mat-label>
      <textarea matInput [(ngModel)]="data.HTML"></textarea>
    </mat-form-field>
    <mat-form-field>
      <mat-label>From</mat-label>
      <input matInput [(ngModel)]="data.from">
    </mat-form-field>
    <mat-form-field>
      <mat-label>Reply-To</mat-label>
      <input matInput [(ngModel)]="data.reply_to">
    </mat-form-field>
  </ng-container>

  <!-- Wait Action Fields -->
  <ng-container *ngIf="data.type === 'wait'">
    <mat-form-field>
      <mat-label>Duration</mat-label>
      <input matInput [(ngModel)]="data.duration">
    </mat-form-field>
  </ng-container>

  <!-- Condition Action Fields -->
  <ng-container *ngIf="data.type === 'condition'">
    <mat-form-field>
      <mat-label>Criteria</mat-label>
      <mat-select [(ngModel)]="data.criteria">
        <mat-option [value]="'read'">read</mat-option>
        <mat-option [value]="'click'">click</mat-option>
      </mat-select>
    </mat-form-field>
    <!-- <mat-form-field>
      <mat-label>Campaign ID</mat-label>
      <input matInput [(ngModel)]="data.campaignID">
    </mat-form-field> -->
    <mat-form-field>
      <mat-label>Duration</mat-label>
      <input matInput [(ngModel)]="data.duration">
    </mat-form-field>

  </ng-container>
</div>
<div mat-dialog-actions>
  <button mat-button (click)="onCancel()">Cancel</button>
  <button mat-button color="primary" (click)="onSave()">Save</button>
</div>

  `,
})
export class EditActionDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<EditActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}

  onSave(): void {
    let actionData: any = {};

    // Check the action type and include the relevant fields
    if (this.data.type === 'email') {
      actionData = {
        subject: this.data.subject,
        track_open: this.data.track_open,
        track_click: this.data.track_click,
        HTML: this.data.HTML,
        from: this.data.from,
        reply_to: this.data.reply_to,
      };
    } else if (this.data.type === 'wait') {
      actionData = {
        duration: this.data.duration,
      };
    } else if (this.data.type === 'condition') {
      actionData = {
        criteria: this.data.criteria,
        duration: this.data.duration,
      };
    }

    this.dialogRef.close(actionData);
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
