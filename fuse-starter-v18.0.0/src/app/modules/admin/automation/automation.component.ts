
import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit, ViewEncapsulation } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { TriggersComponent } from './triggers/triggers.component';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { DataTransferService } from './datatransferservice';


@Component({
  selector: 'automation',
  templateUrl: './automation.component.html',
  standalone   : true,
  imports :[MatButtonModule,MatDialogModule,RouterModule,CommonModule,],
  encapsulation: ViewEncapsulation.None,
  
})
export class AutomationComponent implements OnInit {

  selectedAction: string = '';
  showWorkflow: boolean = false;

  constructor(private router: Router,public dialog: MatDialog,private cdRef: ChangeDetectorRef, private dataTransferService: DataTransferService) {}

  ngOnInit(): void { }

  createWorkflow() {
    const dialogRef = this.dialog.open(TriggersComponent);
    

    dialogRef.componentInstance.formSubmitted.subscribe((formData: any) => {
      // this.handleFormSubmission(formData);
      this.showWorkflow = true; // Close the dialog after handling form submission
      //go to child compoenent workflow
      if(formData){
        this.dataTransferService.setWorkflowData(formData);
        this.router.navigate(['automation/workflow'], { state: {workflowData:JSON.stringify(formData)  } });
      }
     

      // this.router.navigate(['automation', 'workflow'], { state: { data: formData } });
      console.log(formData);
      console.log(this.showWorkflow);
      
      
      dialogRef.close();
      
    });
  }

  handleFormSubmission(formData: any) {

    console.log(this.showWorkflow);
    
    console.log(formData);
    //get the data to children route workflow
    
    
    
  }

}