import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataTransferService } from '../datatransferservice';
import { FuseCardComponent } from "../../../../../@fuse/components/card/card.component";
import {MatCardModule} from '@angular/material/card';

@Component({
    selector: 'workflow',
    templateUrl: './workflow.component.html',
    standalone: true,
    styleUrls: ['./workflow.component.scss'],
    imports: [FuseCardComponent,MatCardModule]
})
export class WorkflowComponent implements OnInit {
  actions: any[] = [];

  addAction() {
    this.actions.push({
      type: '',
      message: 'test message',
    });
  }

  removeAction(index: number) {
    this.actions.splice(index, 1);
  }

  formData: any;

  constructor(private router: Router,private route: ActivatedRoute, private dataTransferService: DataTransferService) { }

  // workflow.component.ts
  ngOnInit() {
    this.formData = this.dataTransferService.getWorkflowData();
    console.log('Form Data in Workflow Component:', this.formData);

    if (!this.formData) {
      console.error('No form data found in state.');
    } else {
      this.dataTransferService.clearWorkflowData(); // Optionally clear data after use
    }
  }
}