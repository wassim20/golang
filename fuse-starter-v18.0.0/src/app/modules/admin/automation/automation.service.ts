import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { DataTransferService } from './datatransferservice';
import { AuthService } from 'app/core/auth/auth.service';
import { jwtDecode } from 'jwt-decode';

@Injectable({
  providedIn: 'root'
})
export class AutomationService {
 
 
  private getCompanyID(): string | null {
    const token = this.auth.accessToken;
    try {
      const decodedToken: { company_id: string } = jwtDecode(token) as { company_id: string };
      return decodedToken.company_id;
    } catch (error) {
      console.error('Failed to decode JWT token:', error);
      return null;
    }
  }

  private baseUrl = 'http://localhost:8080/api';
 // private jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiYjI3ZWU3N2MtOTA0My00MDAwLWI3ZTAtZjFhOTIwZGEyYzJmIiwiZXhwIjoxNzE5OTY2MDE4LCJpYXQiOjE3MTk3OTMyMTgsInJvbGVzIjpbeyJpZCI6ImMyNmQ0MGE2LTc0MzMtNDg1MC1hYmUyLWIwMzIxYzNkNDA3ZiIsIm5hbWUiOiJNYW5hZ2VyIiwiY29tcGFueV9pZCI6ImIyN2VlNzdjLTkwNDMtNDAwMC1iN2UwLWYxYTkyMGRhMmMyZiJ9XSwidXNlcl9pZCI6IjgzODJlYzZkLTI1OWUtNDQxMy04NjIyLWFiYTFlZWM3MzdlOSJ9.VJn51F28ekCyBpKnwMH-9olqj-FtTp715aBg5LGi2O4'

  constructor(private datatransfer:DataTransferService,private http: HttpClient,private auth :AuthService) { }

  private getHeaders(): HttpHeaders {
    let headers = new HttpHeaders()
      .set('Authorization', `Bearer ${this.auth.accessToken}`)
      .set('Content-Type', 'application/json');
    return headers;
  }


  createWorkflow(workflowData: any) {
    
    const companyID = this.getCompanyID() || '';
    return this.http.post<any>(`${this.baseUrl}/${companyID}/workflow`, workflowData, {
      headers: this.getHeaders()
    }).subscribe(
      response => {
        const workflowData = this.datatransfer.getWorkflowData();
        
        this.datatransfer.setWorkflowData({ ...workflowData, id: response.data });
        console.log(this.datatransfer.getWorkflowData());
        console.log("Response with ID:", response);
      },
      error => {
        console.error("Error:", error);
      }
    );
}

createAction(actionData: any,workflowId: any) {
  const companyID = this.getCompanyID() || '';
  return this.http.post<any>(`${this.baseUrl}/${companyID}/workflow/${workflowId}/action`, actionData, {
    headers: this.getHeaders()
  });
}

getAutomations(page: number = 1, limit: number = 10) {
  const companyID = this.getCompanyID() || '';
  console.log("Company ID:", companyID);
  console.log("Headers:", this.getHeaders());
  return this.http.get<any>(`${this.baseUrl}/${companyID}/workflow`, {
    headers: this.getHeaders(),
    params: {
      page: page.toString(),
      limit: limit.toString()
    }
  });
}
getActions(workflowId: string) {
  const companyID = this.getCompanyID() || '';
  return this.http.get<any>(`${this.baseUrl}/${companyID}/workflow/${workflowId}/action`, {
    headers: this.getHeaders()
  });
}

updateAutomation(automationId: string, automationData: any) {
  const companyID = this.getCompanyID() || '';
  return this.http.put<any>(`${this.baseUrl}/${companyID}/workflow/${automationId}`, automationData, {
    headers: this.getHeaders()
  });
}

deleteAutomation(automationId: string) {
  const companyID = this.getCompanyID() || '';
  return this.http.delete<any>(`${this.baseUrl}/${companyID}/workflow/${automationId}`, {
    headers: this.getHeaders()
  });
}

updateAction(worflowId:string, actionData: any) {
  const companyID = this.getCompanyID() || '';
  return this.http.put<any>(`${this.baseUrl}/${companyID}/workflow/${worflowId}/action/${actionData.id}`, {type:actionData.type,data:actionData.data}, {
    headers: this.getHeaders()
  });
}
startWorkflow(workflowId: string) {
  const companyID = this.getCompanyID() || '';
  return this.http.post<any>(`${this.baseUrl}/${companyID}/workflow/${workflowId}/start`, {}, {
    headers: this.getHeaders()
  });
}

getCampaigns() {
  const companyID = this.getCompanyID() || '';
  return this.http.get<any>(`${this.baseUrl}/${companyID}/campaigns/all`, {
    headers: this.getHeaders(),
  });
}
}
