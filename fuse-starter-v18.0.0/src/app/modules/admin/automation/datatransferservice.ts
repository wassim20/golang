import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class DataTransferService {
  private workflowData: any;

  setWorkflowData(data: any) {
    console.log("here");
    
    this.workflowData = data;
  }

  getWorkflowData() {
    return this.workflowData;
  }

  clearWorkflowData() {
    this.workflowData = null;
  }
}
