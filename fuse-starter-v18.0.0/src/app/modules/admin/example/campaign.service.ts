import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { throwError } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CampaignService {

  
 private companyID ="afa35ff6-4de5-4806-9a21-e0c2453d2834" ;
 private baseUrl = 'http://localhost:8080/api';
                    
 private jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55X2lkIjoiYjI3ZWU3N2MtOTA0My00MDAwLWI3ZTAtZjFhOTIwZGEyYzJmIiwiZXhwIjoxNzE5OTY2MDE4LCJpYXQiOjE3MTk3OTMyMTgsInJvbGVzIjpbeyJpZCI6ImMyNmQ0MGE2LTc0MzMtNDg1MC1hYmUyLWIwMzIxYzNkNDA3ZiIsIm5hbWUiOiJNYW5hZ2VyIiwiY29tcGFueV9pZCI6ImIyN2VlNzdjLTkwNDMtNDAwMC1iN2UwLWYxYTkyMGRhMmMyZiJ9XSwidXNlcl9pZCI6IjgzODJlYzZkLTI1OWUtNDQxMy04NjIyLWFiYTFlZWM3MzdlOSJ9.VJn51F28ekCyBpKnwMH-9olqj-FtTp715aBg5LGi2O4'
 constructor(private http: HttpClient) { }
  //api/:companyID/campaigns
  private getHeaders(): HttpHeaders {
    let headers = new HttpHeaders()
      .set('Authorization', `Bearer ${this.jwtToken}`)
      .set('Content-Type', 'application/json');
    return headers;
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
  getCampaigns( page: number = 1, limit: number = 10): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/${this.companyID}/campaigns`, {
      headers: this.getHeaders(),
      params: {
        page: page.toString(),
        limit: limit.toString()
      }
    });
  }
  createCampaign(companyID: string, campaign: any): Observable<any> {
    campaign.html = "<h1>Test</h1>";
    const headers = this.getHeaders();
    return this.http.post<any>(`${this.baseUrl}/${companyID}/campaigns/${campaign.mailingListId}`, campaign, {
      headers: headers
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

  getCampaignByID(companyID: string, campaignID: string): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/${companyID}/campaigns/${campaignID}`, {
      headers: this.getHeaders()
    });
  }

  updateCampaign(companyID: string, campaignID: string, campaign: any): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/${companyID}/campaigns/${campaignID}`, campaign, {
      headers: this.getHeaders()
    });
  }

 
}
