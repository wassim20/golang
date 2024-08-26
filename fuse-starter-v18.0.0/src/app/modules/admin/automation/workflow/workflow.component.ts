import { Component, ElementRef, Inject, ViewChild } from '@angular/core';
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
export class WorkflowComponent {
  @ViewChild('workflowContainer') workflowContainer!: ElementRef;
  @ViewChild('workflowArea') workflowArea!: ElementRef;

  actions: any[] = [];
  zoomLevel: number = 1;
  zoomTransform: string = 'scale(1)';
  constructor(public dialog: MatDialog) {}

  addAction(parentAction?: any, branch?: 'yes' | 'no'): void {
    const dialogRef = this.dialog.open(ActionDialogComponent, {
      width: '300px',
      data: { parentAction }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        if (parentAction) {
          if (branch) {
            parentAction[branch + 'Branch'] = parentAction[branch + 'Branch'] || [];
            parentAction[branch + 'Branch'].push(result);
          } else {
            parentAction.children = parentAction.children || [];
            parentAction.children.push(result);
          }
        } else {
          this.actions.push(result);
        }
      }
    });
  }

  editAction(action: any): void {
    const dialogRef = this.dialog.open(EditActionDialogComponent, {
      width: '300px',
      data: action
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        Object.assign(action, result);
      }
    });
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
  template: `
    <h2 mat-dialog-title>Edit Action</h2>
    <div mat-dialog-content>
      <mat-form-field *ngIf="data.type === 'email'">
        <mat-label>Subject</mat-label>
        <input matInput [(ngModel)]="data.subject" />
      </mat-form-field>
      <mat-form-field *ngIf="data.type === 'email'">
        <mat-label>Body</mat-label>
        <textarea matInput [(ngModel)]="data.body"></textarea>
      </mat-form-field>
      <mat-form-field *ngIf="data.type === 'wait'">
        <mat-label>Duration</mat-label>
        <input matInput [(ngModel)]="data.duration" />
        <mat-select [(ngModel)]="data.unit">
          <mat-option value="minutes">Minutes</mat-option>
          <mat-option value="hours">Hours</mat-option>
          <mat-option value="days">Days</mat-option>
        </mat-select>
      </mat-form-field>
      <mat-form-field *ngIf="data.type === 'condition'">
        <mat-label>Condition Type</mat-label>
        <mat-select [(ngModel)]="data.conditionType">
          <mat-option value="Open Previous Email">Open Previous Email</mat-option>
        </mat-select>
      </mat-form-field>
    </div>
    <div mat-dialog-actions>
      <button mat-button (click)="onCancel()">Cancel</button>
      <button mat-button (click)="onSave()">Save</button>
    </div>
  `,
})
export class EditActionDialogComponent {
  constructor(
    public dialogRef: MatDialogRef<EditActionDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any
  ) {}

  onSave(): void {
    this.dialogRef.close(this.data);
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
      type: actionType.toLowerCase(),
      title: actionType,
      yesBranch: [],
      noBranch: []
    };

    this.dialogRef.close(action);
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}
