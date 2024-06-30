import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { throwError } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CampaignService {
 private companyID ="b27ee77c-9043-4000-b7e0-f1a920da2c2f" ;
 private baseUrl = 'http://localhost:8080/api';
 private jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiYjI3ZWU3N2MtOTA0My00MDAwLWI3ZTAtZjFhOTIwZGEyYzJmIiwiZXhwIjoxNzE5NzA0NzUzLCJpYXQiOjE3MTk1MzE5NTMsInJvbGVzIjpbeyJpZCI6ImMyNmQ0MGE2LTc0MzMtNDg1MC1hYmUyLWIwMzIxYzNkNDA3ZiIsIm5hbWUiOiJNYW5hZ2VyIiwiY29tcGFueV9pZCI6ImIyN2VlNzdjLTkwNDMtNDAwMC1iN2UwLWYxYTkyMGRhMmMyZiJ9XSwidXNlcl9pZCI6IjgzODJlYzZkLTI1OWUtNDQxMy04NjIyLWFiYTFlZWM3MzdlOSJ9.um7HBxK3MH70rIS43ZZ_9ALRxvD5d7Dw8qc3vYg5e-Q'
 constructor(private http: HttpClient) { }
  //api/:companyID/campaigns
  private getHeaders(): HttpHeaders {
    return new HttpHeaders({
      'Authorization': `Bearer ${this.jwtToken}`
    });
  }
  getMailingLists(companyID: string, page: number = 1, limit: number = 10): Observable<any> {
    
    
    return this.http.get<any>(`${this.baseUrl}/${companyID}/mailinglist`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }

  createCampaign(companyID: string, campaign: any): Observable<any> {
    console.log('Sending campaign:', campaign);
    campaign.html = "<h1>Test</h1>";  
    return this.http.post<any>(`${this.baseUrl}/${companyID}/campaigns/${campaign.mailingListId}`, campaign, {
      headers: this.getHeaders()
    }).pipe(
      catchError(error => {
        console.error('Error creating campaign:', error);
        if (error.error && error.error.message) {
          return throwError(error.error.message);  // Throw user-friendly message from backend
        } else if (error.status === 400) {
          return throwError('Bad request. Please check campaign data.');  // Handle 400 error
        } else if (error.status === 401) {
          return throwError('Unauthorized. Please check your credentials.');  // Handle 401 error
        } else {
          return throwError('An unexpected error occurred.');  // Handle other errors
        }
      })
    );
  }

  getFromEmails(): Observable<any[]> {
    return this.http.get<any[]>('/api/fromEmails');
  }
}
